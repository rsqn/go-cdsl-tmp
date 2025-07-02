package dsl

import (
	"log"
	"strconv"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// ValidateCustomerInfoModel represents the model for the ValidateCustomerInfo DSL
type ValidateCustomerInfoModel struct {
	StrictValidation string `json:"strictValidation"`
}

// ValidateCustomerInfo is a DSL that validates customer information
type ValidateCustomerInfo struct {
	DslSupport
}

// Execute implements Dsl
func (d *ValidateCustomerInfo) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var strictValidation bool
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if sv, ok := m.Get("strictValidation").(string); ok {
			strictValidation = sv == "true"
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if sv, ok := props["strictValidation"].(string); ok {
				strictValidation = sv == "true"
			}
		} else {
			if sv, ok := m["strictValidation"].(string); ok {
				strictValidation = sv == "true"
			}
		}
	}
	
	// Get customer information from context
	customerName := ctx.GetVar("customerName")
	customerAge := ctx.GetVar("customerAge")
	transactionValue := ctx.GetVar("transactionValue")
	countryCode := ctx.GetVar("countryCode")
	
	log.Printf("ValidateCustomerInfo: Validating customer %s, age %s, transaction value %s, country code %s", 
		customerName, customerAge, transactionValue, countryCode)
	
	// Perform validation
	valid := true
	var validationErrors []string
	
	// Validate customer name
	if customerName == "" {
		valid = false
		validationErrors = append(validationErrors, "Customer name is required")
	}
	
	// Validate customer age
	age, err := strconv.Atoi(customerAge)
	if err != nil || age <= 0 {
		valid = false
		validationErrors = append(validationErrors, "Invalid customer age")
	}
	
	// Validate transaction value
	value, err := strconv.Atoi(transactionValue)
	if err != nil || value <= 0 {
		valid = false
		validationErrors = append(validationErrors, "Invalid transaction value")
	}
	
	// Validate country code
	if len(countryCode) != 2 && strictValidation {
		valid = false
		validationErrors = append(validationErrors, "Invalid country code")
	}
	
	// Store validation result in context
	if err := ctx.PutVar("infoValid", strconv.FormatBool(valid)); err != nil {
		return nil, err
	}
	
	// Store validation errors in context if any
	if len(validationErrors) > 0 {
		if err := ctx.PutVar("validationErrors", strconv.Itoa(len(validationErrors))); err != nil {
			return nil, err
		}
		
		for i, errMsg := range validationErrors {
			if err := ctx.PutVar("validationError"+strconv.Itoa(i), errMsg); err != nil {
				return nil, err
			}
		}
	}
	
	return nil, nil
}
