package lexer

import (
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

const (
	TermPrefix    = "$"
	BuiltInPrefix = "!"
	LiteralPrefix = "@"
)

// Ident is either the name of a Term to match, or a literal string to match
type SyntaxToken struct {
	Ident       string          `xml:"ident,attr"`
	Cardinality TermCardinality `xml:"cardinality,attr"`
	// This is only used for the plural "TermCardinality"s
	Separator string `xml:"separator,attr"`
}

// List of Tokens, at least one of which must match
type SyntaxPhrase struct {
	Alternatives []SyntaxToken `xml:"Token"`
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

func (l *Lexer) MatchString(str string, term string) ([]string, bool) {
	result, _, matched, _ := l.lexString(&str, term)
	return result, matched && strings.TrimSpace(str) == ""
}

func (l *Lexer) AddTerm(term SyntaxTerm, name string) {
	l.registeredTerms[name] = term
}

func (l *Lexer) ClearSyntax() {
	l.registeredTerms = make(map[string]SyntaxTerm)
}

func (l *Lexer) lexString(str *string, expectedTerm string) ([]string, int, bool, bool) {
	result := make([]string, 0)

	// Parse past whitespace
	for r, _ := utf8.DecodeRuneInString(*str); len(*str) > 0 && unicode.IsSpace(r); {
		*str = (*str)[1:]
		r, _ = utf8.DecodeRuneInString(*str)
	}

	if strings.HasPrefix(expectedTerm, TermPrefix) {
		expectedTermName := expectedTerm[1:]
		return l.matchTerm(str, expectedTermName)
	} else if strings.HasPrefix(expectedTerm, BuiltInPrefix) {
		builtInName := expectedTerm[1:]
		return l.matchBuiltIn(str, builtInName)
	}

	if strings.HasPrefix(expectedTerm, LiteralPrefix) {
		expectedTerm = expectedTerm[1:]
	}

	if strings.HasPrefix(*str, expectedTerm) {
		result = append(result, expectedTerm)
		*str = (*str)[len(expectedTerm):]
	}

	match := len(result) > 0
	return result, len(result), match, match
}

func (l *Lexer) matchTerm(str *string, termName string) ([]string, int, bool, bool) {
	term, found := l.registeredTerms[termName]
	if !found {
		return nil, 0, false, false
	}

	result := make([]string, 0)
	strCopy := strings.Clone(*str)
	partialMatch := false
	for _, phrase := range term.Phrases {
		matchedOne := false
		for _, token := range phrase.Alternatives {
			var matchedStrings []string
			matchedStrings, matchedOne = l.parseToken(str, token)
			if matchedOne {
				partialMatch = true
				// Return on the first matched alternative
				result = append(result, matchedStrings...)
				break
			}
		}

		if !matchedOne {
			*str = strCopy
			return result, 0, false, partialMatch
		}
	}

	return result, 1, true, partialMatch
}

func (l *Lexer) matchBuiltIn(str *string, builtInName string) ([]string, int, bool, bool) {
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
		return []string{labelString}, 1, match, match
	case "match":
		labelString := *str
		*str = ""
		return []string{labelString}, 1, true, true
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
		return []string{intString}, 1, match, match
	}
}

func (l *Lexer) parseToken(str *string, token SyntaxToken) ([]string, bool) {
	switch token.Cardinality {
	default:
		fallthrough
	case CardinalityOne:
		subStrings, numMatched, matches, _ := l.lexString(str, token.Ident)

		if numMatched != 1 {
			return subStrings, false
		}

		return subStrings, matches
	case CardinalityOptional:
		if *str == "" {
			return []string{}, true
		}

		subStrings, numMatched, matches, partialMatch := l.lexString(str, token.Ident)
		if numMatched == 0 {
			return subStrings, !partialMatch
		}

		return subStrings, matches
	case CardinalityAtLeastOne:
		splitString := strings.Split(*str, token.Separator)

		subStrings := make([]string, 0)
		for stringIdx := range splitString {
			splitSubStrings, numMatched, matches, _ := l.lexString(str, token.Ident)
			if !matches || numMatched == 0 {
				return subStrings, false
			}

			subStrings = append(subStrings, splitSubStrings...)

			if stringIdx < len(splitString)-1 {
				_, _, matches, _ = l.lexString(str, token.Separator)
				if !matches || numMatched == 0 {
					return subStrings, false
				}
			}
		}

		return subStrings, len(subStrings) >= 1
	case CardinalityMany:
		if *str == "" {
			return []string{}, true
		}

		splitString := strings.Split(*str, token.Separator)

		subStrings := make([]string, 0)
		for stringIdx := range splitString {
			splitSubStrings, numMatched, matches, _ := l.lexString(str, token.Ident)
			if !matches || numMatched == 0 {
				return subStrings, false
			}

			subStrings = append(subStrings, splitSubStrings...)

			if stringIdx < len(splitString)-1 {
				_, _, matches, _ = l.lexString(str, token.Separator)
				if !matches || numMatched == 0 {
					return subStrings, false
				}
			}
		}

		return subStrings, true
	}
}

func isLabelStartRune(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isLabelRune(r rune) bool {
	return isLabelStartRune(r) || unicode.IsDigit(r) || r == '\''
}
