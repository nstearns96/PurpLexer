package lexer

import (
	_ "embed"
	"testing"
)

//go:embed test_syntax.xml
var syntaxData []byte

//go:embed invalid_test_syntax.xml
var invalidSyntaxData []byte

func TestLexer(t *testing.T) {
	lex := NewLexer()
	err := lex.LoadSyntaxXML(invalidSyntaxData)
	if err == nil {
		t.Error("Succesfully loaded invalid syntax data. This should have failed.")
		return
	}

	err = lex.LoadSyntaxXML(syntaxData)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		matchTerm   string
		matchString string
		expectMatch bool
		tokens      []string
	}{
		//Literal matching
		{
			matchTerm:   "foo",
			matchString: "",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   "foo",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   "foo",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   LiteralPrefix + "foo",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   LiteralPrefix + LiteralPrefix,
			matchString: LiteralPrefix,
			expectMatch: true,
			tokens: []string{
				LiteralPrefix,
			},
		},
		{
			matchTerm:   LiteralPrefix + TermPrefix,
			matchString: TermPrefix,
			expectMatch: true,
			tokens: []string{
				TermPrefix,
			},
		},
		{
			matchTerm:   LiteralPrefix + BuiltInPrefix,
			matchString: BuiltInPrefix,
			expectMatch: true,
			tokens: []string{
				BuiltInPrefix,
			},
		},
		//Term matching
		{
			matchTerm:   TermPrefix,
			matchString: "",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "foo",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "foo",
			matchString: " foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "foo",
			matchString: "foo ",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "bar",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooAndBar",
			matchString: "foobar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "fooAndBar",
			matchString: "foo bar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "fooAndBar",
			matchString: "foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooAndBar",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooOrBar",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "fooOrBar",
			matchString: "bar",
			expectMatch: true,
			tokens: []string{
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "fooOrBar",
			matchString: "foobar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooTerm",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "fooTerm",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooOptionalBar",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "fooOptionalBar",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooOptionalBar",
			matchString: "foobar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "fooOptionalBar",
			matchString: "foobaz",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooOptionalBar",
			matchString: "barfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "foo,bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "foo,foo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "foo, foo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "atLeastOneFoo",
			matchString: "foo ,foo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "fooManyBar",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "fooManyBar",
			matchString: "bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooManyBar",
			matchString: "foobar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "fooManyBar",
			matchString: "foo bar, bar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
				"bar",
			},
		},
		// Built in matching
		{
			matchTerm:   BuiltInPrefix,
			matchString: "",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix,
			matchString: "foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "foo",
			matchString: "",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "foo",
			matchString: "foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "label",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "label",
			matchString: "f00",
			expectMatch: true,
			tokens: []string{
				"f00",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "label",
			matchString: "f00",
			expectMatch: true,
			tokens: []string{
				"f00",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "label",
			matchString: "_",
			expectMatch: true,
			tokens: []string{
				"_",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "label",
			matchString: "_foo",
			expectMatch: true,
			tokens: []string{
				"_foo",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "label",
			matchString: "0",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "match",
			matchString: "",
			expectMatch: true,
			tokens: []string{
				"",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "match",
			matchString: "foo",
			expectMatch: true,
			tokens: []string{
				"foo",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "int",
			matchString: "",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "int",
			matchString: "foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "int",
			matchString: "-",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   BuiltInPrefix + "int",
			matchString: "0",
			expectMatch: true,
			tokens: []string{
				"0",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "int",
			matchString: "-1",
			expectMatch: true,
			tokens: []string{
				"-1",
			},
		},
		{
			matchTerm:   BuiltInPrefix + "int",
			matchString: "0f",
			expectMatch: false,
			tokens:      nil,
		},
		//Complex matching
		{
			matchTerm:   TermPrefix + "fooIntBar",
			matchString: "foo123bar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"123",
				"bar",
			},
		},
	}

	for _, tt := range tests {
		matchToks, matched := lex.MatchString(tt.matchString, tt.matchTerm)
		if matched != tt.expectMatch {
			if tt.expectMatch {
				t.Errorf("expected to match %s with term %s, but didn't", tt.matchString, tt.matchTerm)
				continue
			} else {
				t.Errorf("did not expect to match %s with term %s, but did", tt.matchString, tt.matchTerm)
				continue
			}
		}

		if !tt.expectMatch {
			continue
		}

		if len(tt.tokens) != len(matchToks) {
			t.Errorf("expected tokens do not match. Expected %v but got %v", tt.tokens, matchToks)
			continue
		}

		for tokIndex, tok := range matchToks {
			if tok != tt.tokens[tokIndex] {
				t.Errorf("expected tokens do not match. Expected %v but got %v", tt.tokens, matchToks)
				break
			}
		}
	}

	lex.ClearSyntax()
}
