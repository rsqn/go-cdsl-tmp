package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rsqn/go-cdsl/pkg/concurrency"
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/definitionsource"
	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/execution"
	"github.com/rsqn/go-cdsl/pkg/registry"
	"github.com/rsqn/go-cdsl/pkg/types"
)

func main() {
	// Create the DSL initialization helper
	dslInitHelper := registry.NewDslInitialisationHelper()
	
	// Register DSL implementations
	dslInitHelper.RegisterDsl("setState", func() dsl.Dsl { return &dsl.SetState{} })
	dslInitHelper.RegisterDsl("setVar", func() dsl.Dsl { return &dsl.SetVar{} })
	dslInitHelper.RegisterDsl("routeTo", func() dsl.Dsl { return &dsl.RouteTo{} })
	dslInitHelper.RegisterDsl("endRoute", func() dsl.Dsl { return &dsl.EndRoute{} })
	dslInitHelper.RegisterDsl("await", func() dsl.Dsl { return &dsl.Await{} })
	dslInitHelper.RegisterDsl("sayHello", func() dsl.Dsl { return &dsl.SayHello{} })
	
	// Create the flow registry
	flowRegistry := registry.NewInMemoryFlowRegistry()
	
	// Create the registry loader
	registryLoader := registry.NewRegistryLoader(flowRegistry, dslInitHelper)
	
	// Create the XML definition source
	xmlSource := definitionsource.NewXmlDomDefinitionSource("./resources")
	
	// Load the flow definitions
	doc, err := xmlSource.LoadDocument("test-flow.xml")
	if err != nil {
		log.Fatalf("Failed to load document: %v", err)
	}
	
	// Load the document into the registry
	if err := registryLoader.LoadDocument(doc); err != nil {
		log.Fatalf("Failed to load document into registry: %v", err)
	}
	
	// Create the flow executor
	executor := execution.NewFlowExecutor()
	executor.FlowRegistry = flowRegistry
	executor.DslInitHelper = dslInitHelper
	executor.LockProvider = concurrency.NewLockProviderUnitTestSupport()
	executor.Auditor = context.NewCdslContextAuditorUnitTestSupport()
	executor.ContextRepository = context.NewCdslContextRepositoryUnitTestSupport()
	
	// Get the flow
	flow, err := flowRegistry.GetFlow("shouldRunHelloWorldAndEndRoute")
	if err != nil {
		log.Fatalf("Failed to get flow: %v", err)
	}
	
	// Create an input event
	inputEvent := types.NewCdslInputEvent()
	
	// Execute the flow
	outputEvent, err := executor.Execute(flow, inputEvent)
	if err != nil {
		log.Fatalf("Failed to execute flow: %v", err)
	}
	
	// Print the output
	fmt.Printf("Flow execution completed with context ID: %s\n", outputEvent.ContextID)
	fmt.Printf("Context state: %s\n", outputEvent.ContextState)
	
	// Execute the flow again with the same context
	inputEvent = types.NewCdslInputEvent().WithContextID(outputEvent.ContextID)
	outputEvent, err = executor.Execute(flow, inputEvent)
	if err != nil {
		fmt.Printf("Error executing flow: %v\n", err)
	} else {
		fmt.Printf("Flow execution completed with context ID: %s\n", outputEvent.ContextID)
		fmt.Printf("Context state: %s\n", outputEvent.ContextState)
	}
}

func init() {
	// Create resources directory if it doesn't exist
	resourcesDir := "./resources"
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		log.Fatalf("Failed to create resources directory: %v", err)
	}
	
	// Create a sample flow definition
	flowXml := `<?xml version="1.0" encoding="utf-8" ?>
<cdsl>
    <flow id="shouldRunHelloWorldAndEndRoute" defaultStep="init" errorStep="error">
        <step id="init">
            <setState val="Alive"/>
            <sayHello name="Go"/>
            <setVar name="myVar" val="myVal"/>
            <routeTo target="end"/>
        </step>

        <step id="end">
            <endRoute/>
            <finally>
                <setState val="End"/>
            </finally>
        </step>
    </flow>
</cdsl>`
	
	// Write the flow definition to a file
	flowPath := filepath.Join(resourcesDir, "test-flow.xml")
	if err := os.WriteFile(flowPath, []byte(flowXml), 0644); err != nil {
		log.Fatalf("Failed to write flow definition: %v", err)
	}
}
