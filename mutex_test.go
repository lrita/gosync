// go test -race -bench=.
package sync

import (
	"context"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	mu := NewMutex()
	mu.Lock()
	defer mu.UnLock()
	if mu.TryLock() {
		t.Errorf("cannot fetch mutex !!!")
	}
}

func TestMutexTryLockTimeout(t *testing.T) {
	mu := NewMutex()
	mu.Lock()
	go func() {
		time.Sleep(1 * time.Millisecond)
		mu.UnLock()
	}()
	if mu.TryLockTimeout(500 * time.Microsecond) {
		t.Errorf("cannot fetch mutex in 500us !!!")
	}
	if !mu.TryLockTimeout(5 * time.Millisecond) {
		t.Errorf("should fetch mutex in 5ms !!!")
	}
	mu.UnLock()
}

func TestMutexTryLockContext(t *testing.T) {
	mu := NewMutex()
	ctx, cancel := context.WithCancel(context.Background())
	mu.Lock()
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	if mu.TryLockContext(ctx) {
		t.Errorf("cannot fetch mutex !!!")
	}
}

func BenchmarkMutex(b *testing.B) {
	mu := NewMutex()
	a := 0
	c := 0
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			a++
			mu.UnLock()
			mu.Lock()
			c = a
			mu.UnLock()
		}
	})
	_ = a
	_ = c
}

func TestMutexGroup(t *testing.T) {
	mu := NewMutexGroup()
	mu.Lock("g")
	defer mu.UnLock("g")
	if mu.TryLock("g") {
		t.Errorf("cannot fetch mutex !!!")
	}
}

func TestMutexGroupTryLockTimeout(t *testing.T) {
	mu := NewMutexGroup()
	mu.Lock("g")
	go func() {
		time.Sleep(1 * time.Millisecond)
		mu.UnLock("g")
	}()
	if mu.TryLockTimeout("g", 500*time.Microsecond) {
		t.Errorf("cannot fetch mutex in 500us !!!")
	}
	if !mu.TryLockTimeout("g", 5*time.Millisecond) {
		t.Errorf("should fetch mutex in 5ms !!!")
	}
	mu.UnLock("g")
}

func TestMutexGroupTryLockContext(t *testing.T) {
	mu := NewMutexGroup()
	ctx, cancel := context.WithCancel(context.Background())
	mu.Lock("g")
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	if mu.TryLockContext("g", ctx) {
		t.Errorf("cannot fetch mutex !!!")
	}
}

func BenchmarkMutexGroup(b *testing.B) {
	mu := NewMutexGroup()
	a := 0
	c := 0
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock("g")
			a++
			mu.UnLock("g")
			mu.Lock("g")
			c = a
			mu.UnLock("g")
		}
	})
	_ = a
	_ = c
}
