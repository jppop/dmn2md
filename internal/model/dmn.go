package model

// Definitions est la racine d'un fichier DMN (élément <definitions>).
type Definitions struct {
	Name      string     `xml:"name,attr"`
	Namespace string     `xml:"namespace,attr"`
	Decisions []Decision `xml:"decision"`
}

// Decision représente un élément <decision>.
type Decision struct {
	ID                      string                   `xml:"id,attr"`
	Name                    string                   `xml:"name,attr"`
	InformationRequirements []InformationRequirement `xml:"informationRequirement"`
	DecisionTable           *DecisionTable           `xml:"decisionTable"`
	LiteralExpression       *LiteralExpression       `xml:"literalExpression"`
}

// InformationRequirement représente un élément <informationRequirement>.
type InformationRequirement struct {
	ID               string           `xml:"id,attr"`
	RequiredDecision RequiredDecision `xml:"requiredDecision"`
}

// RequiredDecision représente l'élément <requiredDecision> avec son href.
type RequiredDecision struct {
	Href string `xml:"href,attr"`
}

// LiteralExpression représente un élément <literalExpression>.
type LiteralExpression struct {
	Text string `xml:"text"`
}

// DecisionTable représente un élément <decisionTable>.
type DecisionTable struct {
	ID        string   `xml:"id,attr"`
	HitPolicy string   `xml:"hitPolicy,attr"`
	Inputs    []Input  `xml:"input"`
	Outputs   []Output `xml:"output"`
	Rules     []Rule   `xml:"rule"`
}

// Input représente un élément <input> (colonne d'entrée).
type Input struct {
	ID              string          `xml:"id,attr"`
	Label           string          `xml:"label,attr"`
	InputExpression InputExpression `xml:"inputExpression"`
}

// InputExpression représente l'expression d'une colonne d'entrée.
type InputExpression struct {
	TypeRef string `xml:"typeRef,attr"`
	Text    string `xml:"text"`
}

// Output représente un élément <output> (colonne de sortie).
type Output struct {
	ID      string `xml:"id,attr"`
	Label   string `xml:"label,attr"`
	TypeRef string `xml:"typeRef,attr"`
}

// Rule représente une ligne de la table (<rule>).
type Rule struct {
	ID            string  `xml:"id,attr"`
	InputEntries  []Entry `xml:"inputEntry"`
	OutputEntries []Entry `xml:"outputEntry"`
}

// Entry représente un <inputEntry> ou <outputEntry>.
type Entry struct {
	Text string `xml:"text"`
}
