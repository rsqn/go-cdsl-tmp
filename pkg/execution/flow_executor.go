package execution

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/rsqn/go-cdsl/pkg/concurrency"
	"github.com/rsqn/go-cdsl/pkg/context"
	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/exceptions"
	"github.com/rsqn/go-cdsl/pkg/model"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// FlowRegistry is an interface for retrieving flows
type FlowRegistry interface {
	// GetFlow retrieves a flow by ID
	GetFlow(id string) (*model.Flow, error)
}

// DslInitHelper is an interface for resolving DSL instances
type DslInitHelper interface {
	// Resolve resolves a DSL instance from metadata
	Resolve(metadata types.DslMetadata) dsl.Dsl
}

// FlowExecutor is responsible for executing flows
type FlowExecutor struct {
	FlowRegistry         FlowRegistry
	DslInitHelper        DslInitHelper
	LockProvider         concurrency.LockProvider
	Auditor              context.CdslContextAuditor
	ContextRepository    context.CdslContextRepository
	LockRetries          int
	LockDuration         time.Duration
	LockRetryMaxDuration time.Duration
	MyIdentifier         string
}

// NewFlowExecutor creates a new FlowExecutor
func NewFlowExecutor() *FlowExecutor {
	return &FlowExecutor{
		LockRetries:          3,
		LockDuration:         30 * time.Second,
		LockRetryMaxDuration: 1 * time.Second,
		MyIdentifier:         "<anonymous>",
	}
}

// intersectModel creates a deep copy of the model
func (e *FlowExecutor) intersectModel(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	
	// Debug the model
	e.debugModel(src)
	
	data, err := json.Marshal(src)
	if err != nil {
		return src // fallback to original if marshaling fails
	}
	
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return src // fallback to original if unmarshaling fails
	}
	
	return result
}

// debugModel prints debug information about the model
func (e *FlowExecutor) debugModel(model interface{}) {
	if model == nil {
		log.Printf("DEBUG: Model is nil")
		return
	}
	
	log.Printf("DEBUG: Model type: %T", model)
	
	switch m := model.(type) {
	case *dsl.MapModel:
		log.Printf("DEBUG: MapModel properties: %+v", m.Properties)
	case map[string]interface{}:
		log.Printf("DEBUG: map[string]interface{} contents: %+v", m)
	default:
		// Try to use reflection
		val := reflect.ValueOf(model)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		
		if val.Kind() == reflect.Struct {
			log.Printf("DEBUG: Struct with %d fields", val.NumField())
			for i := 0; i < val.NumField(); i++ {
				field := val.Type().Field(i)
				fieldVal := val.Field(i)
				if fieldVal.CanInterface() {
					log.Printf("DEBUG: Field %s = %v", field.Name, fieldVal.Interface())
				} else {
					log.Printf("DEBUG: Field %s (cannot access value)", field.Name)
				}
			}
		}
	}
}

// obtainOutputs executes a list of DSL elements and returns the first output event
func (e *FlowExecutor) obtainOutputs(
	runtime *context.CdslRuntime,
	ctx *context.CdslContext,
	inputEvent *types.CdslInputEvent,
	flow *model.Flow,
	step *model.FlowStep,
	elements []types.DslMetadata,
) (*types.CdslOutputEvent, error) {
	for _, dslMeta := range elements {
		runtime.GetAuditor().Execute(ctx, flow.ID, step.ID, dslMeta.Name)
		log.Printf("DSL EXECUTE: Flow '%s', Step '%s', Element '%s'", flow.ID, step.ID, dslMeta.Name)
		
		dslInstance := e.DslInitHelper.Resolve(dslMeta)
		if dslInstance == nil {
			return nil, exceptions.NewCdslError(fmt.Sprintf("Failed to resolve DSL %s", dslMeta.Name), nil)
		}
		
		// Build or intersect model
		model := e.intersectModel(dslMeta.Model)
		
		// Execute the step
		output, err := dslInstance.Execute(runtime, ctx, model, inputEvent)
		if err != nil {
			log.Printf("DSL ERROR: Flow '%s', Step '%s', Element '%s': %v", flow.ID, step.ID, dslMeta.Name, err)
			return nil, err
		}
		
		// Handle output if required
		if output != nil {
			log.Printf("DSL OUTPUT: Flow '%s', Step '%s', Element '%s', Action: %s", 
				flow.ID, step.ID, dslMeta.Name, output.Action)
			return output, nil
		}
	}
	
	return nil, nil
}

// Execute executes a flow with the given input event
func (e *FlowExecutor) Execute(flow *model.Flow, inputEvent *types.CdslInputEvent) (*types.CdslFlowOutputEvent, error) {
	if flow == nil {
		return nil, exceptions.NewCdslError("Flow must be provided", nil)
	}
	
	var lock *concurrency.Lock
	var ctx *context.CdslContext
	var runtime *context.CdslRuntime
	
	try := func() (*types.CdslFlowOutputEvent, error) {
		var err error
		
		// Create or load context
		if inputEvent.ContextID == "" {
			// Create a new context
			ctx = context.NewCdslContext()
			ctx.ID = uuid.New().String()
			lock, err = e.LockProvider.Obtain(
				e.MyIdentifier,
				"context/"+ctx.ID,
				e.LockDuration,
				e.LockRetries,
				e.LockRetryMaxDuration,
			)
			if err != nil {
				return nil, err
			}
			
			if err := e.ContextRepository.SaveContext(lock.ID, ctx); err != nil {
				return nil, err
			}
			
			ctx, err = e.ContextRepository.GetContext(lock.ID, ctx.ID)
			if err != nil {
				return nil, err
			}
		} else {
			// Lock and load an existing context
			lock, err = e.LockProvider.Obtain(
				e.MyIdentifier,
				"context/"+inputEvent.ContextID,
				e.LockDuration,
				e.LockRetries,
				e.LockRetryMaxDuration,
			)
			if err != nil {
				return nil, err
			}
			
			ctx, err = e.ContextRepository.GetContext(lock.ID, inputEvent.ContextID)
			if err != nil {
				return nil, err
			}
			
			if ctx.State == context.StateEnd {
				return nil, exceptions.NewCdslError(fmt.Sprintf("State of %s is End", ctx.ID), nil)
			}
		}
		
		// Get or determine current step
		if ctx.CurrentStep == "" {
			ctx.CurrentStep = flow.DefaultStep
		}
		
		runtime = context.NewCdslRuntime()
		runtime.SetAuditor(e.Auditor)
		runtime.SetTransactionID(lock.ID)
		ctx.SetRuntime(runtime)
		
		// Get the step
		var step *model.FlowStep
		nextStep := flow.FetchStep(ctx.CurrentStep)
		var outputEvent *types.CdslFlowOutputEvent
		
		if inputEvent.RequestedStep != "" {
			nextStep = flow.FetchStep(inputEvent.RequestedStep)
			if nextStep == nil {
				return nil, exceptions.NewCdslError(fmt.Sprintf("Requested step %s was not found", inputEvent.RequestedStep), nil)
			}
			ctx.CurrentStep = inputEvent.RequestedStep
		}
		
		for nextStep != nil {
			ctx.CurrentStep = nextStep.ID
			ctx.PushTransition(flow.ID + "/" + nextStep.ID)
			runtime.GetAuditor().Transition(ctx, flow.ID, nextStep.ID)
			
			log.Printf("STEP ENTER: Flow '%s', Step '%s'", flow.ID, nextStep.ID)
			
			step = nextStep
			nextStep = nil
			
			var result *types.CdslOutputEvent
			var err error
			
			// Execute logic elements
			generalOutput, err := e.obtainOutputs(runtime, ctx, inputEvent, flow, step, step.LogicElements)
			if err != nil {
				if flow.ErrorStep != "" {
					nextStep = flow.FetchStep(flow.ErrorStep)
					runtime.GetAuditor().Error(ctx, flow.ID, step.ID, "", err)
					log.Printf("STEP ERROR: Flow '%s', Step '%s': %v", flow.ID, step.ID, err)
					continue
				}
				return nil, err
			}
			
			// Execute final elements
			finalOutput, err := e.obtainOutputs(runtime, ctx, inputEvent, flow, step, step.FinalElements)
			if err != nil {
				if flow.ErrorStep != "" {
					nextStep = flow.FetchStep(flow.ErrorStep)
					runtime.GetAuditor().Error(ctx, flow.ID, step.ID, "", err)
					log.Printf("STEP ERROR: Flow '%s', Step '%s': %v", flow.ID, step.ID, err)
					continue
				}
				return nil, err
			}
			
			// Determine result
			if finalOutput != nil {
				result = finalOutput
			} else if generalOutput != nil {
				result = generalOutput
			}
			
			// Execute post step tasks
			for _, task := range runtime.GetPostStepTasks() {
				func() {
					defer func() {
						if r := recover(); r != nil {
							runtime.GetAuditor().Error(ctx, flow.ID, step.ID, "", fmt.Errorf("panic in post step task: %v", r))
						}
					}()
					
					e.Auditor.ExecutePostStep(ctx, flow.ID, step.ID, task)
					_ = task.RunTask()
				}()
			}
			runtime.ClearPostStepTasks()
			
			if result != nil {
				switch result.Action {
				case types.ActionRoute:
					ctx.CurrentStep = result.NextRoute
					nextStep = flow.FetchStep(result.NextRoute)
					if nextStep == nil {
						return nil, exceptions.NewCdslError(fmt.Sprintf("Invalid Route %s", result.NextRoute), nil)
					}
					log.Printf("STEP EXIT: Flow '%s', Step '%s', Action: Route to '%s'", flow.ID, step.ID, result.NextRoute)
				case types.ActionAwait:
					ctx.State = context.StateAwait
					ctx.CurrentStep = result.NextRoute
					log.Printf("STEP EXIT: Flow '%s', Step '%s', Action: Await at '%s'", flow.ID, step.ID, result.NextRoute)
				case types.ActionEnd:
					ctx.State = context.StateEnd
					log.Printf("STEP EXIT: Flow '%s', Step '%s', Action: End", flow.ID, step.ID)
				case types.ActionReject:
					log.Printf("STEP EXIT: Flow '%s', Step '%s', Action: Reject", flow.ID, step.ID)
					// No special handling for reject
				}
				
				outputEvent = types.NewCdslFlowOutputEvent().With(result)
			} else {
				log.Printf("STEP EXIT: Flow '%s', Step '%s', Action: None", flow.ID, step.ID)
			}
		}
		
		// Save context
		if err := e.ContextRepository.SaveContext(runtime.GetTransactionID(), ctx); err != nil {
			return nil, err
		}
		
		// Release lock
		if err := e.LockProvider.Release(lock); err != nil {
			return nil, err
		}
		lock = nil
		
		// Execute post commit tasks
		for _, task := range runtime.GetPostCommitTasks() {
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Just log the panic, don't fail the whole operation
					}
				}()
				
				e.Auditor.ExecutePostCommit(ctx, flow.ID, task)
				_ = task.RunTask()
			}()
		}
		runtime.ClearPostCommitTasks()
		
		// Output something
		if outputEvent == nil {
			outputEvent = types.NewCdslFlowOutputEvent()
		}
		
		outputEvent.ContextID = ctx.ID
		outputEvent.ContextState = string(ctx.State)
		outputEvent.OutputValues = runtime.GetOutputValueMap()
		
		// Add context variables to output values
		if outputEvent.OutputValues == nil {
			outputEvent.OutputValues = make(map[string]*types.CdslOutputValue)
		}
		
		// Copy all variables from context to output
		for key, value := range ctx.Vars {
			outputEvent.OutputValues[key] = &types.CdslOutputValue{
				Value: value,
			}
		}
		
		runtime = nil
		
		return outputEvent, nil
	}
	
	result, err := try()
	
	// Ensure lock is released
	if lock != nil {
		_ = e.LockProvider.Release(lock)
	}
	
	return result, err
}
