package types

// DslMetadata contains metadata about a DSL element
type DslMetadata struct {
	Name  string
	Model interface{}
}

// Action represents the action to take after executing a DSL element
type Action string

const (
	// ActionRoute indicates that the flow should route to another step
	ActionRoute Action = "Route"
	// ActionAwait indicates that the flow should wait for an event
	ActionAwait Action = "Await"
	// ActionEnd indicates that the flow should end
	ActionEnd Action = "End"
	// ActionReject indicates that the flow should reject the input
	ActionReject Action = "Reject"
)

// CdslInputEvent represents an input event to a flow
type CdslInputEvent struct {
	ContextID     string
	RequestedStep string
	Payload       map[string]interface{}
}

// NewCdslInputEvent creates a new CdslInputEvent
func NewCdslInputEvent() *CdslInputEvent {
	return &CdslInputEvent{
		Payload: make(map[string]interface{}),
	}
}

// WithContextID sets the context ID for this input event
func (e *CdslInputEvent) WithContextID(id string) *CdslInputEvent {
	e.ContextID = id
	return e
}

// WithRequestedStep sets the requested step for this input event
func (e *CdslInputEvent) WithRequestedStep(step string) *CdslInputEvent {
	e.RequestedStep = step
	return e
}

// CdslOutputEvent represents an output event from a flow
type CdslOutputEvent struct {
	Action    Action
	NextRoute string
	Payload   map[string]interface{}
}

// NewCdslOutputEvent creates a new CdslOutputEvent
func NewCdslOutputEvent() *CdslOutputEvent {
	return &CdslOutputEvent{
		Payload: make(map[string]interface{}),
	}
}

// CdslFlowOutputEvent represents the output of a flow execution
type CdslFlowOutputEvent struct {
	ContextID     string
	ContextState  string
	OutputValues  map[string]*CdslOutputValue
	Action        Action
	NextRoute     string
	Payload       map[string]interface{}
}

// NewCdslFlowOutputEvent creates a new CdslFlowOutputEvent
func NewCdslFlowOutputEvent() *CdslFlowOutputEvent {
	return &CdslFlowOutputEvent{
		OutputValues: make(map[string]*CdslOutputValue),
		Payload:      make(map[string]interface{}),
	}
}

// With initializes a CdslFlowOutputEvent from a CdslOutputEvent
func (e *CdslFlowOutputEvent) With(output *CdslOutputEvent) *CdslFlowOutputEvent {
	if output != nil {
		e.Action = output.Action
		e.NextRoute = output.NextRoute
		e.Payload = output.Payload
	}
	return e
}

// CdslOutputValue represents an output value from a flow
type CdslOutputValue struct {
	Value interface{}
}

// NewCdslOutputValue creates a new CdslOutputValue
func NewCdslOutputValue(value interface{}) *CdslOutputValue {
	return &CdslOutputValue{
		Value: value,
	}
}
