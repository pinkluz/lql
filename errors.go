package lql

import "strings"

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
