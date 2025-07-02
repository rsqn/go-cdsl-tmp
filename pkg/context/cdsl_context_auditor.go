package context

// CdslContextAuditor is responsible for auditing context operations
type CdslContextAuditor interface {
	// SetVar audits a variable change
	SetVar(ctx *CdslContext, key string, newValue string, oldValue string)
	
	// Transition audits a transition between steps
	Transition(ctx *CdslContext, flowID string, stepID string)
	
	// Execute audits the execution of a DSL element
	Execute(ctx *CdslContext, flowID string, stepID string, dslName string)
	
	// ExecutePostStep audits the execution of a post-step task
	ExecutePostStep(ctx *CdslContext, flowID string, stepID string, task PostStepTask)
	
	// ExecutePostCommit audits the execution of a post-commit task
	ExecutePostCommit(ctx *CdslContext, flowID string, task PostCommitTask)
	
	// Error audits an error
	Error(ctx *CdslContext, flowID string, stepID string, dslName string, err error)
}

// CdslContextAuditorUnitTestSupport is a simple implementation of CdslContextAuditor for unit tests
type CdslContextAuditorUnitTestSupport struct{}

// NewCdslContextAuditorUnitTestSupport creates a new CdslContextAuditorUnitTestSupport
func NewCdslContextAuditorUnitTestSupport() *CdslContextAuditorUnitTestSupport {
	return &CdslContextAuditorUnitTestSupport{}
}

// SetVar implements CdslContextAuditor
func (a *CdslContextAuditorUnitTestSupport) SetVar(ctx *CdslContext, key string, newValue string, oldValue string) {}

// Transition implements CdslContextAuditor
func (a *CdslContextAuditorUnitTestSupport) Transition(ctx *CdslContext, flowID string, stepID string) {}

// Execute implements CdslContextAuditor
func (a *CdslContextAuditorUnitTestSupport) Execute(ctx *CdslContext, flowID string, stepID string, dslName string) {}

// ExecutePostStep implements CdslContextAuditor
func (a *CdslContextAuditorUnitTestSupport) ExecutePostStep(ctx *CdslContext, flowID string, stepID string, task PostStepTask) {}

// ExecutePostCommit implements CdslContextAuditor
func (a *CdslContextAuditorUnitTestSupport) ExecutePostCommit(ctx *CdslContext, flowID string, task PostCommitTask) {}

// Error implements CdslContextAuditor
func (a *CdslContextAuditorUnitTestSupport) Error(ctx *CdslContext, flowID string, stepID string, dslName string, err error) {}
