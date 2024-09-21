package lexer

import (
	"encoding/xml"
	"fmt"
)

type syntaxTermXML struct {
	Name    string         `xml:"name,attr"`
	Phrases []SyntaxPhrase `xml:"Phrase"`
}

type syntaxXML struct {
	Terms []syntaxTermXML `xml:"Term"`
}

func (tc *TermCardinality) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "one":
		*tc = CardinalityOne
		return nil
	case "optional":
		*tc = CardinalityOptional
		return nil
	case "atLeastOne":
		*tc = CardinalityAtLeastOne
		return nil
	case "many":
		*tc = CardinalityMany
		return nil
	}

	return fmt.Errorf("failed to parse string to TermCardinality %s", attr.Value)
}

func (tok *SyntaxToken) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type resultType SyntaxToken
	result := resultType{
		Cardinality: CardinalityOne,
	}
	if err := d.DecodeElement(&result, &start); err != nil {
		return err
	}
	*tok = (SyntaxToken)(result)
	return nil
}

func (lex *Lexer) LoadSyntaxXML(syntaxData []byte) error {
	var resultXML syntaxXML
	err := xml.Unmarshal(syntaxData, &resultXML)
	if err != nil {
		return err
	}

	lex.ClearSyntax()
	for _, term := range resultXML.Terms {
		lex.AddTerm(SyntaxTerm{term.Phrases}, term.Name)
	}

	return nil
}
