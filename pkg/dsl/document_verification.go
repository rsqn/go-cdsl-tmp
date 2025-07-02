package dsl

import (
	"log"
	"strconv"
	
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// DocumentVerificationModel represents the model for the DocumentVerification DSL
type DocumentVerificationModel struct {
	DocumentType string `json:"documentType"`
	DocumentID   string `json:"documentId"`
}

// DocumentVerification is a DSL that verifies customer documents
type DocumentVerification struct {
	DslSupport
}

// Execute implements Dsl
func (d *DocumentVerification) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *types.CdslInputEvent) (*types.CdslOutputEvent, error) {
	var documentType, documentID string
	
	// Try to extract attributes from different model types
	switch m := model.(type) {
	case *MapModel:
		if dt, ok := m.Get("documentType").(string); ok {
			documentType = dt
		}
		if did, ok := m.Get("documentId").(string); ok {
			documentID = did
		}
	case map[string]interface{}:
		// Check if there's a Properties key
		if props, ok := m["Properties"].(map[string]interface{}); ok {
			if dt, ok := props["documentType"].(string); ok {
				documentType = dt
			}
			if did, ok := props["documentId"].(string); ok {
				documentID = did
			}
		} else {
			if dt, ok := m["documentType"].(string); ok {
				documentType = dt
			}
			if did, ok := m["documentId"].(string); ok {
				documentID = did
			}
		}
	}
	
	// Remove any trailing quotes
	if len(documentType) >= 2 && documentType[len(documentType)-1] == '"' {
		documentType = documentType[:len(documentType)-1]
	}
	if len(documentID) >= 2 && documentID[len(documentID)-1] == '"' {
		documentID = documentID[:len(documentID)-1]
	}
	
	// Default values if not provided
	if documentType == "" {
		documentType = "passport"
	}
	if documentID == "" {
		documentID = "123456789"
	}
	
	log.Printf("DocumentVerification: Verifying document type %s with ID %s", documentType, documentID)
	
	// Get customer information from context
	customerName := ctx.GetVar("customerName")
	riskLevel := ctx.GetVar("riskLevel")
	
	// Perform document verification
	verified := true
	
	// For high-risk customers, we might want to perform additional verification
	if riskLevel == "high" {
		log.Printf("DocumentVerification: Performing additional verification for high-risk customer %s", customerName)
		// In a real implementation, this would perform additional verification steps
	}
	
	// Store verification result in context
	if err := ctx.PutVar("documentsVerified", strconv.FormatBool(verified)); err != nil {
		return nil, err
	}
	
	// Store document information in context
	if err := ctx.PutVar("documentType", documentType); err != nil {
		return nil, err
	}
	if err := ctx.PutVar("documentID", documentID); err != nil {
		return nil, err
	}
	
	return nil, nil
}
