package dsl

import (
	"log"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// SetStateModel represents the model for the SetState DSL
type SetStateModel struct {
	Val string `json:"val"`
}

// SetState is a DSL that sets the state of the context
type SetState struct {
	DslSupport
}

// Execute implements Dsl
func (d *SetState) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var val string
	
	// Try to extract the val attribute from different model types
	switch m := model.(type) {
	case *MapModel:
		if v, ok := m.Get("val").(string); ok {
			val = v
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if v, ok := props["val"].(string); ok {
				val = v
			}
		} else if v, ok := m["val"].(string); ok {
			val = v
		}
	}
	
	if val == "" {
		log.Printf("SetState: Val is not a string or not found")
		return nil, nil
	}
	
	// Remove any trailing quotes
	if len(val) >= 2 && val[len(val)-1] == '"' {
		val = val[:len(val)-1]
	}
	
	log.Printf("SetState: Setting state to '%s'", val)
	
	switch val {
	case "Undefined":
		ctx.State = context.StateUndefined
		log.Printf("STATE CHANGE: Context '%s', New State: Undefined", ctx.ID)
	case "Alive":
		ctx.State = context.StateAlive
		log.Printf("STATE CHANGE: Context '%s', New State: Alive", ctx.ID)
	case "Await":
		ctx.State = context.StateAwait
		log.Printf("STATE CHANGE: Context '%s', New State: Await", ctx.ID)
	case "End":
		ctx.State = context.StateEnd
		log.Printf("STATE CHANGE: Context '%s', New State: End", ctx.ID)
	case "Error":
		ctx.State = context.StateError
		log.Printf("STATE CHANGE: Context '%s', New State: Error", ctx.ID)
	}
	
	return nil, nil
}
