package dsl

import (
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// ErrorDsl is a DSL that always returns an error
type ErrorDsl struct {
	DslSupport
}

// Execute implements Dsl
func (d *ErrorDsl) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	return nil, &types.CdslError{Message: "Simulated error"}
}
