package execution

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/rsqn/go-cdsl/pkg/concurrency"
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/registry"
	"github.com/rsqn/go-cdsl/pkg/types"
)

func TestFlowExecutor_Execute(t *testing.T) {
	// Create the DSL initialization helper
	dslInitHelper := registry.NewDslInitialisationHelper()
	
	// Register DSL implementations
	dslInitHelper.RegisterDsl("setState", func() dsl.Dsl { return &dsl.SetState{} })
	dslInitHelper.RegisterDsl("setVar", func() dsl.Dsl { return &dsl.SetVar{} })
	dslInitHelper.RegisterDsl("routeTo", func() dsl.Dsl { return &dsl.RouteTo{} })
	dslInitHelper.RegisterDsl("endRoute", func() dsl.Dsl { return &dsl.EndRoute{} })
	dslInitHelper.RegisterDsl("sayHello", func() dsl.Dsl { return &dsl.SayHello{} })
	
	// Create a flow
	flow := NewFlow()
	flow.ID = "testFlow"
	flow.DefaultStep = "init"
	flow.ErrorStep = "error"
	
	// Create steps
	initStep := NewFlowStep("init")
	
	// Add logic elements to init step
	setStateModel := dsl.NewMapModel()
	setStateModel.Set("val", "Alive")
	initStep.LogicElements = append(initStep.LogicElements, types.DslMetadata{
		Name:  "setState",
		Model: setStateModel,
	})
	
	sayHelloModel := dsl.NewMapModel()
	sayHelloModel.Set("name", "Test")
	initStep.LogicElements = append(initStep.LogicElements, types.DslMetadata{
		Name:  "sayHello",
		Model: sayHelloModel,
	})
	
	setVarModel := dsl.NewMapModel()
	setVarModel.Set("name", "testVar")
	setVarModel.Set("val", "testValue")
	initStep.LogicElements = append(initStep.LogicElements, types.DslMetadata{
		Name:  "setVar",
		Model: setVarModel,
	})
	
	routeToModel := dsl.NewMapModel()
	routeToModel.Set("target", "end")
	initStep.LogicElements = append(initStep.LogicElements, types.DslMetadata{
		Name:  "routeTo",
		Model: routeToModel,
	})
	
	// Create end step
	endStep := NewFlowStep("end")
	
	// Add logic elements to end step
	endStep.LogicElements = append(endStep.LogicElements, types.DslMetadata{
		Name:  "endRoute",
		Model: dsl.NewMapModel(),
	})
	
	// Add final elements to end step
	setEndStateModel := dsl.NewMapModel()
	setEndStateModel.Set("val", "End")
	endStep.FinalElements = append(endStep.FinalElements, types.DslMetadata{
		Name:  "setState",
		Model: setEndStateModel,
	})
	
	// Add steps to flow
	flow.PutStep("init", initStep)
	flow.PutStep("end", endStep)
	
	// Create the flow registry
	flowRegistry := registry.NewInMemoryFlowRegistry()
	flowRegistry.RegisterFlow(flow)
	
	// Create the flow executor
	executor := NewFlowExecutor()
	executor.FlowRegistry = flowRegistry
	executor.DslInitHelper = dslInitHelper
	executor.LockProvider = concurrency.NewLockProviderUnitTestSupport()
	executor.Auditor = context.NewCdslContextAuditorUnitTestSupport()
	executor.ContextRepository = context.NewCdslContextRepositoryUnitTestSupport()
	
	// Create an input event
	inputEvent := types.NewCdslInputEvent()
	
	// Execute the flow
	outputEvent, err := executor.Execute(flow, inputEvent)
	
	// Assert no error
	assert.NoError(t, err)
	
	// Assert output event
	assert.NotNil(t, outputEvent)
	assert.NotEmpty(t, outputEvent.ContextID)
	assert.Equal(t, "End", outputEvent.ContextState)
	
	// Execute the flow again with the same context
	inputEvent = types.NewCdslInputEvent().WithContextID(outputEvent.ContextID)
	_, err = executor.Execute(flow, inputEvent)
	
	// Assert error because context is in End state
	assert.Error(t, err)
}
