package lexer

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TermCardinality int64

const (
	CardinalityOne TermCardinality = iota
	CardinalityOptional
	CardinalityAtLeastOne
	CardinalityMany
)

type SpaceParsing int64

const (
	SpaceAny SpaceParsing = iota
	SpaceRequired
	SpaceNone
)

const (
	TermPrefix    = "$"
	BuiltInPrefix = "!"
	LiteralPrefix = "@"
)

type MatchNode struct {
	MatchIdent   string
	MatchStrings []string
	Children     []MatchNode
}

type MatchResult struct {
	Head    *MatchNode
	Matched bool
}

// Ident is either the name of a Term to match, or a literal string to match
type SyntaxToken struct {
	Ident       string          `xml:"ident,attr"`
	Cardinality TermCardinality `xml:"cardinality,attr"`
	// These are only used for the plural "TermCardinality"s
	Separator            string       `xml:"separator,attr"`
	SpaceBeforeSeparator SpaceParsing `xml:"spaceBeforeSeparator,attr"`
	SpaceAfterSeparator  SpaceParsing `xml:"spaceAfterSeparator,attr"`
}

// List of Tokens, at least one of which must match
type SyntaxPhrase struct {
	Alternatives []SyntaxToken `xml:"Token"`
	SpaceAfter   SpaceParsing  `xml:"spaceAfter,attr"`
}

// Whitespace separated phrases
type SyntaxTerm struct {
	Phrases []SyntaxPhrase
}

type Lexer struct {
	registeredTerms map[string]SyntaxTerm
}

func NewLexer() *Lexer {
	return &Lexer{
		make(map[string]SyntaxTerm),
	}
}

func (l *Lexer) MatchString(str string, term string) MatchResult {
	result, _, matched, _ := l.lexString(&str, term, SpaceAny)
	if !matched {
		result = nil
	}

	return MatchResult{result, matched && strings.TrimSpace(str) == ""}
}

func (l *Lexer) AddTerm(term SyntaxTerm, name string) {
	l.registeredTerms[name] = term
}

func (l *Lexer) GetTerm(name string) (SyntaxTerm, error) {
	term, exists := l.registeredTerms[name]
	if !exists {
		return SyntaxTerm{}, errors.New("term does not exist")
	}

	return term, nil
}

func (l *Lexer) ClearSyntax() {
	l.registeredTerms = make(map[string]SyntaxTerm)
}

func (l *Lexer) lexString(str *string, expectedTerm string, spacing SpaceParsing) (*MatchNode, int, bool, bool) {
	// Parse past whitespace
	switch spacing {
	case SpaceRequired:
		if r, _ := utf8.DecodeRuneInString(*str); len(*str) > 0 && !unicode.IsSpace(r) {
			return nil, 0, false, false
		}
		fallthrough
	case SpaceAny:
		for r, _ := utf8.DecodeRuneInString(*str); len(*str) > 0 && unicode.IsSpace(r); {
			*str = (*str)[1:]
			r, _ = utf8.DecodeRuneInString(*str)
		}
	case SpaceNone:
		if r, _ := utf8.DecodeRuneInString(*str); len(*str) > 0 && unicode.IsSpace(r) {
			return nil, 0, false, false
		}
	}

	if strings.HasPrefix(expectedTerm, TermPrefix) {
		return l.matchTerm(str, expectedTerm)
	} else if strings.HasPrefix(expectedTerm, BuiltInPrefix) {
		return l.matchBuiltIn(str, expectedTerm)
	}

	result := MatchNode{
		expectedTerm,
		make([]string, 0),
		nil,
	}

	if strings.HasPrefix(expectedTerm, LiteralPrefix) {
		expectedTerm = expectedTerm[1:]
	}

	if strings.HasPrefix(*str, expectedTerm) {
		result.MatchStrings = append(result.MatchStrings, expectedTerm)
		*str = (*str)[len(expectedTerm):]
	}

	match := len(result.MatchStrings) > 0
	return &result, len(result.MatchStrings), match, match
}

func (l *Lexer) matchTerm(str *string, expectedTerm string) (*MatchNode, int, bool, bool) {
	termName := expectedTerm[1:]
	term, found := l.registeredTerms[termName]
	if !found {
		return nil, 0, false, false
	}

	result := MatchNode{
		expectedTerm,
		nil,
		nil,
	}

	strCopy := strings.Clone(*str)
	partialMatch := false
	lastPhraseSpace := SpaceAny
	for _, phrase := range term.Phrases {
		matchedOne := false
		for _, token := range phrase.Alternatives {
			var matchNode *MatchNode
			matchNode, matchedOne = l.parseToken(str, token, lastPhraseSpace)
			if matchedOne {
				if matchNode != nil {
					result.MatchStrings = append(result.MatchStrings, matchNode.MatchStrings...)
					result.Children = append(result.Children, *matchNode)
				}
				partialMatch = true
				break
			}
		}

		if !matchedOne {
			*str = strCopy
			return &result, 0, false, partialMatch
		}

		lastPhraseSpace = phrase.SpaceAfter
	}

	return &result, 1, true, partialMatch
}

func (l *Lexer) matchBuiltIn(str *string, expectedBuiltIn string) (*MatchNode, int, bool, bool) {
	builtInName := expectedBuiltIn[1:]
	switch builtInName {
	default:
		return nil, 0, false, false
	case "label":
		labelString := ""
		for runeIdx, r := range *str {
			if runeIdx == 0 && !isLabelStartRune(r) {
				break
			} else if !isLabelRune(r) {
				break
			}

			labelString += string(r)
		}

		*str = (*str)[len(labelString):]

		match := len(labelString) > 0
		return &MatchNode{
				expectedBuiltIn,
				[]string{labelString},
				nil,
			},
			1,
			match,
			match
	case "match":
		labelString := *str
		*str = ""
		return &MatchNode{
				expectedBuiltIn,
				[]string{labelString},
				nil,
			},
			1,
			true,
			true
	case "int":
		var intString string
		for runeIdx, r := range *str {
			if runeIdx == 0 {
				if r != '-' && !unicode.IsDigit(r) {
					break
				}
			} else if !unicode.IsDigit(r) {
				break
			}

			intString += string(r)
		}

		if intString == "-" {
			return nil, 0, false, false
		}

		*str = (*str)[len(intString):]

		match := len(intString) > 0
		return &MatchNode{
				expectedBuiltIn,
				[]string{intString},
				nil,
			},
			1,
			match,
			match
	}
}

func (l *Lexer) parseToken(str *string, token SyntaxToken, spacing SpaceParsing) (*MatchNode, bool) {
	switch token.Cardinality {
	default:
		fallthrough
	case CardinalityOne:
		matchNode, numMatched, matches, _ := l.lexString(str, token.Ident, spacing)

		if numMatched != 1 {
			return nil, false
		}

		return &MatchNode{token.Ident, matchNode.MatchStrings, matchNode.Children}, matches
	case CardinalityOptional:
		if *str == "" {
			return nil, true
		}

		matchNode, numMatched, matches, partialMatch := l.lexString(str, token.Ident, spacing)
		if numMatched == 0 {
			return nil, !partialMatch
		}

		return &MatchNode{token.Ident, matchNode.MatchStrings, matchNode.Children}, matches
	case CardinalityAtLeastOne:
		splitString := strings.Split(*str, token.Separator)

		var subStrings []string
		var subNodes []MatchNode
		firstMatch := true
		for stringIdx := range splitString {
			splitSpacing := spacing
			if !firstMatch {
				splitSpacing = token.SpaceAfterSeparator
			}

			matchNode, numMatched, matches, _ := l.lexString(str, token.Ident, splitSpacing)
			if !matches || numMatched == 0 {
				return nil, false
			}

			subStrings = append(subStrings, matchNode.MatchStrings...)
			subNodes = append(subNodes, matchNode.Children...)

			if stringIdx < len(splitString)-1 {
				_, _, matches, _ = l.lexString(str, token.Separator, token.SpaceBeforeSeparator)
				if !matches || numMatched == 0 {
					return nil, false
				}
			}

			firstMatch = false
		}

		return &MatchNode{token.Ident, subStrings, subNodes}, len(subStrings) >= 1
	case CardinalityMany:
		if *str == "" {
			return nil, true
		}

		splitString := strings.Split(*str, token.Separator)

		var subStrings []string
		var subNodes []MatchNode
		firstMatch := true
		for stringIdx := range splitString {
			splitSpacing := spacing
			if !firstMatch {
				splitSpacing = token.SpaceAfterSeparator
			}
			matchNode, numMatched, matches, _ := l.lexString(str, token.Ident, splitSpacing)
			if !matches || numMatched == 0 {
				return nil, false
			}

			subStrings = append(subStrings, matchNode.MatchStrings...)
			subNodes = append(subNodes, matchNode.Children...)

			if stringIdx < len(splitString)-1 {
				_, _, matches, _ = l.lexString(str, token.Separator, token.SpaceBeforeSeparator)
				if !matches || numMatched == 0 {
					return nil, false
				}
			}

			firstMatch = false
		}

		return &MatchNode{token.Ident, subStrings, subNodes}, true
	}
}

func isLabelStartRune(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isLabelRune(r rune) bool {
	return isLabelStartRune(r) || unicode.IsDigit(r) || r == '\''
}
