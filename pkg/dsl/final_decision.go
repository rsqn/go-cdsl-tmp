package dsl

import (
	"log"
	"strconv"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// FinalDecisionModel represents the model for the FinalDecision DSL
type FinalDecisionModel struct {
	AutoApprove string `json:"autoApprove"`
}

// FinalDecision is a DSL that makes the final KYC decision
type FinalDecision struct {
	DslSupport
}

// Execute implements Dsl
func (d *FinalDecision) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var autoApprove bool
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if aa, ok := m.Get("autoApprove").(string); ok {
			autoApprove = aa == "true"
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if aa, ok := props["autoApprove"].(string); ok {
				autoApprove = aa == "true"
			}
		} else {
			if aa, ok := m["autoApprove"].(string); ok {
				autoApprove = aa == "true"
			}
		}
	}
	
	// Get customer information and check results from context
	customerName := ctx.GetVar("customerName")
	riskLevel := ctx.GetVar("riskLevel")
	infoValid := ctx.GetVar("infoValid") == "true"
	documentsVerified := ctx.GetVar("documentsVerified") == "true"
	sanctionsCheckPassed := ctx.GetVar("sanctionsCheckPassed") == "true"
	amlCheckPassed := ctx.GetVar("amlCheckPassed") == "true"
	
	log.Printf("FinalDecision: Making decision for customer %s with risk level %s", customerName, riskLevel)
	log.Printf("FinalDecision: Info valid: %v, Documents verified: %v, Sanctions check passed: %v, AML check passed: %v", 
		infoValid, documentsVerified, sanctionsCheckPassed, amlCheckPassed)
	
	// Make decision
	approved := false
	
	// Auto-approve if all checks passed and either auto-approve is enabled or risk level is low
	if infoValid && documentsVerified && sanctionsCheckPassed && amlCheckPassed {
		if autoApprove || riskLevel == "low" {
			approved = true
		} else if riskLevel == "medium" {
			// For medium risk, we might want to perform additional checks
			log.Printf("FinalDecision: Medium risk customer %s requires manual review", customerName)
			approved = true // For demonstration purposes, we'll approve medium risk customers
		} else {
			// For high risk, we might want to require manual approval
			log.Printf("FinalDecision: High risk customer %s requires manual approval", customerName)
			approved = true // For demonstration purposes, we'll approve high risk customers
		}
	}
	
	// Store decision in context
	if err := ctx.PutVar("kycApproved", strconv.FormatBool(approved)); err != nil {
		return nil, err
	}
	
	// Store decision reason in context
	var reason string
	if approved {
		reason = "All checks passed"
	} else {
		reason = "One or more checks failed"
	}
	
	if err := ctx.PutVar("kycDecisionReason", reason); err != nil {
		return nil, err
	}
	
	return nil, nil
}
