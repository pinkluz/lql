package lql

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"strings"

	bsq "github.com/blugelabs/bluge"
)

type lex struct {
	// Input string from the user
	input *bytes.Reader
	errs  []string

	// search is the end result of the lexer parser and what is given to the user
	query  bsq.Query
	fields []string

	// Hold current value of strings and other multi-byte operators (and, or, not, ...)
	lval      string
	lastToken int

	inQuote  bool
	inEscape bool
}

// Lex parses tokens and if needed sets the lval which will
// be what is used in the $1, $2, ..., $n vars in your goyacc
// parser.y file.
func (l *lex) Lex(lval *yySymType) int {
	var b byte
	var err error
	for {
		b, err = l.input.ReadByte()
		if err == io.EOF {
			// not a real error
			return 0
		}
		if err != nil {
			l.errs = append(l.errs, err.Error())
			return 0
		}

		// Keep going until we get non whitespace
		if l.byteToToken(b) != tWHITESPACE {
			break
		}
	}

	token := l.byteToToken(b)
	if rtoken := l.singleByteTokenToRetVal(token); rtoken != 0 {
		l.lastToken = rtoken
		return rtoken
	}

	// Looks like we are in a string and we will now keep
	// eating bytes until we are to the end of the string.
	// Then hack off any whitespace which should only happen
	// at the start of the string.
	var retStr string
	if token == tDOUBLEQUOTE {
		l.inQuote = true
		retStr = l.eatString()
	} else if token == tBACKSLASH {
		l.inEscape = true
		retStr = l.eatString()
	} else {
		retStr = string(b) + l.eatString()
	}

	// DEBUG: fmt.Println("RET:", retStr)

	// Always trim any space between the string and the next token
	retStr = strings.TrimSpace(retStr)

	// Check if this is a multi byte token such as and or or before just
	// freely passing it back and calling it a string.
	if rtoken := l.multiByteTokenToRetVal(retStr); rtoken != 0 {
		l.lastToken = rtoken
		return rtoken
	}

	// Set the left val for the parser to use
	lval.str = retStr

	// Clear everything for next Lex call
	l.lastToken = tSTRING
	l.inQuote = false
	l.inEscape = false

	return tSTRING
}

// Add a field that will limit the amount of field returned
// from our search.
func (l *lex) addReturnField(field string) {
	l.fields = append(l.fields, field)
}

// Final method that is called by the lexer to set
// the query that has been built.
func (l *lex) setSearchQuery(q bsq.Query) {
	l.query = q
}

func (l *lex) eatString() string {
	var retVal string
	for {
		b, err := l.input.ReadByte()
		if err != nil {
			return retVal
		}

		token := l.byteToToken(b)

		if token == tBACKSLASH && !l.inEscape {
			l.inEscape = true
		} else {
			// Hit the end of a double quote so cleanup and return
			// without including the " in the retVal
			if !l.inEscape {
				if token == tDOUBLEQUOTE && l.inQuote {
					break
				}

				// If we hit a known token we are no longer in a string
				// and should backout the last read byte and break
				if token != tSTRING && !l.inQuote {
					l.input.UnreadByte()
					break
				}
			}

			retVal = retVal + string(b)
			l.inEscape = false
		}
	}

	return retVal
}

func (l *lex) multiByteTokenToRetVal(val string) int {
	// check if we are on the right hand side of a k/v
	// pair before testing the value. or k=and causes issues
	switch l.lastToken {
	case tEQUAL, tGREATERTHAN,
		tLESSTHAN, tTILDE, tAND, tOR:

		return 0
	case 0: // Last token has not been set
		return 0
	}

	switch val {
	case "returns":
		return tRETURNS
	case "and":
		return tAND
	case "or":
		return tOR
	}

	return 0
}

// given a single token decide what if anything should be returned
// to the parser.
func (l *lex) singleByteTokenToRetVal(t int) int {

	if l.lastToken == tEQUAL && t == tLSQUAREBRACKET {
		return tLSQUAREBRACKET
	}

	// Special case for <= and >=
	if (l.lastToken == tGREATERTHAN || l.lastToken == tLESSTHAN) && t == tEQUAL {
		return tEQUAL
	}

	// !~
	if l.lastToken == tEXCLAMATION && t == tTILDE {
		return tTILDE
	}

	// =~
	if l.lastToken == tEQUAL && t == tTILDE {
		return tTILDE
	}

	// check if we are on the right hand side of a k/v
	// pair before testing the value. or k=\" causes issues
	switch l.lastToken {
	case tEQUAL, tGREATERTHAN,
		tLESSTHAN, tTILDE:

		return 0
	}

	// tSTRING and tDOUBLEQUOTE mean that we are about to enter
	// into a string.
	if t != tSTRING && t != tDOUBLEQUOTE {
		return t
	}

	return 0
}

// token is for checking tokens that can be decided by only looking
// at a single byte and possibly some other context.
func (l *lex) byteToToken(b byte) int {
	// always eat whitespace
	if b == ' ' {
		return tWHITESPACE
	}

	switch b {
	case '~':
		return tTILDE
	case '>':
		return tGREATERTHAN
	case '<':
		return tLESSTHAN
	case '(':
		return tLBRACKET
	case ')':
		return tRBRACKET
	case '=':
		return tEQUAL
	case '\\':
		return tBACKSLASH
	case '"':
		return tDOUBLEQUOTE
	case '!':
		return tEXCLAMATION
	case ',':
		return tCOMMA
	case '[':
		return tLSQUAREBRACKET
	case ']':
		return tRSQUAREBRACKET
	default:
		return tSTRING
	}
}

func (l *lex) Error(s string) {
	l.errs = append(l.errs, s)
}

// handle building a range query since we just take a tSTRING tCOMMA tSTRING
// for the value to keep it simple. We verify here that tSTRING is a valid
// float or is * which means 0 for min or max float for max.
func (l *lex) makeRangeQuery(min string, max string) *bsq.NumericRangeQuery {
	var fmin float64
	if min == "*" {
		fmin = 0
	} else {
		fmin = l.strToFloat64(min)
	}

	var fmax float64
	if max == "*" {
		fmax = math.MaxFloat64
	} else {
		fmax = l.strToFloat64(max)
	}

	query := bsq.NewNumericRangeQuery(fmin, fmax)

	return query
}

func (l *lex) strToFloat64(val string) float64 {
	lval, err := strconv.ParseFloat(val, 64)
	if err != nil {
		l.errs = append(l.errs, err.Error())
	}

	return lval
}
