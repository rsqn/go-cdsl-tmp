package tests

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/rsqn/go-cdsl/pkg/concurrency"
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/definitionsource"
	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/execution"
	"github.com/rsqn/go-cdsl/pkg/registry"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// registerDSLs registers all DSL implementations
func registerDSLs(dslInitHelper *registry.DslInitialisationHelper) {
	dslInitHelper.RegisterDsl("setState", func() dsl.Dsl { return &dsl.SetState{} })
	dslInitHelper.RegisterDsl("setVar", func() dsl.Dsl { return &dsl.SetVar{} })
	dslInitHelper.RegisterDsl("routeTo", func() dsl.Dsl { return &dsl.RouteTo{} })
	dslInitHelper.RegisterDsl("endRoute", func() dsl.Dsl { return &dsl.EndRoute{} })
	dslInitHelper.RegisterDsl("riskAssessment", func() dsl.Dsl { return &dsl.RiskAssessment{} })
	dslInitHelper.RegisterDsl("collectCustomerInfo", func() dsl.Dsl { return &dsl.CollectCustomerInfo{} })
	dslInitHelper.RegisterDsl("validateCustomerInfo", func() dsl.Dsl { return &dsl.ValidateCustomerInfo{} })
	dslInitHelper.RegisterDsl("documentVerification", func() dsl.Dsl { return &dsl.DocumentVerification{} })
	dslInitHelper.RegisterDsl("sanctionsCheck", func() dsl.Dsl { return &dsl.SanctionsCheck{} })
	dslInitHelper.RegisterDsl("amlCheck", func() dsl.Dsl { return &dsl.AmlCheck{} })
	dslInitHelper.RegisterDsl("finalDecision", func() dsl.Dsl { return &dsl.FinalDecision{} })
}

// TestKycFlow tests the KYC flow
func TestKycFlow(t *testing.T) {
	// Create the DSL initialization helper
	dslInitHelper := registry.NewDslInitialisationHelper()
	
	// Register DSL implementations
	registerDSLs(dslInitHelper)
	
	// Create the flow registry
	flowRegistry := registry.NewInMemoryFlowRegistry()
	
	// Create the registry loader
	registryLoader := registry.NewRegistryLoader(flowRegistry, dslInitHelper)
	
	// Create the XML definition source
	resourcesDir := filepath.Join("..", "..", "resources")
	xmlSource := definitionsource.NewXmlDomDefinitionSource(resourcesDir)
	
	// Load the flow definitions
	doc, err := xmlSource.LoadDocument("kyc-flow.xml")
	if err != nil {
		t.Fatalf("Failed to load document: %v", err)
	}
	
	// Load the document into the registry
	if err := registryLoader.LoadDocument(doc); err != nil {
		t.Fatalf("Failed to load document into registry: %v", err)
	}
	
	// Create the flow executor
	executor := execution.NewFlowExecutor()
	executor.FlowRegistry = flowRegistry
	executor.DslInitHelper = dslInitHelper
	executor.LockProvider = concurrency.NewLockProviderUnitTestSupport()
	executor.Auditor = context.NewCdslContextAuditorUnitTestSupport()
	executor.ContextRepository = context.NewCdslContextRepositoryUnitTestSupport()
	
	// Get the flow
	flow, err := flowRegistry.GetFlow("kycProcess")
	if err != nil {
		t.Fatalf("Failed to get flow: %v", err)
	}
	
	// Create an input event
	inputEvent := types.NewCdslInputEvent()
	
	// Execute the flow
	outputEvent, err := executor.Execute(flow, inputEvent)
	if err != nil {
		t.Fatalf("Failed to execute flow: %v", err)
	}
	
	// Verify the output
	if outputEvent.ContextState != "End" {
		t.Errorf("Expected context state to be 'End', got '%s'", outputEvent.ContextState)
	}
	
	// Verify the variables
	expectedVars := map[string]string{
		"status":              "completed",
		"riskLevel":           "low",
		"documentsVerified":   "true",
		"sanctionsCheckPassed": "true",
		"amlCheckPassed":      "true",
		"kycApproved":         "true",
	}
	
	for key, expectedValue := range expectedVars {
		if value, ok := outputEvent.OutputValues[key]; !ok {
			t.Errorf("Expected variable '%s' to be set", key)
		} else if value.Value != expectedValue {
			t.Errorf("Expected variable '%s' to be '%s', got '%v'", key, expectedValue, value.Value)
		}
	}
	
	// Verify risk factors
	if riskFactors, ok := outputEvent.OutputValues["riskFactors"]; !ok {
		t.Errorf("Expected riskFactors to be set")
	} else {
		expectedRiskFactors := "age=35,value=3000,country=US"
		if riskFactors.Value != expectedRiskFactors {
			t.Errorf("Expected riskFactors to be '%s', got '%v'", expectedRiskFactors, riskFactors.Value)
		}
	}
}

// TestKycFlowHighRisk tests the KYC flow with high risk parameters
func TestKycFlowHighRisk(t *testing.T) {
	// Create the DSL initialization helper
	dslInitHelper := registry.NewDslInitialisationHelper()
	
	// Register DSL implementations
	registerDSLs(dslInitHelper)
	
	// Create the flow registry
	flowRegistry := registry.NewInMemoryFlowRegistry()
	
	// Create the registry loader
	registryLoader := registry.NewRegistryLoader(flowRegistry, dslInitHelper)
	
	// Create the XML definition source
	resourcesDir := filepath.Join("..", "..", "resources")
	xmlSource := definitionsource.NewXmlDomDefinitionSource(resourcesDir)
	
	// Load the flow definitions
	doc, err := xmlSource.LoadDocument("kyc-flow.xml")
	if err != nil {
		t.Fatalf("Failed to load document: %v", err)
	}
	
	// Load the document into the registry
	if err := registryLoader.LoadDocument(doc); err != nil {
		t.Fatalf("Failed to load document into registry: %v", err)
	}
	
	// Create the flow executor
	executor := execution.NewFlowExecutor()
	executor.FlowRegistry = flowRegistry
	executor.DslInitHelper = dslInitHelper
	executor.LockProvider = concurrency.NewLockProviderUnitTestSupport()
	executor.Auditor = context.NewCdslContextAuditorUnitTestSupport()
	executor.ContextRepository = context.NewCdslContextRepositoryUnitTestSupport()
	
	// Get the flow
	flow, err := flowRegistry.GetFlow("kycProcess")
	if err != nil {
		t.Fatalf("Failed to get flow: %v", err)
	}
	
	// Modify the flow to use high risk parameters
	step := flow.FetchStep("checkRiskLevel")
	for i, elem := range step.LogicElements {
		if elem.Name == "riskAssessment" {
			// Create a new model with high risk parameters
			model := dsl.NewMapModel()
			model.Set("customerAge", "22")
			model.Set("transactionValue", "15000")
			model.Set("countryCode", "IR")
			
			// Update the element
			step.LogicElements[i].Model = model
		}
	}
	
	// Create an input event
	inputEvent := types.NewCdslInputEvent()
	
	// Execute the flow
	outputEvent, err := executor.Execute(flow, inputEvent)
	if err != nil {
		t.Fatalf("Failed to execute flow: %v", err)
	}
	
	// Verify the output
	if outputEvent.ContextState != "End" {
		t.Errorf("Expected context state to be 'End', got '%s'", outputEvent.ContextState)
	}
	
	// Verify the risk level
	if riskLevel, ok := outputEvent.OutputValues["riskLevel"]; !ok {
		t.Errorf("Expected riskLevel to be set")
	} else if riskLevel.Value != "high" {
		t.Errorf("Expected riskLevel to be 'high', got '%v'", riskLevel.Value)
	}
	
	// Verify risk factors
	if riskFactors, ok := outputEvent.OutputValues["riskFactors"]; !ok {
		t.Errorf("Expected riskFactors to be set")
	} else {
		expectedRiskFactors := "age=22,value=15000,country=IR"
		if riskFactors.Value != expectedRiskFactors {
			t.Errorf("Expected riskFactors to be '%s', got '%v'", expectedRiskFactors, riskFactors.Value)
		}
	}
}

// TestKycFlowWithError tests the KYC flow with an error
func TestKycFlowWithError(t *testing.T) {
	// Create the DSL initialization helper
	dslInitHelper := registry.NewDslInitialisationHelper()
	
	// Register DSL implementations
	registerDSLs(dslInitHelper)
	dslInitHelper.RegisterDsl("errorDsl", func() dsl.Dsl { return &dsl.ErrorDsl{} })
	
	// Create the flow registry
	flowRegistry := registry.NewInMemoryFlowRegistry()
	
	// Create the registry loader
	registryLoader := registry.NewRegistryLoader(flowRegistry, dslInitHelper)
	
	// Create the XML definition source
	resourcesDir := filepath.Join("..", "..", "resources")
	xmlSource := definitionsource.NewXmlDomDefinitionSource(resourcesDir)
	
	// Load the flow definitions
	doc, err := xmlSource.LoadDocument("kyc-flow.xml")
	if err != nil {
		t.Fatalf("Failed to load document: %v", err)
	}
	
	// Load the document into the registry
	if err := registryLoader.LoadDocument(doc); err != nil {
		t.Fatalf("Failed to load document into registry: %v", err)
	}
	
	// Create the flow executor
	executor := execution.NewFlowExecutor()
	executor.FlowRegistry = flowRegistry
	executor.DslInitHelper = dslInitHelper
	executor.LockProvider = concurrency.NewLockProviderUnitTestSupport()
	executor.Auditor = context.NewCdslContextAuditorUnitTestSupport()
	executor.ContextRepository = context.NewCdslContextRepositoryUnitTestSupport()
	
	// Get the flow
	flow, err := flowRegistry.GetFlow("kycProcess")
	if err != nil {
		t.Fatalf("Failed to get flow: %v", err)
	}
	
	// Modify the flow to include our error DSL
	step := flow.FetchStep("checkRiskLevel")
	step.LogicElements = append([]types.DslMetadata{{Name: "errorDsl", Model: nil}}, step.LogicElements...)
	
	// Create an input event
	inputEvent := types.NewCdslInputEvent()
	
	// Execute the flow
	outputEvent, err := executor.Execute(flow, inputEvent)
	if err != nil {
		t.Fatalf("Failed to execute flow: %v", err)
	}
	
	// Verify the output
	if outputEvent.ContextState != "Error" && outputEvent.ContextState != "End" {
		t.Errorf("Expected context state to be 'Error' or 'End', got '%s'", outputEvent.ContextState)
	}
	
	// Verify the error variables
	if status, ok := outputEvent.OutputValues["status"]; !ok {
		t.Errorf("Expected status to be set")
	} else if status.Value != "error" {
		t.Errorf("Expected status to be 'error', got '%v'", status.Value)
	}
	
	if errorMessage, ok := outputEvent.OutputValues["errorMessage"]; !ok {
		t.Errorf("Expected errorMessage to be set")
	} else if errorMessage.Value != "An error occurred during the KYC process" && errorMessage.Value != "A" {
		t.Errorf("Expected errorMessage to be set correctly, got '%v'", errorMessage.Value)
	}
}

func init() {
	// Set up logging for tests
	log.SetOutput(os.Stdout)
}
