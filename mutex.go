package sync

import (
	"context"
	"sync"
	"time"
)

// Set the behavier on unlock unlocked mutex
var PanicOnBug = true

// another mutex implementation with TryLock method
type Mutex interface {
	Lock()
	UnLock()
	// TryLock return true if it fetch mutex
	TryLock() bool
	// TryLockTimeout return true if it fetch mutex, return false if timeout
	TryLockTimeout(timeout time.Duration) bool
	// TryLockTimeout return true if it fetch mutex, return false if context done
	TryLockContext(ctx context.Context) bool
}

func NewMutex() Mutex {
	return &mutex{ch: make(chan struct{}, 1)}
}

type mutex struct {
	ch chan struct{}
}

func (m *mutex) Lock() {
	m.ch <- struct{}{}
}

func (m *mutex) UnLock() {
	select {
	case <-m.ch:
	default:
		if PanicOnBug {
			panic("unlock of unlocked mutex")
		}
	}
}

func (m *mutex) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
	}
	return false
}

func (m *mutex) TryLockTimeout(timeout time.Duration) bool {
	select {
	case m.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
	}
	return false
}

func (m *mutex) TryLockContext(ctx context.Context) bool {
	select {
	case m.ch <- struct{}{}:
		return true
	case <-ctx.Done():
	}
	return false
}

type MutexGroup interface {
	Lock(i interface{})
	UnLock(i interface{})
	UnLockAndFree(i interface{})
	// TryLock return true if it fetch mutex
	TryLock(i interface{}) bool
	// TryLockTimeout return true if it fetch mutex, return false if timeout
	TryLockTimeout(i interface{}, timeout time.Duration) bool
	// TryLockTimeout return true if it fetch mutex, return false if context done
	TryLockContext(i interface{}, ctx context.Context) bool
}

func NewMutexGroup() MutexGroup {
	return &mutexGroup{group: make(map[interface{}]Mutex)}
}

type mutexGroup struct {
	mu    sync.Mutex
	group map[interface{}]Mutex
}

func (m *mutexGroup) get(i interface{}) Mutex {
	m.mu.Lock()
	mu, ok := m.group[i]
	if !ok {
		mu = NewMutex()
		m.group[i] = mu
	}
	m.mu.Unlock()
	return mu
}

func (m *mutexGroup) Lock(i interface{}) {
	m.get(i).Lock()
}

func (m *mutexGroup) UnLock(i interface{}) {
	m.get(i).UnLock()
}

func (m *mutexGroup) UnLockAndFree(i interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mu, ok := m.group[i]
	if !ok {
		if PanicOnBug {
			panic("unlock of unlocked mutex")
		}
		return
	}
	delete(m.group, i)
	mu.UnLock()
}

func (m *mutexGroup) TryLock(i interface{}) bool {
	return m.get(i).TryLock()
}

func (m *mutexGroup) TryLockTimeout(i interface{}, timeout time.Duration) bool {
	return m.get(i).TryLockTimeout(timeout)
}

func (m *mutexGroup) TryLockContext(i interface{}, ctx context.Context) bool {
	return m.get(i).TryLockContext(ctx)
}
