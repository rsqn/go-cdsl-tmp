package concurrency

import (
	"fmt"
	"time"
)

// Lock represents a lock on a resource
type Lock struct {
	ID        string
	Resource  string
	Owner     string
	ExpiresAt time.Time
}

// NewLock creates a new Lock
func NewLock(id string, resource string, owner string, duration time.Duration) *Lock {
	return &Lock{
		ID:        id,
		Resource:  resource,
		Owner:     owner,
		ExpiresAt: time.Now().Add(duration),
	}
}

// IsExpired checks if the lock has expired
func (l *Lock) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// LockRejectedException is thrown when a lock cannot be obtained
type LockRejectedException struct {
	Resource string
	Owner    string
	Message  string
}

// Error implements the error interface
func (e *LockRejectedException) Error() string {
	return fmt.Sprintf("Lock rejected for resource %s by owner %s: %s", e.Resource, e.Owner, e.Message)
}

// NewLockRejectedException creates a new LockRejectedException
func NewLockRejectedException(resource string, owner string, message string) *LockRejectedException {
	return &LockRejectedException{
		Resource: resource,
		Owner:    owner,
		Message:  message,
	}
}

// LockProvider is responsible for obtaining and releasing locks
type LockProvider interface {
	// Obtain attempts to obtain a lock
	Obtain(owner string, resource string, duration time.Duration, retries int, retryMaxDuration time.Duration) (*Lock, error)
	
	// Release releases a lock
	Release(lock *Lock) error
}

// LockProviderUnitTestSupport is a simple implementation of LockProvider for unit tests
type LockProviderUnitTestSupport struct {
	locks map[string]*Lock
}

// NewLockProviderUnitTestSupport creates a new LockProviderUnitTestSupport
func NewLockProviderUnitTestSupport() *LockProviderUnitTestSupport {
	return &LockProviderUnitTestSupport{
		locks: make(map[string]*Lock),
	}
}

// Obtain implements LockProvider
func (p *LockProviderUnitTestSupport) Obtain(owner string, resource string, duration time.Duration, retries int, retryMaxDuration time.Duration) (*Lock, error) {
	lock, exists := p.locks[resource]
	if exists && !lock.IsExpired() {
		return nil, NewLockRejectedException(resource, owner, "Resource is already locked")
	}
	
	lock = NewLock(resource, resource, owner, duration)
	p.locks[resource] = lock
	return lock, nil
}

// Release implements LockProvider
func (p *LockProviderUnitTestSupport) Release(lock *Lock) error {
	delete(p.locks, lock.Resource)
	return nil
}
