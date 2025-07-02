package dsl

import (
	"fmt"
	"log"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// SayHelloModel represents the model for the SayHello DSL
type SayHelloModel struct {
	Name string `json:"name"`
}

// SayHello is a DSL that prints a greeting
type SayHello struct {
	DslSupport
}

// Execute implements Dsl
func (d *SayHello) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var name string
	
	// Try to extract the name attribute from different model types
	switch m := model.(type) {
	case *MapModel:
		if n, ok := m.Get("name").(string); ok {
			name = n
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if n, ok := props["name"].(string); ok {
				name = n
			}
		} else if n, ok := m["name"].(string); ok {
			name = n
		}
	}
	
	if name == "" {
		name = "World"
	} else {
		// Remove any trailing quotes
		if len(name) >= 2 && name[len(name)-1] == '"' {
			name = name[:len(name)-1]
		}
	}
	
	message := fmt.Sprintf("Hello, %s!", name)
	log.Printf("SayHello: %s", message)
	
	if err := ctx.PutVar("greeting", message); err != nil {
		return nil, err
	}
	
	return nil, nil
}
