package registry

import (
	"log"
	
	"github.com/rsqn/go-cdsl/pkg/definitionsource"
	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/model"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// RegistryLoader is responsible for loading flow definitions into the registry
type RegistryLoader struct {
	flowRegistry  FlowRegistry
	dslInitHelper *DslInitialisationHelper
}

// NewRegistryLoader creates a new RegistryLoader
func NewRegistryLoader(flowRegistry FlowRegistry, dslInitHelper *DslInitialisationHelper) *RegistryLoader {
	return &RegistryLoader{
		flowRegistry:  flowRegistry,
		dslInitHelper: dslInitHelper,
	}
}

// LoadDocument loads a document into the registry
func (l *RegistryLoader) LoadDocument(doc *definitionsource.DocumentDefinition) error {
	for _, flowDef := range doc.Flows {
		flow := model.NewFlow().From(*flowDef)
		
		// Process steps
		for stepID, stepDef := range flowDef.Steps {
			step := model.NewFlowStep(stepID)
			
			// Process logic elements
			for _, elemDef := range stepDef.Elements {
				meta := types.DslMetadata{
					Name:  elemDef.Name,
					Model: l.buildModel(elemDef),
				}
				step.LogicElements = append(step.LogicElements, meta)
			}
			
			// Process finally elements
			for _, elemDef := range stepDef.Finally {
				meta := types.DslMetadata{
					Name:  elemDef.Name,
					Model: l.buildModel(elemDef),
				}
				step.FinalElements = append(step.FinalElements, meta)
			}
			
			flow.PutStep(stepID, step)
		}
		
		if err := l.flowRegistry.RegisterFlow(flow); err != nil {
			return err
		}
	}
	
	return nil
}

// buildModel builds a model from an element definition
func (l *RegistryLoader) buildModel(elemDef definitionsource.ElementDefinition) interface{} {
	model := dsl.NewMapModel()
	
	// Add attributes
	for k, v := range elemDef.Attributes {
		model.Set(k, v)
		log.Printf("Setting attribute in model: %s = %s", k, v)
	}
	
	// Add elements
	for k, v := range elemDef.Elements {
		model.Set(k, v)
		log.Printf("Setting element in model: %s = %v", k, v)
	}
	
	// Add content if present
	if elemDef.Content != "" {
		model.Set("content", elemDef.Content)
		log.Printf("Setting content in model: %s", elemDef.Content)
	}
	
	return model
}
