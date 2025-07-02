package context

import (
	"sync"

	"github.com/rsqn/go-cdsl/pkg/types"
)

// PostCommitTask represents a task to be executed after a transaction is committed
type PostCommitTask interface {
	RunTask() error
}

// PostStepTask represents a task to be executed after a step is completed
type PostStepTask interface {
	RunTask() error
}

// CdslRuntime represents the runtime environment for a flow execution
type CdslRuntime struct {
	auditor         CdslContextAuditor
	transactionID   string
	postCommitTasks []PostCommitTask
	postStepTasks   []PostStepTask
	outputValues    map[string]*types.CdslOutputValue
	mu              sync.RWMutex
}

// NewCdslRuntime creates a new CdslRuntime
func NewCdslRuntime() *CdslRuntime {
	return &CdslRuntime{
		postCommitTasks: make([]PostCommitTask, 0),
		postStepTasks:   make([]PostStepTask, 0),
		outputValues:    make(map[string]*types.CdslOutputValue),
	}
}

// SetAuditor sets the auditor for this runtime
func (r *CdslRuntime) SetAuditor(auditor CdslContextAuditor) {
	r.auditor = auditor
}

// GetAuditor returns the auditor for this runtime
func (r *CdslRuntime) GetAuditor() CdslContextAuditor {
	return r.auditor
}

// SetTransactionID sets the transaction ID for this runtime
func (r *CdslRuntime) SetTransactionID(id string) {
	r.transactionID = id
}

// GetTransactionID returns the transaction ID for this runtime
func (r *CdslRuntime) GetTransactionID() string {
	return r.transactionID
}

// AddPostCommitTask adds a task to be executed after the transaction is committed
func (r *CdslRuntime) AddPostCommitTask(task PostCommitTask) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.postCommitTasks = append(r.postCommitTasks, task)
}

// GetPostCommitTasks returns the list of post-commit tasks
func (r *CdslRuntime) GetPostCommitTasks() []PostCommitTask {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.postCommitTasks
}

// ClearPostCommitTasks clears the list of post-commit tasks
func (r *CdslRuntime) ClearPostCommitTasks() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.postCommitTasks = make([]PostCommitTask, 0)
}

// AddPostStepTask adds a task to be executed after a step is completed
func (r *CdslRuntime) AddPostStepTask(task PostStepTask) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.postStepTasks = append(r.postStepTasks, task)
}

// GetPostStepTasks returns the list of post-step tasks
func (r *CdslRuntime) GetPostStepTasks() []PostStepTask {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.postStepTasks
}

// ClearPostStepTasks clears the list of post-step tasks
func (r *CdslRuntime) ClearPostStepTasks() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.postStepTasks = make([]PostStepTask, 0)
}

// AddOutputValue adds an output value to the runtime
func (r *CdslRuntime) AddOutputValue(key string, value *types.CdslOutputValue) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.outputValues[key] = value
}

// GetOutputValueMap returns the map of output values
func (r *CdslRuntime) GetOutputValueMap() map[string]*types.CdslOutputValue {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return r.outputValues
}
