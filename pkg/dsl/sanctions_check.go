package dsl

import (
	"log"
	"strconv"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// SanctionsCheckModel represents the model for the SanctionsCheck DSL
type SanctionsCheckModel struct {
	CheckType string `json:"checkType"`
}

// SanctionsCheck is a DSL that checks customer against sanctions lists
type SanctionsCheck struct {
	DslSupport
}

// Execute implements Dsl
func (d *SanctionsCheck) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var checkType string
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if ct, ok := m.Get("checkType").(string); ok {
			checkType = ct
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if ct, ok := props["checkType"].(string); ok {
				checkType = ct
			}
		} else {
			if ct, ok := m["checkType"].(string); ok {
				checkType = ct
			}
		}
	}
	
	// Remove any trailing quotes
	if len(checkType) >= 2 && checkType[len(checkType)-1] == '"' {
		checkType = checkType[:len(checkType)-1]
	}
	
	// Default value if not provided
	if checkType == "" {
		checkType = "standard"
	}
	
	// Get customer information from context
	customerName := ctx.GetVar("customerName")
	countryCode := ctx.GetVar("countryCode")
	riskLevel := ctx.GetVar("riskLevel")
	
	log.Printf("SanctionsCheck: Checking customer %s from %s with risk level %s against sanctions lists", 
		customerName, countryCode, riskLevel)
	
	// Perform sanctions check
	passed := true
	
	// High-risk countries might require additional checks
	highRiskCountries := map[string]bool{
		"AF": true, // Afghanistan
		"IR": true, // Iran
		"KP": true, // North Korea
		"SY": true, // Syria
	}
	
	if highRiskCountries[countryCode] && checkType != "enhanced" {
		log.Printf("SanctionsCheck: Warning - Customer from high-risk country %s but not using enhanced checks", countryCode)
	}
	
	// For enhanced checks, we might want to perform additional verification
	if checkType == "enhanced" {
		log.Printf("SanctionsCheck: Performing enhanced sanctions check for customer %s", customerName)
		// In a real implementation, this would perform additional checks
	}
	
	// Store check result in context
	if err := ctx.PutVar("sanctionsCheckPassed", strconv.FormatBool(passed)); err != nil {
		return nil, err
	}
	
	// Store check type in context
	if err := ctx.PutVar("sanctionsCheckType", checkType); err != nil {
		return nil, err
	}
	
	return nil, nil
}
