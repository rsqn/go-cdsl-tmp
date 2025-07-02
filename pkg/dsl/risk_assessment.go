package dsl

import (
	"log"
	"strconv"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// RiskAssessmentModel represents the model for the RiskAssessment DSL
type RiskAssessmentModel struct {
	CustomerAge      string `json:"customerAge"`
	TransactionValue string `json:"transactionValue"`
	CountryCode      string `json:"countryCode"`
}

// RiskAssessment is a DSL that performs risk assessment
type RiskAssessment struct {
	DslSupport
}

// Execute implements Dsl
func (d *RiskAssessment) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var customerAge, transactionValue, countryCode string
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if ca, ok := m.Get("customerAge").(string); ok {
			customerAge = ca
		}
		if tv, ok := m.Get("transactionValue").(string); ok {
			transactionValue = tv
		}
		if cc, ok := m.Get("countryCode").(string); ok {
			countryCode = cc
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if ca, ok := props["customerAge"].(string); ok {
				customerAge = ca
			}
			if tv, ok := props["transactionValue"].(string); ok {
				transactionValue = tv
			}
			if cc, ok := props["countryCode"].(string); ok {
				countryCode = cc
			}
		} else {
			if ca, ok := m["customerAge"].(string); ok {
				customerAge = ca
			}
			if tv, ok := m["transactionValue"].(string); ok {
				transactionValue = tv
			}
			if cc, ok := m["countryCode"].(string); ok {
				countryCode = cc
			}
		}
	}
	
	// Remove any trailing quotes
	if len(customerAge) >= 2 && customerAge[len(customerAge)-1] == '"' {
		customerAge = customerAge[:len(customerAge)-1]
	}
	if len(transactionValue) >= 2 && transactionValue[len(transactionValue)-1] == '"' {
		transactionValue = transactionValue[:len(transactionValue)-1]
	}
	if len(countryCode) >= 2 && countryCode[len(countryCode)-1] == '"' {
		countryCode = countryCode[:len(countryCode)-1]
	}
	
	// Default values if not provided
	if customerAge == "" {
		customerAge = "30"
	}
	if transactionValue == "" {
		transactionValue = "1000"
	}
	if countryCode == "" {
		countryCode = "US"
	}
	
	log.Printf("RiskAssessment: Assessing risk for customer age %s, transaction value %s, country code %s", 
		customerAge, transactionValue, countryCode)
	
	// Convert values
	age, _ := strconv.Atoi(customerAge)
	value, _ := strconv.Atoi(transactionValue)
	
	// Simple risk assessment logic
	var riskLevel string
	
	// High-risk countries
	highRiskCountries := map[string]bool{
		"AF": true, // Afghanistan
		"IR": true, // Iran
		"KP": true, // North Korea
		"SY": true, // Syria
	}
	
	if highRiskCountries[countryCode] {
		riskLevel = "high"
	} else if value > 10000 {
		riskLevel = "high"
	} else if value > 5000 || age < 25 {
		riskLevel = "medium"
	} else {
		riskLevel = "low"
	}
	
	log.Printf("RiskAssessment: Risk level determined as %s", riskLevel)
	
	// Store the risk level in the context
	if err := ctx.PutVar("riskLevel", riskLevel); err != nil {
		return nil, err
	}
	
	// Store additional risk factors
	if err := ctx.PutVar("riskFactors", "age="+customerAge+",value="+transactionValue+",country="+countryCode); err != nil {
		return nil, err
	}
	
	return nil, nil
}
