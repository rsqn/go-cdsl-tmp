package dsl

import (
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// EndRoute is a DSL that ends the flow
type EndRoute struct {
	DslSupport
}

// Execute implements Dsl
func (d *EndRoute) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	output := types.NewCdslOutputEvent()
	output.Action = types.ActionEnd
	return output, nil
}
