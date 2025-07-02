package dsl

import (
	"log"
	"strconv"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// AmlCheckModel represents the model for the AmlCheck DSL
type AmlCheckModel struct {
	CheckLevel string `json:"checkLevel"`
}

// AmlCheck is a DSL that performs Anti-Money Laundering checks
type AmlCheck struct {
	DslSupport
}

// Execute implements Dsl
func (d *AmlCheck) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var checkLevel string
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if cl, ok := m.Get("checkLevel").(string); ok {
			checkLevel = cl
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if cl, ok := props["checkLevel"].(string); ok {
				checkLevel = cl
			}
		} else {
			if cl, ok := m["checkLevel"].(string); ok {
				checkLevel = cl
			}
		}
	}
	
	// Remove any trailing quotes
	if len(checkLevel) >= 2 && checkLevel[len(checkLevel)-1] == '"' {
		checkLevel = checkLevel[:len(checkLevel)-1]
	}
	
	// Default value if not provided
	if checkLevel == "" {
		checkLevel = "standard"
	}
	
	// Get customer information from context
	customerName := ctx.GetVar("customerName")
	transactionValue := ctx.GetVar("transactionValue")
	riskLevel := ctx.GetVar("riskLevel")
	
	log.Printf("AmlCheck: Performing AML check for customer %s with transaction value %s and risk level %s", 
		customerName, transactionValue, riskLevel)
	
	// Perform AML check
	passed := true
	
	// For high-risk customers or large transactions, we might want to perform additional checks
	if riskLevel == "high" || checkLevel == "enhanced" {
		log.Printf("AmlCheck: Performing enhanced AML check for high-risk customer %s", customerName)
		// In a real implementation, this would perform additional checks
	}
	
	// Store check result in context
	if err := ctx.PutVar("amlCheckPassed", strconv.FormatBool(passed)); err != nil {
		return nil, err
	}
	
	// Store check level in context
	if err := ctx.PutVar("amlCheckLevel", checkLevel); err != nil {
		return nil, err
	}
	
	return nil, nil
}
