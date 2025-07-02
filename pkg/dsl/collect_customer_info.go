package dsl

import (
	"log"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// CollectCustomerInfoModel represents the model for the CollectCustomerInfo DSL
type CollectCustomerInfoModel struct {
	Name            string `json:"name"`
	Age             string `json:"age"`
	TransactionValue string `json:"transactionValue"`
	CountryCode     string `json:"countryCode"`
}

// CollectCustomerInfo is a DSL that collects customer information
type CollectCustomerInfo struct {
	DslSupport
}

// Execute implements Dsl
func (d *CollectCustomerInfo) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var name, age, transactionValue, countryCode string
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if n, ok := m.Get("name").(string); ok {
			name = n
		}
		if a, ok := m.Get("age").(string); ok {
			age = a
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
			if n, ok := props["name"].(string); ok {
				name = n
			}
			if a, ok := props["age"].(string); ok {
				age = a
			}
			if tv, ok := props["transactionValue"].(string); ok {
				transactionValue = tv
			}
			if cc, ok := props["countryCode"].(string); ok {
				countryCode = cc
			}
		} else {
			if n, ok := m["name"].(string); ok {
				name = n
			}
			if a, ok := m["age"].(string); ok {
				age = a
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
	if len(name) >= 2 && name[len(name)-1] == '"' {
		name = name[:len(name)-1]
	}
	if len(age) >= 2 && age[len(age)-1] == '"' {
		age = age[:len(age)-1]
	}
	if len(transactionValue) >= 2 && transactionValue[len(transactionValue)-1] == '"' {
		transactionValue = transactionValue[:len(transactionValue)-1]
	}
	if len(countryCode) >= 2 && countryCode[len(countryCode)-1] == '"' {
		countryCode = countryCode[:len(countryCode)-1]
	}
	
	// Default values if not provided
	if name == "" {
		name = "John Doe"
	}
	if age == "" {
		age = "30"
	}
	if transactionValue == "" {
		transactionValue = "1000"
	}
	if countryCode == "" {
		countryCode = "US"
	}
	
	log.Printf("CollectCustomerInfo: Collecting information for customer %s, age %s, transaction value %s, country code %s", 
		name, age, transactionValue, countryCode)
	
	// Store the customer information in the context
	if err := ctx.PutVar("customerName", name); err != nil {
		return nil, err
	}
	if err := ctx.PutVar("customerAge", age); err != nil {
		return nil, err
	}
	if err := ctx.PutVar("transactionValue", transactionValue); err != nil {
		return nil, err
	}
	if err := ctx.PutVar("countryCode", countryCode); err != nil {
		return nil, err
	}
	
	return nil, nil
}
