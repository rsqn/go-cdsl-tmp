package dsl

import (
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// AwaitModel represents the model for the Await DSL
type AwaitModel struct {
	At string `json:"at"`
}

// Await is a DSL that pauses execution and waits for an event
type Await struct {
	DslSupport
}

// Execute implements Dsl
func (d *Await) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	m, ok := model.(*MapModel)
	if !ok {
		return nil, nil
	}
	
	at, ok := m.Get("at").(string)
	if !ok {
		return nil, nil
	}
	
	output := types.NewCdslOutputEvent()
	output.Action = types.ActionAwait
	output.NextRoute = at
	return output, nil
}
