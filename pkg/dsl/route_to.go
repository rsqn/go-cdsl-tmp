package dsl

import (
	"log"
	"reflect"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// RouteToModel represents the model for the RouteTo DSL
type RouteToModel struct {
	Target string `json:"target"`
}

// RouteTo is a DSL that routes to another step
type RouteTo struct {
	DslSupport
}

// Execute implements Dsl
func (d *RouteTo) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	log.Printf("RouteTo: Model type: %v", reflect.TypeOf(model))
	
	// Try to handle different model types
	var target string
	
	switch m := model.(type) {
	case *MapModel:
		log.Printf("RouteTo: Model is a MapModel with properties: %+v", m.Properties)
		if t, ok := m.Get("target").(string); ok {
			target = t
		}
	case map[string]interface{}:
		log.Printf("RouteTo: Model is a map[string]interface{} with keys: %v", m)
		
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			log.Printf("RouteTo: Found Properties map: %v", props)
			if t, ok := props["target"].(string); ok {
				target = t
			}
		} else if t, ok := m["target"].(string); ok {
			target = t
		}
	default:
		// Try to access as a struct with reflection
		log.Printf("RouteTo: Model is of unknown type, trying reflection")
		val := reflect.ValueOf(model)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		
		if val.Kind() == reflect.Struct {
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i)
				if field.Name == "Properties" {
					propVal := val.Field(i)
					if propVal.Kind() == reflect.Map {
						for _, key := range propVal.MapKeys() {
							log.Printf("RouteTo: Found property: %v", key.String())
							if key.String() == "target" {
								targetVal := propVal.MapIndex(key)
								if targetVal.IsValid() && targetVal.CanInterface() {
									if t, ok := targetVal.Interface().(string); ok {
										target = t
									}
								}
							}
						}
					}
				}
			}
		}
	}
	
	if target == "" {
		log.Printf("RouteTo: Target is not a string or not found")
		return nil, nil
	}
	
	// Remove any trailing quotes
	if len(target) >= 2 && target[len(target)-1] == '"' {
		target = target[:len(target)-1]
	}
	
	log.Printf("RouteTo: Routing to target '%s'", target)
	
	output := types.NewCdslOutputEvent()
	output.Action = types.ActionRoute
	output.NextRoute = target
	return output, nil
}
