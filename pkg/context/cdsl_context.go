package context

import (
	"errors"
	"log"
	"sync"
)

// State represents the state of a CdslContext
type State string

const (
	// StateUndefined is the initial state
	StateUndefined State = "Undefined"
	// StateAlive indicates the context is active
	StateAlive State = "Alive"
	// StateAwait indicates the context is waiting for an event
	StateAwait State = "Await"
	// StateEnd indicates the context has completed
	StateEnd State = "End"
	// StateError indicates the context has encountered an error
	StateError State = "Error"

	// MaxTransitionsHistory is the maximum number of transitions to keep in history
	MaxTransitionsHistory = 1000
)

// CdslContext represents the execution context for a flow
type CdslContext struct {
	runtime       *CdslRuntime
	ID            string            `json:"id"`
	State         State             `json:"state"`
	CurrentFlow   string            `json:"currentFlow"`
	CurrentStep   string            `json:"currentStep"`
	TransientVars map[string]interface{} `json:"-"`
	Vars          map[string]string `json:"vars"`
	Transitions   []string          `json:"transitions"`
	mu            sync.RWMutex
}

// NewCdslContext creates a new CdslContext
func NewCdslContext() *CdslContext {
	return &CdslContext{
		Vars:          make(map[string]string),
		TransientVars: make(map[string]interface{}),
		State:         StateUndefined,
		Transitions:   make([]string, 0),
	}
}

// PushTransition adds a transition to the history
func (c *CdslContext) PushTransition(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.Transitions = append(c.Transitions, s)
	if len(c.Transitions) > MaxTransitionsHistory {
		c.Transitions = c.Transitions[1:]
	}
}

// SetRuntime sets the runtime for this context
func (c *CdslContext) SetRuntime(runtime *CdslRuntime) {
	c.runtime = runtime
}

// GetRuntime returns the runtime for this context
func (c *CdslContext) GetRuntime() *CdslRuntime {
	return c.runtime
}

// GetVar retrieves a variable from the context
func (c *CdslContext) GetVar(name string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.Vars[name]
}

// PutVar sets a variable in the context
func (c *CdslContext) PutVar(key string, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.runtime == nil {
		return errors.New("CdslRuntime not present")
	}
	
	oldValue := c.Vars[key]
	c.runtime.GetAuditor().SetVar(c, key, value, oldValue)
	
	if oldValue == "" {
		log.Printf("VAR SET: Context '%s', Key '%s', Value '%s'", c.ID, key, value)
	} else if oldValue != value {
		log.Printf("VAR CHANGE: Context '%s', Key '%s', Old Value '%s', New Value '%s'", c.ID, key, oldValue, value)
	}
	
	c.Vars[key] = value
	return nil
}

// PutTransient sets a transient variable in the context
func (c *CdslContext) PutTransient(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	_, exists := c.TransientVars[key]
	if !exists {
		log.Printf("TRANSIENT SET: Context '%s', Key '%s'", c.ID, key)
	} else {
		log.Printf("TRANSIENT CHANGE: Context '%s', Key '%s'", c.ID, key)
	}
	
	c.TransientVars[key] = value
}

// FetchTransient retrieves a transient variable from the context
func (c *CdslContext) FetchTransient(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.TransientVars[key]
}
