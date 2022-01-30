// Code generated by goyacc -l -o parser.go parser.y. DO NOT EDIT.
package lql

import __yyfmt__ "fmt"

import (
	"fmt"
	"math"

	bsq "github.com/blugelabs/bluge"
)

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go parser.y

var debugGrammer bool

func logDebugGrammar(msg interface{}) {
	if debugGrammer {
		fmt.Printf("%+v\n", msg)
	}
}

type yySymType struct {
	yys   int
	num   float64
	str   string
	query bsq.Query
}

const tOR = 57346
const tAND = 57347
const tRETURNS = 57348
const tGREATERTHAN = 57349
const tLESSTHAN = 57350
const tNUMBER = 57351
const tSTRING = 57352
const tEQUAL = 57353
const tEXCLAMATION = 57354
const tLBRACKET = 57355
const tRBRACKET = 57356
const tASTERISK = 57357
const tTILDE = 57358
const tCOMMA = 57359
const tLSQUAREBRACKET = 57360
const tRSQUAREBRACKET = 57361
const tWHITESPACE = 57362
const tBACKSLASH = 57363
const tDOUBLEQUOTE = 57364

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"tOR",
	"tAND",
	"tRETURNS",
	"tGREATERTHAN",
	"tLESSTHAN",
	"tNUMBER",
	"tSTRING",
	"tEQUAL",
	"tEXCLAMATION",
	"tLBRACKET",
	"tRBRACKET",
	"tASTERISK",
	"tTILDE",
	"tCOMMA",
	"tLSQUAREBRACKET",
	"tRSQUAREBRACKET",
	"tWHITESPACE",
	"tBACKSLASH",
	"tDOUBLEQUOTE",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 55

var yyAct = [...]int{
	14, 21, 20, 51, 50, 49, 48, 43, 23, 41,
	22, 42, 30, 30, 24, 46, 44, 31, 32, 25,
	15, 47, 45, 16, 33, 7, 8, 36, 35, 12,
	13, 40, 39, 10, 11, 19, 5, 2, 28, 4,
	29, 26, 9, 27, 37, 17, 18, 34, 38, 7,
	8, 6, 8, 1, 3,
}

var yyPact = [...]int{
	26, -1000, 45, -1000, 26, 22, 10, 26, 26, 21,
	-8, 3, 32, 29, -4, -1000, 10, 47, -1000, -1000,
	-1000, -1000, 9, 37, 18, 34, -1000, 39, -1000, 23,
	10, -5, -6, -10, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 7, 6, -13, -14, -15, -16, -1000, -1000,
	-1000, -1000,
}

var yyPgo = [...]int{
	0, 54, 37, 53, 0,
}

var yyR1 = [...]int{
	0, 3, 3, 4, 4, 4, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1,
}

var yyR2 = [...]int{
	0, 1, 3, 1, 3, 3, 1, 3, 3, 3,
	3, 3, 4, 4, 7, 7, 7, 7, 3, 3,
	4, 4, 4, 4,
}

var yyChk = [...]int{
	-1000, -3, -2, -1, 13, 10, 6, 4, 5, -2,
	11, 12, 7, 8, -4, 10, 13, -2, -2, 14,
	10, 9, 18, 16, 11, 16, 9, 11, 9, 11,
	17, -4, 9, 15, 10, 10, 9, 10, 9, 9,
	-4, 14, 17, 17, 9, 15, 9, 15, 19, 19,
	19, 19,
}

var yyDef = [...]int{
	0, -2, 1, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 2, 3, 0, 7, 8, 9,
	10, 11, 0, 0, 0, 0, 18, 0, 19, 0,
	0, 0, 0, 0, 22, 12, 13, 23, 20, 21,
	4, 5, 0, 0, 0, 0, 0, 0, 14, 16,
	15, 17,
}

var yyTok1 = [...]int{
	1,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22,
}

var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			logDebugGrammar("searchParts")
			// TODO: add finish function here that can do any cleanup
			// such as closing out any remaining logical operators that
			// were in the queue.
			query := yyDollar[1].query

			yylex.(*lex).setSearchQuery(query)
		}
	case 2:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			logDebugGrammar("searchParts tRETURNS returnPart")
			// TODO: add finish function here that can do any cleanup
			// such as closing out any remaining logical operators that
			// were in the queue.
			query := yyDollar[1].query

			yylex.(*lex).setSearchQuery(query)
		}
	case 3:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			field := yyDollar[1].str
			yylex.(*lex).addReturnField(field)
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			logDebugGrammar("returnPart tCOMMA returnPart")
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			logDebugGrammar("tLBRACKET returnPart tRBRACKET")
		}
	case 6:
		yyDollar = yyS[yypt-1 : yypt+1]
		{
			logDebugGrammar("lql")
			yyVAL.query = yyDollar[1].query
		}
	case 7:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			logDebugGrammar("searchParts tOR searchParts")
			queryOne := yyDollar[1].query
			queryTwo := yyDollar[3].query

			query := bsq.NewBooleanQuery()
			query.AddShould(queryOne)
			query.AddShould(queryTwo)

			yyVAL.query = query
		}
	case 8:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			logDebugGrammar("searchParts tAND searchParts")
			queryOne := yyDollar[1].query
			queryTwo := yyDollar[3].query

			query := bsq.NewBooleanQuery()
			query.AddMust(queryOne)
			query.AddMust(queryTwo)

			yyVAL.query = query
		}
	case 9:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			logDebugGrammar("tLBRACKET searchParts tRBRACKET")
			queryPar := yyDollar[2].query

			query := bsq.NewBooleanQuery()
			query.AddMust(queryPar)

			yyVAL.query = query
		}
	case 10:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			// k=v
			logDebugGrammar("tSTRING tEQUAL tSTRING")
			key := yyDollar[1].str
			value := yyDollar[3].str

			query := bsq.NewMatchQuery(value)
			query.SetField(key)

			yyVAL.query = query
		}
	case 11:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			// k=1
			logDebugGrammar("tSTRING tEQUAL tNUMBER")
			key := yyDollar[1].str
			value := yyDollar[3].num

			query := bsq.NewNumericRangeInclusiveQuery(value, value, true, true)
			query.SetField(key)

			yyVAL.query = query
		}
	case 12:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			// k!=v
			logDebugGrammar("tSTRING tEXCLAMATION tEQUAL tSTRING")
			key := yyDollar[1].str
			value := yyDollar[4].str

			match := bsq.NewMatchQuery(value)
			match.SetField(key)
			query := bsq.NewBooleanQuery()
			query.AddMustNot(match)

			yyVAL.query = query
		}
	case 13:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			// k!=1
			logDebugGrammar("tSTRING tEXCLAMATION tEQUAL tNUMBER")
			key := yyDollar[1].str
			value := yyDollar[4].num

			match := bsq.NewNumericRangeInclusiveQuery(value, value, true, true)
			match.SetField(key)
			query := bsq.NewBooleanQuery()
			query.AddMustNot(match)

			yyVAL.query = query
		}
	case 14:
		yyDollar = yyS[yypt-7 : yypt+1]
		{
			// k=[v1,v2]
			logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tNUMBER tCOMMA tNUMBER tRSQUAREBRACKET")
			key := yyDollar[1].str
			min := yyDollar[4].num
			max := yyDollar[6].num

			// Read comment about func for more info
			query := bsq.NewNumericRangeQuery(min, max)
			query.SetField(key)

			yyVAL.query = query
		}
	case 15:
		yyDollar = yyS[yypt-7 : yypt+1]
		{
			// k=[*,v2]
			logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tASTERISK tCOMMA tNUMBER tRSQUAREBRACKET")
			key := yyDollar[1].str
			max := yyDollar[6].num

			// Read comment about func for more info
			query := bsq.NewNumericRangeQuery(0.0, max)
			query.SetField(key)

			yyVAL.query = query
		}
	case 16:
		yyDollar = yyS[yypt-7 : yypt+1]
		{
			// k=[v1,*]
			logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tNUMBER tCOMMA tASTERISK tRSQUAREBRACKET")
			key := yyDollar[1].str
			min := yyDollar[4].num

			// Read comment about func for more info
			query := bsq.NewNumericRangeQuery(min, math.MaxFloat64)
			query.SetField(key)

			yyVAL.query = query
		}
	case 17:
		yyDollar = yyS[yypt-7 : yypt+1]
		{
			// k=[*,*]
			logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tASTERISK tCOMMA tASTERISK tRSQUAREBRACKET")
			key := yyDollar[1].str

			// Read comment about func for more info
			query := bsq.NewNumericRangeQuery(0.0, math.MaxFloat64)
			query.SetField(key)

			yyVAL.query = query
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			// k>v
			logDebugGrammar("tSTRING tGREATERTHAN tNUMBER")
			key := yyDollar[1].str
			value := yyDollar[3].num
			setInclusiveRangeQuery := false

			query := bsq.NewNumericRangeInclusiveQuery(value, math.MaxFloat64,
				setInclusiveRangeQuery, setInclusiveRangeQuery)
			query.SetField(key)

			yyVAL.query = query
		}
	case 19:
		yyDollar = yyS[yypt-3 : yypt+1]
		{
			// k<v
			logDebugGrammar("tSTRING tLESSTHAN tNUMBER")
			key := yyDollar[1].str
			value := yyDollar[3].num
			setInclusiveRangeQuery := false

			query := bsq.NewNumericRangeInclusiveQuery(0.0, value,
				setInclusiveRangeQuery, setInclusiveRangeQuery)
			query.SetField(key)

			yyVAL.query = query
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			// k>=v
			logDebugGrammar("tSTRING tGREATERTHAN tEQUAL tNUMBER")
			key := yyDollar[1].str
			value := yyDollar[4].num
			setInclusiveRangeQuery := true

			query := bsq.NewNumericRangeInclusiveQuery(value, math.MaxFloat64,
				setInclusiveRangeQuery, setInclusiveRangeQuery)
			query.SetField(key)

			yyVAL.query = query
		}
	case 21:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			// k<=v
			logDebugGrammar("tSTRING tLESSTHAN tEQUAL tNUMBER")
			key := yyDollar[1].str
			value := yyDollar[4].num
			setInclusiveRangeQuery := true

			query := bsq.NewNumericRangeInclusiveQuery(0.0, value,
				setInclusiveRangeQuery, setInclusiveRangeQuery)
			query.SetField(key)

			yyVAL.query = query
		}
	case 22:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			// k=~v
			logDebugGrammar("tSTRING  tEQUAL tTILDE tSTRING")
			key := yyDollar[1].str
			value := yyDollar[4].str

			query := bsq.NewRegexpQuery(value)
			query.SetField(key)

			yyVAL.query = query
		}
	case 23:
		yyDollar = yyS[yypt-4 : yypt+1]
		{
			// k!~v
			logDebugGrammar("tSTRING tEXCLAMATION tTILDE tSTRING")
			key := yyDollar[1].str
			value := yyDollar[4].str

			match := bsq.NewRegexpQuery(value)
			match.SetField(key)
			query := bsq.NewBooleanQuery()
			query.AddMustNot(match)

			yyVAL.query = query
		}
	}
	goto yystack /* stack new state and value */
}
