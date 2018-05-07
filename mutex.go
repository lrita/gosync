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
	m := &mutex{ch: make(chan struct{}, 1)}
	m.ch <- struct{}{}
	return m
}

type mutex struct {
	ch chan struct{}
}

func (m *mutex) Lock() {
	<-m.ch
}

func (m *mutex) UnLock() {
	select {
	case m.ch <- struct{}{}:
	default:
		if PanicOnBug {
			panic("unlock of unlocked mutex")
		}
	}
}

func (m *mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *mutex) TryLockTimeout(timeout time.Duration) bool {
	tm := time.NewTimer(timeout)
	select {
	case <-m.ch:
		tm.Stop()
		return true
	case <-tm.C:
	}
	return false
}

func (m *mutex) TryLockContext(ctx context.Context) bool {
	select {
	case <-m.ch:
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
	return &mutexGroup{group: make(map[interface{}]*entry)}
}

type mutexGroup struct {
	mu    sync.Mutex
	group map[interface{}]*entry
}

type entry struct {
	ref int
	mu  Mutex
}

func (m *mutexGroup) get(i interface{}, ref int) Mutex {
	m.mu.Lock()
	defer m.mu.Unlock()
	en, ok := m.group[i]
	if !ok {
		if ref > 0 {
			en = &entry{mu: NewMutex()}
			m.group[i] = en
		} else if PanicOnBug {
			panic("unlock of unlocked mutex")
		} else {
			return nil
		}
	}
	en.ref += ref
	return en.mu
}

func (m *mutexGroup) Lock(i interface{}) {
	m.get(i, 1).Lock()
}

func (m *mutexGroup) UnLock(i interface{}) {
	mu := m.get(i, -1)
	if mu != nil {
		mu.UnLock()
	}
}

func (m *mutexGroup) UnLockAndFree(i interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	en, ok := m.group[i]
	if !ok {
		if PanicOnBug {
			panic("unlock of unlocked mutex")
		}
		return
	}
	en.ref--
	if en.ref == 0 {
		delete(m.group, i)
	}
	en.mu.UnLock()
}

func (m *mutexGroup) TryLock(i interface{}) bool {
	locked := m.get(i, 1).TryLock()
	if !locked {
		m.get(i, -1)
	}
	return locked
}

func (m *mutexGroup) TryLockTimeout(i interface{}, timeout time.Duration) bool {
	locked := m.get(i, 1).TryLockTimeout(timeout)
	if !locked {
		m.get(i, -1)
	}
	return locked
}

func (m *mutexGroup) TryLockContext(i interface{}, ctx context.Context) bool {
	locked := m.get(i, 1).TryLockContext(ctx)
	if !locked {
		m.get(i, -1)
	}
	return locked
}
