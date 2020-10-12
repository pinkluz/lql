%{
package query

import (
	"fmt"

	bsq "github.com/blevesearch/bleve/search/query"
)

//go:generate go run golang.org/x/tools/cmd/goyacc -l -o parser.go parser.y

type queryOperator int

const (
	nop queryOperator = iota
	andQuery
	orQuery
)

var debugGrammer bool

func logDebugGrammar(msg interface{}) {
	if debugGrammer {
		fmt.Printf("%+v\n", msg)
	}
}

%}

%union{
	str string
	query bsq.Query
}

// multi byte tokens
%token tOR tAND tRETURNS tGREATERTHAN tLESSTHAN

// single byte tokens
%token tSTRING tEQUAL tEXCLAMATION tLBRACKET tRBRACKET
%token tTILDE  tCOMMA tLSQUAREBRACKET tRSQUAREBRACKET

// Tokens only used internally and not in the grammar def
%token tWHITESPACE tBACKSLASH tDOUBLEQUOTE

// https://stackoverflow.com/questions/12876543/left-and-right-in-yacc
// https://www.gnu.org/software/bison/manual/html_node/Precedence-Decl.html
%left tCOMMA
%left tOR
%left tAND
%left tLBRACKET tRBRACKET

%type <str>      tSTRING
%type <query>    lql
%type <query>    searchParts

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

	query := bsq.NewDisjunctionQuery([]bsq.Query{
		queryOne, queryTwo})
	$$ = query
}
|
searchParts tAND searchParts {
	logDebugGrammar("searchParts tAND searchParts")
	queryOne := $1
	queryTwo := $3

	query := bsq.NewConjunctionQuery([]bsq.Query{
		queryOne, queryTwo})
	$$ = query
}
|
tLBRACKET searchParts tRBRACKET {
	logDebugGrammar("tLBRACKET searchParts tRBRACKET")
	queryPar := $2

	query := bsq.NewBooleanQuery([]bsq.Query{}, []bsq.Query{}, []bsq.Query{})
	query.AddMust(queryPar)

	$$ = query
};

lql:
tSTRING tEQUAL tSTRING {
	// k=v
	logDebugGrammar("tSTRING tEQUAL tSTRING")
	key := $1
	value := $3

	query := bsq.NewRegexpQuery(value)
	query.SetField(key)

	$$ = query
}
|
tSTRING tEXCLAMATION tEQUAL tSTRING {
	// k!=v
	logDebugGrammar("tSTRING tEXCLAMATION tEQUAL tSTRING")
	key := $1
	value := $4

	match := bsq.NewMatchPhraseQuery(value)
	match.SetField(key)
	query := bsq.NewBooleanQuery([]bsq.Query{}, []bsq.Query{}, []bsq.Query{})
	query.AddMustNot(match)

	$$ = query
}
|
tSTRING tEQUAL tLSQUAREBRACKET tSTRING tCOMMA tSTRING tRSQUAREBRACKET {
	// k=[v1,v2]
	logDebugGrammar("tSTRING tEQUAL tLSQUAREBRACKET tSTRING tCOMMA tSTRING tRSQUAREBRACKET")
	key := $1
	min := $4
	max := $6

	// Read comment about func for more info
	query := yylex.(*lex).makeRangeQuery(min, max)
	query.SetField(key)

	$$ = query
}
|
tSTRING tGREATERTHAN tSTRING {
	// k>v
	logDebugGrammar("tSTRING tGREATERTHAN tSTRING")
	key := $1
	value := $3
	setInclusiveRangeQuery := false

	fval := yylex.(*lex).strToFloat64(value)

	query := bsq.NewNumericRangeInclusiveQuery(&fval, nil,
		&setInclusiveRangeQuery, &setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING tLESSTHAN tSTRING {
	// k<v
	logDebugGrammar("tSTRING tLESSTHAN tSTRING")
	key := $1
	value := $3
	setInclusiveRangeQuery := false

	fval := yylex.(*lex).strToFloat64(value)

	query := bsq.NewNumericRangeInclusiveQuery(nil, &fval,
		&setInclusiveRangeQuery, &setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING tGREATERTHAN tEQUAL tSTRING {
	// k>=v
	logDebugGrammar("tSTRING tGREATERTHAN tEQUAL tSTRING")
	key := $1
	value := $4
	setInclusiveRangeQuery := true

	fval := yylex.(*lex).strToFloat64(value)

	query := bsq.NewNumericRangeInclusiveQuery(&fval, nil,
		&setInclusiveRangeQuery, &setInclusiveRangeQuery)
	query.SetField(key)

	$$ = query
}
|
tSTRING tLESSTHAN tEQUAL tSTRING {
	// k<=v
	logDebugGrammar("tSTRING tLESSTHAN tEQUAL tSTRING")
	key := $1
	value := $4
	setInclusiveRangeQuery := true

	fval := yylex.(*lex).strToFloat64(value)

	query := bsq.NewNumericRangeInclusiveQuery(nil, &fval,
		&setInclusiveRangeQuery, &setInclusiveRangeQuery)
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
	query := bsq.NewBooleanQuery([]bsq.Query{}, []bsq.Query{}, []bsq.Query{})
	query.AddMustNot(match)

	$$ = query
}
;