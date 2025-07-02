package registry

import (
	"sync"

	"github.com/rsqn/go-cdsl/pkg/dsl"
	"github.com/rsqn/go-cdsl/pkg/types"
)

// DslInitialisationHelper is responsible for resolving DSL instances
type DslInitialisationHelper struct {
	dslFactories map[string]func() dsl.Dsl
	mu           sync.RWMutex
}

// NewDslInitialisationHelper creates a new DslInitialisationHelper
func NewDslInitialisationHelper() *DslInitialisationHelper {
	return &DslInitialisationHelper{
		dslFactories: make(map[string]func() dsl.Dsl),
	}
}

// RegisterDsl registers a DSL factory
func (h *DslInitialisationHelper) RegisterDsl(name string, factory func() dsl.Dsl) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.dslFactories[name] = factory
}

// Resolve resolves a DSL instance from metadata
func (h *DslInitialisationHelper) Resolve(metadata types.DslMetadata) dsl.Dsl {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	factory, exists := h.dslFactories[metadata.Name]
	if !exists {
		return nil
	}
	
	return factory()
}
