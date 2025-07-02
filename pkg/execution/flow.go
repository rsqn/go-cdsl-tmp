package execution

import (
	"github.com/rsqn/go-cdsl/pkg/model"
)

// This file provides compatibility with the model package
// and re-exports the types for backward compatibility

// Flow is an alias for model.Flow
type Flow = model.Flow

// FlowStep is an alias for model.FlowStep
type FlowStep = model.FlowStep

// NewFlow creates a new Flow
func NewFlow() *Flow {
	return model.NewFlow()
}

// NewFlowStep creates a new FlowStep
func NewFlowStep(id string) *FlowStep {
	return model.NewFlowStep(id)
}
