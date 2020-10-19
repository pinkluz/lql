<h1 align="center">Lily Query Language</h1>

<p align="center">A Solr like query language for the bluge index (beta)</p>

<br/>

Lily Query Language (LQL) for short is used to build a bluge search result. It's
somewhat modeled after the Solr query language but is currently much more simple.
In addition to building a query the language also lets you limit the result set.

### Examples

```
# Find any match of string for a given asset
tag=U123

# Whitespace doesn't matter this is the same as above
tag  =   U123

# Additinally double quotes are used to escape invalid characters
"tag!=" = "U123"

# Not equal operator is also valid
tag!=U123

# Check greater than or less than
value>5

# use logical operators like and/or
value>5 and tag=U123 or value=10

# override operator precidence with brackets
value>5 and (tag=U123 or value=10)

# limit the field that is returned
tag=U123 returns tag

# have multiple returned fields
tag=U123 returns tag, value, test

# Find something between two values
value=[1.3,2.2]

# Use a regex query
tag=~U.*

# Invert your regex
tag!~U.*
```

### Specification

#### Lexing

Everything in LQL is a string or an operator. Although under the hood the lexer will take
a given string and convert it to a float64 for operators `>`, `<`, `<=`, `>=`, and `=[,]`.
If it can not be parsed to a float64 an error will be returned.

`strings` - The `tSTRING` token defines a string in the LQL grammar. A string is defined as
any valid character that is not a single character operator. Valid single character operators
are `=`, `!`, `(`, `)`, `~`, `,`, `[`, `]`, `\`, and `"`. `\` is reserved for escaping characters
and a string that starts in an unescaped `"` must end in it. You only need to escape
double quotes inside a quoted string.

`operators` - The following operators are defined in the LQL grammer. Multi-byte operators are
`and`, `or`, `returns`, `>=`, `<=`, `=~`, and`!~`. Single-byte operators are `=`, `!`, `(`,
`)`, `~`, `,`, `[`, `]`, `\`, and `"`.

#### Parser

The up to date grammer can be seen in `parser.y`.

### Development

To rebuild `parser.go` you can run the following command

```
go run golang.org/x/tools/cmd/goyacc -l -o parser.go parser.y
```

Additionally if you want to add/delete/change something parser.y has been documented with
lots of good external resouorces to read.

Before putting up a chanage against this make sure to run the following.

```
go generate ./... && go test -v ./...
```