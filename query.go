package query

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

// UnknownParseError if you get this something has gone very wrong
type UnknownParseError struct{}

func (u *UnknownParseError) Error() string {
	return "Unknown parse error"
}

// ParseError is for exposing known parse errors
type ParseError struct {
	messages []string
}

func (p *ParseError) Error() string {
	return strings.Join(p.messages, " | ")
}

// ParseNOP means you gave us nothing to do
type ParseNOP struct{}

func (p *ParseNOP) Error() string {
	return "Nothing was given to parse"
}

// Parse the passed in query
func Parse(query []byte) (*bleve.SearchRequest, error) {
	if query == nil || len(query) == 0 {
		return nil, &ParseNOP{}
	}

	toLex := &lex{
		search: &bleve.SearchRequest{
			Query:   nil,
			Size:    10,
			From:    0,
			Explain: false,
			Sort: search.SortOrder{
				&search.SortScore{Desc: true}},
		},
		// Bleve ignores case so we just downcase to make life easy
		input: bytes.NewReader(bytes.ToLower(query)),
		errs:  []string{},
	}

	// Start the parser generated by goyacc
	code := yyParse(toLex)

	if len(toLex.errs) > 0 {
		return nil, fmt.Errorf(strings.Join(toLex.errs, "\n"))
	}

	if code != 0 {
		return nil, &UnknownParseError{}
	}

	return toLex.search, nil
}