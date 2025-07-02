package registry

import (
	"sync"

	"github.com/rsqn/go-cdsl/pkg/model"
)

// FlowRegistry is responsible for storing and retrieving flows
type FlowRegistry interface {
	// RegisterFlow registers a flow
	RegisterFlow(flow *model.Flow) error
	
	// GetFlow retrieves a flow by ID
	GetFlow(id string) (*model.Flow, error)
}

// InMemoryFlowRegistry is an in-memory implementation of FlowRegistry
type InMemoryFlowRegistry struct {
	flows map[string]*model.Flow
	mu    sync.RWMutex
}

// NewInMemoryFlowRegistry creates a new InMemoryFlowRegistry
func NewInMemoryFlowRegistry() *InMemoryFlowRegistry {
	return &InMemoryFlowRegistry{
		flows: make(map[string]*model.Flow),
	}
}

// RegisterFlow implements FlowRegistry
func (r *InMemoryFlowRegistry) RegisterFlow(flow *model.Flow) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.flows[flow.ID] = flow
	return nil
}

// GetFlow implements FlowRegistry
func (r *InMemoryFlowRegistry) GetFlow(id string) (*model.Flow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.flows[id], nil
}
