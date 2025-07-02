# go-cdsl

Provides a simple mechanism for creating a DSL to support your business logic, and the framework to run it in.

## Overview

go-cdsl is a Go implementation of a state machine that allows you to write business logic as a domain-specific language (DSL). It provides a framework for defining, loading, and executing flows composed of DSL elements.

## Features

- Define flows using XML
- Create custom DSL elements
- Execute flows with context persistence
- Support for concurrency with locking
- Auditing of flow execution

## Usage

### Define a Flow

Flows are defined in XML:

```xml
<?xml version="1.0" encoding="utf-8" ?>
<cdsl>
    <flow id="myFlow" defaultStep="init" errorStep="error">
        <step id="init">
            <setState val="Alive"/>
            <sayHello name="World"/>
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
</cdsl>
```

### Create a Custom DSL Element

```go
package mydsl

import (
    "github.com/rsqn/go-cdsl/pkg/context"
    "github.com/rsqn/go-cdsl/pkg/dsl"
    "github.com/rsqn/go-cdsl/pkg/model"
)

// MyCustomDsl is a custom DSL element
type MyCustomDsl struct {
    dsl.DslSupport
}

// Execute implements dsl.Dsl
func (d *MyCustomDsl) Execute(runtime *context.CdslRuntime, ctx *context.CdslContext, model interface{}, input *model.CdslInputEvent) (*model.CdslOutputEvent, error) {
    m, ok := model.(*dsl.MapModel)
    if !ok {
        return nil, nil
    }
    
    // Do something with the model and context
    
    return nil, nil
}
```

### Execute a Flow

```go
package main

import (
    "github.com/rsqn/go-cdsl/pkg/concurrency"
    "github.com/rsqn/go-cdsl/pkg/context"
    "github.com/rsqn/go-cdsl/pkg/definitionsource"
    "github.com/rsqn/go-cdsl/pkg/dsl"
    "github.com/rsqn/go-cdsl/pkg/execution"
    "github.com/rsqn/go-cdsl/pkg/model"
    "github.com/rsqn/go-cdsl/pkg/registry"
)

func main() {
    // Create the DSL initialization helper
    dslInitHelper := registry.NewDslInitialisationHelper()
    
    // Register DSL implementations
    dslInitHelper.RegisterDsl("setState", func() dsl.Dsl { return &dsl.SetState{} })
    dslInitHelper.RegisterDsl("setVar", func() dsl.Dsl { return &dsl.SetVar{} })
    dslInitHelper.RegisterDsl("routeTo", func() dsl.Dsl { return &dsl.RouteTo{} })
    dslInitHelper.RegisterDsl("endRoute", func() dsl.Dsl { return &dsl.EndRoute{} })
    
    // Create the flow registry
    flowRegistry := registry.NewInMemoryFlowRegistry()
    
    // Create the registry loader
    registryLoader := registry.NewRegistryLoader(flowRegistry, dslInitHelper)
    
    // Create the XML definition source
    xmlSource := definitionsource.NewXmlDomDefinitionSource("./resources")
    
    // Load the flow definitions
    doc, err := xmlSource.LoadDocument("flows.xml")
    if err != nil {
        panic(err)
    }
    
    // Load the document into the registry
    if err := registryLoader.LoadDocument(doc); err != nil {
        panic(err)
    }
    
    // Create the flow executor
    executor := execution.NewFlowExecutor()
    executor.FlowRegistry = flowRegistry
    executor.DslInitHelper = dslInitHelper
    executor.LockProvider = concurrency.NewLockProviderUnitTestSupport()
    executor.Auditor = context.NewCdslContextAuditorUnitTestSupport()
    executor.ContextRepository = context.NewCdslContextRepositoryUnitTestSupport()
    
    // Get the flow
    flow, err := flowRegistry.GetFlow("myFlow")
    if err != nil {
        panic(err)
    }
    
    // Create an input event
    inputEvent := model.NewCdslInputEvent()
    
    // Execute the flow
    outputEvent, err := executor.Execute(flow, inputEvent)
    if err != nil {
        panic(err)
    }
    
    // Use the output
    println("Flow execution completed with context ID:", outputEvent.ContextID)
}
```

## Dependencies

This project has minimal dependencies:

- Standard Go libraries
- github.com/google/uuid for generating unique identifiers

## License

This project is licensed under the GNU General Public License, Version 3.0.
