package definitionsource

// ElementDefinition represents a DSL element definition
type ElementDefinition struct {
	Name       string                 `xml:",name" json:"name" yaml:"name"`
	Attributes map[string]string      `xml:",attr" json:"attributes" yaml:"attributes"`
	Elements   map[string]interface{} `xml:",any" json:"elements" yaml:"elements"`
	Content    string                 `xml:",chardata" json:"content" yaml:"content"`
}

// StepDefinition represents a step definition
type StepDefinition struct {
	ID       string             `xml:"id,attr" json:"id" yaml:"id"`
	Elements []ElementDefinition `xml:",any" json:"elements" yaml:"elements"`
	Finally  []ElementDefinition `xml:"finally>*" json:"finally" yaml:"finally"`
}

// FlowDefinition represents a flow definition
type FlowDefinition struct {
	ID          string                     `xml:"id,attr" json:"id" yaml:"id"`
	DefaultStep string                     `xml:"defaultStep,attr" json:"defaultStep" yaml:"defaultStep"`
	ErrorStep   string                     `xml:"errorStep,attr" json:"errorStep" yaml:"errorStep"`
	Steps       map[string]*StepDefinition `xml:"-" json:"steps" yaml:"steps"`
	StepsList   []StepDefinition          `xml:"step" json:"-" yaml:"-"`
}

// DocumentDefinition represents a document containing flow definitions
type DocumentDefinition struct {
	Flows map[string]*FlowDefinition `xml:"flow" json:"flows" yaml:"flows"`
}
