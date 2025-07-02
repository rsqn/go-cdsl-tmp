package model

import (
	"github.com/rsqn/go-cdsl/pkg/definitionsource"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// Flow represents a flow definition
type Flow struct {
	ID          string
	DefaultStep string
	ErrorStep   string
	Steps       map[string]*FlowStep
}

// NewFlow creates a new Flow
func NewFlow() *Flow {
	return &Flow{
		Steps: make(map[string]*FlowStep),
	}
}

// From initializes a Flow from a FlowDefinition
func (f *Flow) From(def definitionsource.FlowDefinition) *Flow {
	f.ID = def.ID
	f.DefaultStep = def.DefaultStep
	f.ErrorStep = def.ErrorStep
	return f
}

// PutStep adds a step to the flow
func (f *Flow) PutStep(id string, step *FlowStep) {
	f.Steps[id] = step
}

// FetchStep retrieves a step by ID
func (f *Flow) FetchStep(id string) *FlowStep {
	return f.Steps[id]
}

// FlowStep represents a step in a flow
type FlowStep struct {
	ID            string
	LogicElements []types.DslMetadata
	FinalElements []types.DslMetadata
}

// NewFlowStep creates a new FlowStep
func NewFlowStep(id string) *FlowStep {
	return &FlowStep{
		ID:            id,
		LogicElements: make([]types.DslMetadata, 0),
		FinalElements: make([]types.DslMetadata, 0),
	}
}
