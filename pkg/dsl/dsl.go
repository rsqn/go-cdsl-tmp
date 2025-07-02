package dsl

import (
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// Dsl is the interface that all DSL implementations must satisfy
// If a DSL returns an Output, execution will stop at that point and an action will be taken based on the output.
// If you wish to return a value, put it in the context
type Dsl interface {
	Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error)
}

// DslSupport provides common functionality for DSL implementations
type DslSupport struct{}

// ValidatingDsl is a DSL that validates its input
type ValidatingDsl interface {
	Dsl
	Validate() error
}

// MapModel represents a model that can be populated from a map
type MapModel struct {
	Properties map[string]interface{}
}

// NewMapModel creates a new MapModel
func NewMapModel() *MapModel {
	return &MapModel{
		Properties: make(map[string]interface{}),
	}
}

// Get retrieves a property from the model
func (m *MapModel) Get(key string) interface{} {
	return m.Properties[key]
}

// Set sets a property in the model
func (m *MapModel) Set(key string, value interface{}) {
	m.Properties[key] = value
}
