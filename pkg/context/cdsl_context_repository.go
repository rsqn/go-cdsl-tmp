package context

// CdslContextRepository is responsible for storing and retrieving contexts
type CdslContextRepository interface {
	// SaveContext saves a context
	SaveContext(transactionID string, ctx *CdslContext) error
	
	// GetContext retrieves a context
	GetContext(transactionID string, contextID string) (*CdslContext, error)
}

// CdslContextRepositoryUnitTestSupport is a simple implementation of CdslContextRepository for unit tests
type CdslContextRepositoryUnitTestSupport struct {
	contexts map[string]*CdslContext
}

// NewCdslContextRepositoryUnitTestSupport creates a new CdslContextRepositoryUnitTestSupport
func NewCdslContextRepositoryUnitTestSupport() *CdslContextRepositoryUnitTestSupport {
	return &CdslContextRepositoryUnitTestSupport{
		contexts: make(map[string]*CdslContext),
	}
}

// SaveContext implements CdslContextRepository
func (r *CdslContextRepositoryUnitTestSupport) SaveContext(transactionID string, ctx *CdslContext) error {
	r.contexts[ctx.ID] = ctx
	return nil
}

// GetContext implements CdslContextRepository
func (r *CdslContextRepositoryUnitTestSupport) GetContext(transactionID string, contextID string) (*CdslContext, error) {
	return r.contexts[contextID], nil
}
