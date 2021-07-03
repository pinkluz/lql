%{
package lql

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

%}

%union{
	num   float64
	str   string
	query bsq.Query
}

// multi byte tokens
%token tOR tAND tRETURNS tGREATERTHAN tLESSTHAN tNUMBER tSTRING

// single byte tokens
%token tEQUAL tEXCLAMATION tLBRACKET tRBRACKET tASTERISK
%token tTILDE  tCOMMA tLSQUAREBRACKET tRSQUAREBRACKET

// Tokens only used internally and not in the grammar def
%token tWHITESPACE tBACKSLASH tDOUBLEQUOTE

// https://stackoverflow.com/questions/12876543/left-and-right-in-yacc
// https://www.gnu.org/software/bison/manual/html_node/Precedence-Decl.html
%left tCOMMA
%left tOR
%left tAND
%left tLBRACKET tRBRACKET

%type <num>    tNUMBER
%type <str>    tSTRING
%type <str>    tASTERISK
%type <query>  lql
%type <query>  searchParts

%start main

// Read this for some info although figuring out exactly how yacc works
//
// http://dinosaur.compilertools.net/yacc/
// https://arcb.csc.ncsu.edu/~mueller/codeopt/codeopt00/y_man.pdf

%%

main:
searchParts {
	logDebugGrammar("searchParts")
	// TODO: add finish function here that can do any cleanup
	// such as closing out any remaining logical operators that
	// were in the queue.
	query := $1

	yylex.(*lex).setSearchQuery(query)
}
|
searchParts tRETURNS returnPart {
	logDebugGrammar("searchParts tRETURNS returnPart")
	// TODO: add finish function here that can do any cleanup
	// such as closing out any remaining logical operators that
	// were in the queue.
	query := $1

	yylex.(*lex).setSearchQuery(query)
}
;

returnPart:
tSTRING {
	field := $1
	yylex.(*lex).addReturnField(field)
}
|
returnPart tCOMMA returnPart {
	logDebugGrammar("returnPart tCOMMA returnPart")
}
|
tLBRACKET returnPart tRBRACKET {
	logDebugGrammar("tLBRACKET returnPart tRBRACKET")
}
;


searchParts:
lql {
	logDebugGrammar("lql")
	$$ = $1
}
|
searchParts tOR searchParts {
	logDebugGrammar("searchParts tOR searchParts")
	queryOne := $1
	queryTwo := $3

	query := bsq.NewBooleanQuery()
	query.AddShould(queryOne)
	query.AddShould(queryTwo)

	$$ = query
}
|
searchParts tAND searchParts {
	logDebugGrammar("searchParts tAND searchParts")
	queryOne := $1
	queryTwo := $3

	query := bsq.NewBooleanQuery()
	query.AddMust(queryOne)
	query.AddMust(queryTwo)

	$$ = query
}
|
tLBRACKET searchParts tRBRACKET {
	logDebugGrammar("tLBRACKET searchParts tRBRACKET")
	queryPar := $2

	query := bsq.NewBooleanQuery()
	query.AddMust(queryPar)

	$$ = query
};

lql:
tSTRING tEQUAL tSTRING {
	// k=v
	logDebugGrammar("tSTRING tEQUAL tSTRING")
	key := $1
	value := $3

	query := bsq.NewMatchQuery(value)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEQUAL tNUMBER {
	// k=1
	logDebugGrammar("tSTRING tEQUAL tNUMBER")
	key := $1
	value := $3

	query := bsq.NewNumericRangeInclusiveQuery(value, value, true, true)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEXCLAMATION tEQUAL tSTRING {
	// k!=v
	logDebugGrammar("tSTRING tEXCLAMATION tEQUAL tSTRING")
	key := $1
	value := $4

	match := bsq.NewMatchQuery(value)
	match.SetField(key)
	query := bsq.NewBooleanQuery()
	query.AddMustNot(match)

	$$ = query
}
|
tSTRING tEXCLAMATION tEQUAL tNUMBER {
	// k!=1
	logDebugGrammar("tSTRING tEXCLAMATION tEQUAL tNUMBER")
	key := $1
	value := $4

	match := bsq.NewNumericRangeInclusiveQuery(value, value, true, true)
	match.SetField(key)
	query := bsq.NewBooleanQuery()
	query.AddMustNot(match)

	$$ = query
}
|
tSTRING tEQUAL tLSQUAREBRACKET tNUMBER tCOMMA tNUMBER tRSQUAREBRACKET {
	// k=[v1,v2]
	logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tNUMBER tCOMMA tNUMBER tRSQUAREBRACKET")
	key := $1
	min := $4
	max := $6

	// Read comment about func for more info
	query := bsq.NewNumericRangeQuery(min, max)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEQUAL tLSQUAREBRACKET tASTERISK tCOMMA tNUMBER tRSQUAREBRACKET {
	// k=[*,v2]
	logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tASTERISK tCOMMA tNUMBER tRSQUAREBRACKET")
	key := $1
	max := $6

	// Read comment about func for more info
	query := bsq.NewNumericRangeQuery(0.0, max)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEQUAL tLSQUAREBRACKET tNUMBER tCOMMA tASTERISK tRSQUAREBRACKET {
	// k=[v1,*]
	logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tNUMBER tCOMMA tASTERISK tRSQUAREBRACKET")
	key := $1
	min := $4

	// Read comment about func for more info
	query := bsq.NewNumericRangeQuery(min, math.MaxFloat64)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEQUAL tLSQUAREBRACKET tASTERISK tCOMMA tASTERISK tRSQUAREBRACKET {
	// k=[*,*]
	logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tASTERISK tCOMMA tASTERISK tRSQUAREBRACKET")
	key := $1

	// Read comment about func for more info
	query := bsq.NewNumericRangeQuery(0.0, math.MaxFloat64)
	query.SetField(key)

	$$ = query
}
|
tSTRING tGREATERTHAN tNUMBER {
	// k>v
	logDebugGrammar("tSTRING tGREATERTHAN tNUMBER")
	key := $1
	value := $3
	setInclusiveRangeQuery := false

	query := bsq.NewNumericRangeInclusiveQuery(value, math.MaxFloat64,
		setInclusiveRangeQuery, setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING tLESSTHAN tNUMBER {
	// k<v
	logDebugGrammar("tSTRING tLESSTHAN tNUMBER")
	key := $1
	value := $3
	setInclusiveRangeQuery := false

	query := bsq.NewNumericRangeInclusiveQuery(0.0, value,
		setInclusiveRangeQuery, setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING tGREATERTHAN tEQUAL tNUMBER {
	// k>=v
	logDebugGrammar("tSTRING tGREATERTHAN tEQUAL tNUMBER")
	key := $1
	value := $4
	setInclusiveRangeQuery := true

	query := bsq.NewNumericRangeInclusiveQuery(value, math.MaxFloat64,
		setInclusiveRangeQuery, setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING tLESSTHAN tEQUAL tNUMBER {
	// k<=v
	logDebugGrammar("tSTRING tLESSTHAN tEQUAL tNUMBER")
	key := $1
	value := $4
	setInclusiveRangeQuery := true

	query := bsq.NewNumericRangeInclusiveQuery(0.0, value,
		setInclusiveRangeQuery, setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING  tEQUAL tTILDE tSTRING {
	// k=~v
	logDebugGrammar("tSTRING  tEQUAL tTILDE tSTRING")
	key := $1
	value := $4

	query := bsq.NewRegexpQuery(value)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEXCLAMATION tTILDE tSTRING {
	// k!~v
	logDebugGrammar("tSTRING tEXCLAMATION tTILDE tSTRING")
	key := $1
	value := $4

	match := bsq.NewRegexpQuery(value)
	match.SetField(key)
	query := bsq.NewBooleanQuery()
	query.AddMustNot(match)

	$$ = query
}
;
