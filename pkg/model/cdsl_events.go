package model

import (
	"encoding/json"
)

// CdslInputEvent represents an input event to a flow
type CdslInputEvent struct {
	ContextID     string          `json:"contextId"`
	RequestedStep string          `json:"requestedStep"`
	Payload       json.RawMessage `json:"payload"`
}

// NewCdslInputEvent creates a new CdslInputEvent
func NewCdslInputEvent() *CdslInputEvent {
	return &CdslInputEvent{}
}

// WithContextID sets the context ID for this event
func (e *CdslInputEvent) WithContextID(contextID string) *CdslInputEvent {
	e.ContextID = contextID
	return e
}

// WithRequestedStep sets the requested step for this event
func (e *CdslInputEvent) WithRequestedStep(step string) *CdslInputEvent {
	e.RequestedStep = step
	return e
}

// WithPayload sets the payload for this event
func (e *CdslInputEvent) WithPayload(payload interface{}) *CdslInputEvent {
	data, _ := json.Marshal(payload)
	e.Payload = data
	return e
}

// GetPayload unmarshals the payload into the provided target
func (e *CdslInputEvent) GetPayload(target interface{}) error {
	return json.Unmarshal(e.Payload, target)
}

// Action represents the action to take after a DSL execution
type Action string

const (
	// ActionRoute indicates that execution should continue at the specified route
	ActionRoute Action = "Route"
	// ActionAwait indicates that execution should pause and wait for an event
	ActionAwait Action = "Await"
	// ActionEnd indicates that execution should end
	ActionEnd Action = "End"
	// ActionReject indicates that the input was rejected
	ActionReject Action = "Reject"
)

// CdslOutputEvent represents the output of a DSL execution
type CdslOutputEvent struct {
	Action    Action `json:"action"`
	NextRoute string `json:"nextRoute"`
	Message   string `json:"message"`
}

// NewCdslOutputEvent creates a new CdslOutputEvent
func NewCdslOutputEvent() *CdslOutputEvent {
	return &CdslOutputEvent{}
}

// WithAction sets the action for this event
func (e *CdslOutputEvent) WithAction(action Action) *CdslOutputEvent {
	e.Action = action
	return e
}

// WithNextRoute sets the next route for this event
func (e *CdslOutputEvent) WithNextRoute(route string) *CdslOutputEvent {
	e.NextRoute = route
	return e
}

// WithMessage sets the message for this event
func (e *CdslOutputEvent) WithMessage(message string) *CdslOutputEvent {
	e.Message = message
	return e
}

// CdslOutputValue represents a value output from a flow
type CdslOutputValue struct {
	Value interface{} `json:"value"`
}

// NewCdslOutputValue creates a new CdslOutputValue
func NewCdslOutputValue(value interface{}) *CdslOutputValue {
	return &CdslOutputValue{
		Value: value,
	}
}

// CdslFlowOutputEvent represents the output of a flow execution
type CdslFlowOutputEvent struct {
	CdslOutputEvent
	ContextID     string                     `json:"contextId"`
	ContextState  string                     `json:"contextState"`
	OutputValues  map[string]*CdslOutputValue `json:"outputValues"`
}

// NewCdslFlowOutputEvent creates a new CdslFlowOutputEvent
func NewCdslFlowOutputEvent() *CdslFlowOutputEvent {
	return &CdslFlowOutputEvent{
		OutputValues: make(map[string]*CdslOutputValue),
	}
}

// With copies the properties from the provided CdslOutputEvent
func (e *CdslFlowOutputEvent) With(output *CdslOutputEvent) *CdslFlowOutputEvent {
	e.Action = output.Action
	e.NextRoute = output.NextRoute
	e.Message = output.Message
	return e
}
