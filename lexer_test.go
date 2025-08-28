package lexer

import (
	_ "embed"
	"testing"
)

//go:embed test_syntax.xml
var syntaxData []byte

//go:embed invalid_test_syntax.xml
var invalidSyntaxData []byte

type MatchingTestTable []struct {
	matchTerm   string
	matchString string
	expectMatch bool
	tokens      []string
}

func MakeTestLexer() (*Lexer, error) {
	lex := NewLexer()
	err := lex.LoadSyntaxXML(syntaxData)
	if err != nil {
		return nil, err
	}

	return lex, nil
}

func ValidateSyntax(t *testing.T, lex *Lexer) {
	type TokensData struct {
		tokenCount            int
		idents                []string
		cardinalities         []TermCardinality
		separators            []string
		spaceBeforeSeparators []SpaceParsing
		spaceAfterSeparators  []SpaceParsing
	}

	tests := []struct {
		termName    string
		phraseCount int
		tokens      []TokensData
		spacesAfter []SpaceParsing
	}{
		{
			"foo",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"fooAndBar",
			2,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"bar",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
				SpaceAny,
			},
		},
		{
			"fooOrBar",
			1,
			[]TokensData{
				{
					2,
					[]string{
						"foo",
						"bar",
					},
					[]TermCardinality{
						CardinalityOne,
						CardinalityOne,
					},
					[]string{
						"",
						"",
					},
					[]SpaceParsing{
						SpaceAny,
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"fooOptionalBar",
			2,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"bar",
					},
					[]TermCardinality{
						CardinalityOptional,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
				SpaceAny,
			},
		},
		{
			"atLeastOneFoo",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityAtLeastOne,
					},
					[]string{
						",",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"fooManyBar",
			2,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"bar",
					},
					[]TermCardinality{
						CardinalityMany,
					},
					[]string{
						",",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
				SpaceAny,
			},
		},
		{
			"fooTerm",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"$foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"fooIntBar",
			3,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"!int",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"bar",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
				SpaceAny,
				SpaceAny,
			},
		},
		{
			"fooNoSpaceBar",
			2,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"bar",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceNone,
				SpaceAny,
			},
		},
		{
			"fooSpaceBar",
			2,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
				{
					1,
					[]string{
						"bar",
					},
					[]TermCardinality{
						CardinalityOne,
					},
					[]string{
						"",
					},
					[]SpaceParsing{
						SpaceAny,
					},
					[]SpaceParsing{
						SpaceAny,
					},
				},
			},
			[]SpaceParsing{
				SpaceRequired,
				SpaceAny,
			},
		},
		{
			"manyFooBarSeparatorNoSpaceBeforeAfter",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityMany,
					},
					[]string{
						"bar",
					},
					[]SpaceParsing{
						SpaceNone,
					},
					[]SpaceParsing{
						SpaceNone,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"manyFooBarSeparatorSpaceBeforeAfter",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityMany,
					},
					[]string{
						"bar",
					},
					[]SpaceParsing{
						SpaceRequired,
					},
					[]SpaceParsing{
						SpaceRequired,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"manyFooBarSeparatorNoSpaceBeforeSpaceAfter",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityMany,
					},
					[]string{
						"bar",
					},
					[]SpaceParsing{
						SpaceNone,
					},
					[]SpaceParsing{
						SpaceRequired,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
		{
			"manyFooBarSeparatorSpaceBeforeNoSpaceAfter",
			1,
			[]TokensData{
				{
					1,
					[]string{
						"foo",
					},
					[]TermCardinality{
						CardinalityMany,
					},
					[]string{
						"bar",
					},
					[]SpaceParsing{
						SpaceRequired,
					},
					[]SpaceParsing{
						SpaceNone,
					},
				},
			},
			[]SpaceParsing{
				SpaceAny,
			},
		},
	}

	for _, test := range tests {
		term, err := lex.GetTerm(test.termName)
		if err != nil {
			t.Error(err)
		}

		if len(term.Phrases) != test.phraseCount {
			t.Errorf("expected %d phrases in term %s, found %d", len(term.Phrases), test.termName, test.phraseCount)
		}

		for phraseIdx, phrase := range term.Phrases {
			phraseTokens := test.tokens[phraseIdx]
			if phraseTokens.tokenCount != len(phrase.Alternatives) {
				t.Errorf("expected %d alternatives in phrase number %d of term %s, found %d",
					phraseTokens.tokenCount, phraseIdx, test.termName, len(phrase.Alternatives))
			}

			for phraseTokenIdx, phraseToken := range phrase.Alternatives {
				if phraseToken.Ident != phraseTokens.idents[phraseTokenIdx] {
					t.Errorf("expected token %s at token number %d in phrase number %d of term %s, found %s",
						phraseTokens.idents[phraseTokenIdx], phraseTokenIdx, phraseIdx, test.termName, phraseToken.Ident)
				}

				if phraseToken.Cardinality != phraseTokens.cardinalities[phraseTokenIdx] {
					t.Errorf("expected cardinality %v at token number %d in phrase number %d of term %s, found %v",
						phraseTokens.cardinalities[phraseTokenIdx], phraseTokenIdx, phraseIdx, test.termName, phraseToken.Cardinality)
				}

				if phraseToken.Separator != phraseTokens.separators[phraseTokenIdx] {
					t.Errorf("expected separator %q at token number %d in phrase number %d of term %s, found %q",
						phraseTokens.separators[phraseTokenIdx], phraseTokenIdx, phraseIdx, test.termName, phraseToken.Separator)
				}

				if phraseToken.SpaceBeforeSeparator != phraseTokens.spaceBeforeSeparators[phraseTokenIdx] {
					t.Errorf("expected spaceBefore %v at token number %d in phrase number %d of term %s, found %q",
						phraseTokens.spaceBeforeSeparators[phraseTokenIdx], phraseTokenIdx, phraseIdx, test.termName, phraseToken.SpaceBeforeSeparator)
				}

				if phraseToken.SpaceAfterSeparator != phraseTokens.spaceAfterSeparators[phraseTokenIdx] {
					t.Errorf("expected spaceAfter %v at token number %d in phrase number %d of term %s, found %q",
						phraseTokens.spaceAfterSeparators[phraseTokenIdx], phraseTokenIdx, phraseIdx, test.termName, phraseToken.SpaceAfterSeparator)
				}
			}

			if phrase.SpaceAfter != test.spacesAfter[phraseIdx] {
				t.Errorf("expected phrase spaceAfter %v in phrase number %d of term %s, found %v",
					test.spacesAfter[phraseIdx], phraseIdx, test.termName, phrase.SpaceAfter)
			}
		}
	}
}

func TestSyntaxLoading(t *testing.T) {
	lex := NewLexer()
	err := lex.LoadSyntaxXML(invalidSyntaxData)
	if err == nil {
		t.Error("succesfully loaded invalid syntax data. This should have failed")
		return
	}

	err = lex.LoadSyntaxXML(syntaxData)
	if err != nil {
		t.Error(err)
		return
	}

	ValidateSyntax(t, lex)

	lex.ClearSyntax()

	_, err = lex.GetTerm("foo")
	if err == nil {
		t.Error("failed to clear syntax")
		return
	}
}

func RunMatchingTests(t *testing.T, lex *Lexer, tests MatchingTestTable) {
	for _, tt := range tests {
		matchRes := lex.MatchString(tt.matchString, tt.matchTerm)
		if matchRes.matched != tt.expectMatch {
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

		if len(tt.tokens) != len(matchRes.head.matchStrings) {
			t.Errorf("expected tokens do not match for string %s with term %s. Expected %v but got %v", tt.matchString, tt.matchTerm, tt.tokens, matchRes.head.matchStrings)
			continue
		}

		for tokIndex, tok := range matchRes.head.matchStrings {
			if tok != tt.tokens[tokIndex] {
				t.Errorf("expected tokens do not match for string %s with term %s. Expected %v but got %v", tt.matchString, tt.matchTerm, tt.tokens, matchRes.head.matchStrings)
				break
			}
		}
	}
}

func TestLiteralMatching(t *testing.T) {
	lex, err := MakeTestLexer()
	if err != nil {
		t.Error(err)
		return
	}

	tests := MatchingTestTable{
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
	}

	RunMatchingTests(t, lex, tests)
}

func TestTermMatching(t *testing.T) {
	lex, err := MakeTestLexer()
	if err != nil {
		t.Error(err)
		return
	}

	tests := MatchingTestTable{
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
	}

	RunMatchingTests(t, lex, tests)
}

func TestCardinalityMatching(t *testing.T) {
	lex, err := MakeTestLexer()
	if err != nil {
		t.Error(err)
		return
	}

	tests := MatchingTestTable{
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
			matchString: "foobar,",
			expectMatch: false,
			tokens:      nil,
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
	}

	RunMatchingTests(t, lex, tests)
}

func TestBuiltInMatching(t *testing.T) {
	lex, err := MakeTestLexer()
	if err != nil {
		t.Error(err)
		return
	}

	tests := MatchingTestTable{
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
			matchString: "0f",
			expectMatch: false,
			tokens:      nil,
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
			matchTerm:   BuiltInPrefix + "label",
			matchString: "foo~",
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
	}

	RunMatchingTests(t, lex, tests)
}

func TestComplexMatching(t *testing.T) {
	lex, err := MakeTestLexer()
	if err != nil {
		t.Error(err)
		return
	}

	tests := MatchingTestTable{
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

	RunMatchingTests(t, lex, tests)
}

func TestSpacingMatching(t *testing.T) {
	lex, err := MakeTestLexer()
	if err != nil {
		t.Error(err)
		return
	}

	tests := MatchingTestTable{
		// Spacing matching
		{
			matchTerm:   TermPrefix + "fooNoSpaceBar",
			matchString: "foobar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "fooNoSpaceBar",
			matchString: "foo bar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooSpaceBar",
			matchString: "foobar",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "fooSpaceBar",
			matchString: "foo bar",
			expectMatch: true,
			tokens: []string{
				"foo",
				"bar",
			},
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeAfter",
			matchString: "foobarfoo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeAfter",
			matchString: "foo barfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeAfter",
			matchString: "foobar foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeAfter",
			matchString: "foobarfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeAfter",
			matchString: "foo barfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeAfter",
			matchString: "foobar foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeAfter",
			matchString: "foo bar foo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeSpaceAfter",
			matchString: "foobarfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeSpaceAfter",
			matchString: "foo barfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeSpaceAfter",
			matchString: "foobar foo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorNoSpaceBeforeSpaceAfter",
			matchString: "foo bar foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeNoSpaceAfter",
			matchString: "foobarfoo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeNoSpaceAfter",
			matchString: "foo barfoo",
			expectMatch: true,
			tokens: []string{
				"foo",
				"foo",
			},
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeNoSpaceAfter",
			matchString: "foobar foo",
			expectMatch: false,
			tokens:      nil,
		},
		{
			matchTerm:   TermPrefix + "manyFooBarSeparatorSpaceBeforeNoSpaceAfter",
			matchString: "foo bar foo",
			expectMatch: false,
			tokens:      nil,
		},
	}

	RunMatchingTests(t, lex, tests)
}
