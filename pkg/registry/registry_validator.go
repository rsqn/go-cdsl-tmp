package registry

import (
	"fmt"

	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/exceptions"
	"github.com/rsqn/go-cdsl/pkg/model"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// RegistryValidator is responsible for validating flows in the registry
type RegistryValidator struct {
	flowRegistry  FlowRegistry
	dslInitHelper *DslInitialisationHelper
}

// NewRegistryValidator creates a new RegistryValidator
func NewRegistryValidator(flowRegistry FlowRegistry, dslInitHelper *DslInitialisationHelper) *RegistryValidator {
	return &RegistryValidator{
		flowRegistry:  flowRegistry,
		dslInitHelper: dslInitHelper,
	}
}

// ValidateFlow validates a flow
func (v *RegistryValidator) ValidateFlow(flow *model.Flow) error {
	// Validate flow has an ID
	if flow.ID == "" {
		return exceptions.NewCdslValidationError("Flow must have an ID", nil)
	}
	
	// Validate flow has a default step
	if flow.DefaultStep == "" {
		return exceptions.NewCdslValidationError(fmt.Sprintf("Flow %s must have a default step", flow.ID), nil)
	}
	
	// Validate default step exists
	if flow.FetchStep(flow.DefaultStep) == nil {
		return exceptions.NewCdslValidationError(
			fmt.Sprintf("Flow %s default step %s does not exist", flow.ID, flow.DefaultStep),
			nil,
		)
	}
	
	// Validate error step exists if specified
	if flow.ErrorStep != "" && flow.FetchStep(flow.ErrorStep) == nil {
		return exceptions.NewCdslValidationError(
			fmt.Sprintf("Flow %s error step %s does not exist", flow.ID, flow.ErrorStep),
			nil,
		)
	}
	
	// Validate steps
	for stepID, step := range flow.Steps {
		// Validate step has an ID
		if step.ID == "" {
			return exceptions.NewCdslValidationError(
				fmt.Sprintf("Step in flow %s must have an ID", flow.ID),
				nil,
			)
		}
		
		// Validate step ID matches key
		if step.ID != stepID {
			return exceptions.NewCdslValidationError(
				fmt.Sprintf("Step ID %s does not match key %s in flow %s", step.ID, stepID, flow.ID),
				nil,
			)
		}
		
		// Validate logic elements
		for _, elemMeta := range step.LogicElements {
			if err := v.validateDslElement(elemMeta); err != nil {
				return exceptions.NewCdslValidationError(
					fmt.Sprintf("Invalid logic element %s in step %s of flow %s", elemMeta.Name, step.ID, flow.ID),
					err,
				)
			}
		}
		
		// Validate final elements
		for _, elemMeta := range step.FinalElements {
			if err := v.validateDslElement(elemMeta); err != nil {
				return exceptions.NewCdslValidationError(
					fmt.Sprintf("Invalid final element %s in step %s of flow %s", elemMeta.Name, step.ID, flow.ID),
					err,
				)
			}
		}
	}
	
	return nil
}

// validateDslElement validates a DSL element
func (v *RegistryValidator) validateDslElement(elemMeta types.DslMetadata) error {
	// Validate element has a name
	if elemMeta.Name == "" {
		return exceptions.NewCdslValidationError("DSL element must have a name", nil)
	}
	
	// Validate element can be resolved
	dslInstance := v.dslInitHelper.Resolve(elemMeta)
	if dslInstance == nil {
		return exceptions.NewCdslValidationError(fmt.Sprintf("DSL %s could not be resolved", elemMeta.Name), nil)
	}
	
	// Validate element if it's a validating DSL
	if validatingDsl, ok := dslInstance.(dsl.ValidatingDsl); ok {
		if err := validatingDsl.Validate(); err != nil {
			return err
		}
	}
	
	return nil
}
