package lql

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/blugelabs/bluge"
)

func populateBleve(num int) (*bluge.Writer, error) {
	config := bluge.InMemoryOnlyConfig()
	index, err := bluge.OpenWriter(config)
	if err != nil {
		return nil, err
	}

	// REMINDER: bleve lowercases everything when you are performing a search
	for i := 1; i <= num; i++ {
		doc := bluge.NewDocument(fmt.Sprintf("u%d", i))
		doc.AddField(bluge.NewTextField("name", fmt.Sprintf("u%d", i)).StoreValue())
		doc.AddField(bluge.NewTextField("test", "test").StoreValue())
		doc.AddField(bluge.NewTextField("extra", fmt.Sprintf("u%d", i)).StoreValue())
		doc.AddField(bluge.NewNumericField("range", float64(i)).StoreValue())

		index.Insert(doc)
	}

	return index, nil
}

func TestQueryResults(t *testing.T) {
	index, err := populateBleve(5)
	if err != nil {
		t.Errorf("failed to populate test index: %s", err)
		return
	}

	tests := []struct {
		input       []byte
		documents   []string
		fields      []string
		shouldError bool
	}{
		// Basic query test
		{
			input: []byte("name=u1"),
			documents: []string{
				"u1",
			},
		},
		// Test basic or logic
		{
			input: []byte("name=u1 or name=u2"),
			documents: []string{
				"u1",
				"u2",
			},
		},
		// Test basic and logic
		{
			input: []byte("name=u1 and test=test"),
			documents: []string{
				"u1",
			},
		},
		// Test basic and and or logic
		{
			input: []byte("name=u1 and test=test or name=u2"),
			documents: []string{
				"u1",
				"u2",
			},
		},
		// More complex and and or logic
		{
			input: []byte("name=u1 and test=test or name=u2 and test=test"),
			documents: []string{
				"u1",
				"u2",
			},
		},
		// More complex and and or logic with braces
		{
			input: []byte("(name=u1 and test=test or name=u3) and name=u3"),
			documents: []string{
				"u3",
			},
		},
		// More complex and and or logic with braces 2
		{
			input: []byte("(name=u1 and test=test or name=u3 and extra=u3) and name=u3"),
			documents: []string{
				"u3",
			},
		},
		// Test not logic
		{
			input: []byte("name!=u1"),
			documents: []string{
				"u2",
				"u3",
				"u4",
				"u5",
			},
		},
		// Test not logic with and
		{
			input: []byte("name!=u1 and name!=u2"),
			documents: []string{
				"u3",
				"u4",
				"u5",
			},
		},
		// Test not logic with and and or
		{
			input: []byte("name!=u1 and (name!=u2 or test=wow)"),
			documents: []string{
				"u3",
				"u4",
				"u5",
			},
		},
		// Test not logic with and and or
		{
			input: []byte("name!=u1 returns test"),
			documents: []string{
				"u2",
				"u3",
				"u4",
				"u5",
			},
			fields: []string{
				"test",
			},
		},
		// Test not logic with and and or
		{
			input: []byte("name!=u1 returns (test, extra, name)"),
			documents: []string{
				"u2",
				"u3",
				"u4",
				"u5",
			},
			fields: []string{
				"test",
				"extra",
				"name",
			},
		},
		// Test not logic with and and or and do quotes
		{
			input: []byte("name!=\"u1\" returns (\"test\", extra, name)"),
			documents: []string{
				"u2",
				"u3",
				"u4",
				"u5",
			},
			fields: []string{
				"test",
				"extra",
				"name",
			},
		},
		// Test not logic with and and or and do quotes
		{
			input: []byte("name!=\"u1\" returns (\"test\", extra, name)"),
			documents: []string{
				"u2",
				"u3",
				"u4",
				"u5",
			},
			fields: []string{
				"test",
				"extra",
				"name",
			},
		},
		// Test not logic with and and or and do quotes for key and value.
		// This matches all documents because all documents do not have a
		// name! field not equal to u1.
		{
			input: []byte("\"name!\"!=\"u1\""),
			documents: []string{
				"u1",
				"u2",
				"u3",
				"u4",
				"u5",
			},
		},
		// Make sure you can use a double bracket as a value
		{
			input:     []byte("name=\\\""),
			documents: []string{},
		},
		{
			input:     []byte("name=\"\\\"\""),
			documents: []string{},
		},
		// We don't care about the whitespace but it's ugly so don't
		{
			input: []byte("name =   u5"),
			documents: []string{
				"u5",
			},
		},
		// basic range query
		{
			input: []byte("range=[1.2,3.4] and name=u3"),
			documents: []string{
				"u3",
			},
		},
		// wildcard range. * for min sets 0 and * for max sets max float val
		{
			input: []byte("range=[*,*] and (name=u3 or name=u5)"),
			documents: []string{
				"u3",
				"u5",
			},
		},
		// wildcard range. * for min sets 0 and * for max sets max float val
		{
			input: []byte("range<4 and range>2"),
			documents: []string{
				"u3",
			},
		},
		// wildcard range. * for min sets 0 and * for max sets max float val
		{
			input: []byte("range<=4 and range>=2"),
			documents: []string{
				"u2",
				"u3",
				"u4",
			},
		},
		// Test regex query
		{
			input: []byte("name=~\"u.*\""),
			documents: []string{
				"u1",
				"u2",
				"u3",
				"u4",
				"u5",
			},
		},
		// Test regex query inverse
		{
			input: []byte("name!~u1"),
			documents: []string{
				"u2",
				"u3",
				"u4",
				"u5",
			},
		},
		// make sure keyword on right or left side of = doesn't cause invalid syntax
		{
			input:     []byte("and=and and or=or or returns=returns"),
			documents: []string{},
		},
		// Should error
		{
			input:       []byte("name"),
			shouldError: true,
		},
		// Should error
		{
			input:       []byte("range=[r,k]"),
			shouldError: true,
		},
		{
			// Verify that the lexer can tell the difference between a number and a string
			input: []byte("range=1"),
			documents: []string{
				"u1",
			},
		},
		{
			// Verify that the lexer can tell the difference between a float with point and a string
			input: []byte("range=1.0"),
			documents: []string{
				"u1",
			},
		},
		{
			// This is not a valid decimal and should do a string match and find nothing
			input:     []byte("range=1.0.0"),
			documents: []string{},
		},
		{
			// Since range is a numeric field and putting a field in double quotes forces
			// it to search for a string value only this should fail.
			input:     []byte("range=\"1\""),
			documents: []string{},
		},
	}

	// DEBUG: enable the following when debugging parser/lexer issues
	// yyDebug = 10
	// yyErrorVerbose = true
	// debugGrammer = true

	for _, test := range tests {
		// fmt.Println("testing:", string(test.input))
		query, fields, err := Parse(test.input)
		if err == nil && test.shouldError {
			t.Errorf("Should have got an error for '%s' but did not", string(test.input))
			continue
		}
		if err != nil && test.shouldError {
			// this should error so we don't care
			continue
		}
		if err != nil {
			t.Errorf("Parse Error: %v for '%s'", err, string(test.input))
			continue
		}

		reader, _ := index.Reader()
		req := bluge.NewAllMatches(query)

		searchResults, err := reader.Search(context.Background(), req)
		if err != nil {
			t.Errorf("failed to get results for query '%s': %s", string(test.input), err)
		}

		if len(test.fields) > 0 {
			if len(test.fields) != len(fields) {
				t.Errorf("%d fields found for query '%s' but was expecting %d",
					len(fields), string(test.input), len(test.fields))
			} else {
				sort.Strings(test.fields)
				sort.Strings(fields)

				for i, v := range fields {
					if v != test.fields[i] {
						t.Errorf("Looking for field '%s' in search query but got '%s' for query '%s'",
							test.fields[i], v, string(test.input))
						break
					}
				}
			}
		}

		// Pull all the doc names out so we can match against them
		docsFound := []string{}
		for {
			doc, err := searchResults.Next()
			if err != nil || doc == nil {
				break
			}

			doc.VisitStoredFields(func(field string, value []byte) bool {
				if field == "name" {
					docsFound = append(docsFound, string(value))
				}
				return true
			})
		}

		if len(docsFound) != len(test.documents) {
			t.Errorf("%d documents found for query '%s' but was expecting %d",
				len(docsFound), string(test.input), len(test.documents))
		} else {
			sort.Strings(docsFound)
			sort.Strings(test.documents)

			for i, v := range docsFound {
				if v != test.documents[i] {
					t.Errorf("Looking for document %s in result but got %s for query '%s'",
						test.documents[i], v, string(test.input))
					break
				}
			}

			// fmt.Println("Results Matched")
		}
	}
}
