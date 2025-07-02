package dsl

import (
	"log"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// SetVarModel represents the model for the SetVar DSL
type SetVarModel struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

// SetVar is a DSL that sets a variable in the context
type SetVar struct {
	DslSupport
}

// Execute implements Dsl
func (d *SetVar) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var name, val string
	
	// Try to extract the name and val attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if n, ok := m.Get("name").(string); ok {
			name = n
		}
		if v, ok := m.Get("val").(string); ok {
			val = v
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if n, ok := props["name"].(string); ok {
				name = n
			}
			if v, ok := props["val"].(string); ok {
				val = v
			}
		} else {
			if n, ok := m["name"].(string); ok {
				name = n
			}
			if v, ok := m["val"].(string); ok {
				val = v
			}
		}
	}
	
	if name == "" {
		log.Printf("SetVar: Name is not a string or not found")
		return nil, nil
	}
	
	if val == "" {
		log.Printf("SetVar: Val is not a string or not found")
		return nil, nil
	}
	
	// Remove any trailing quotes
	if len(name) >= 2 && name[len(name)-1] == '"' {
		name = name[:len(name)-1]
	}
	if len(val) >= 2 && val[len(val)-1] == '"' {
		val = val[:len(val)-1]
	}
	
	log.Printf("SetVar: Setting variable '%s' to '%s'", name, val)
	
	if err := ctx.PutVar(name, val); err != nil {
		return nil, err
	}
	
	return nil, nil
}
