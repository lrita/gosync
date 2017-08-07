// go test -race -bench=.
package sync

import (
	"context"
	"sync"
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

func TestMutexUnlockTwice(t *testing.T) {
	mu := NewMutex()
	mu.Lock()
	defer func() {
		if x := recover(); x != nil {
			if x != "unlock of unlocked mutex" {
				t.Errorf("unexpect panic")
			}
		} else {
			t.Errorf("should panic after unlock twice")
		}
	}()
	mu.UnLock()
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

func TestMutexGroupMutliWaitLock(t *testing.T) {
	var (
		wg sync.WaitGroup
		mu = NewMutexGroup()
		cn = 3
	)

	for i := 0; i < cn; i++ {
		wg.Add(1)
		go func() {
			mu.Lock("h")
			time.Sleep(1e7)
			mu.UnLock("h")
			wg.Done()
		}()
	}
	wg.Wait()

	for i := 0; i < cn; i++ {
		wg.Add(1)
		go func() {
			mu.Lock("g")
			time.Sleep(1e7)
			mu.UnLockAndFree("g")
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestMutexGroupUnLockAndFree(t *testing.T) {
	var (
		wg sync.WaitGroup
		mu = NewMutexGroup()
		mg = mu.(*mutexGroup)
	)

	for j := 1; j < 5; j++ {
		for i := 0; i < j; i++ {
			wg.Add(1)
			go func() {
				mu.Lock("h")
				time.Sleep(1e6)
				mu.UnLockAndFree("h")
				wg.Done()
			}()
		}
		wg.Wait()
		mg.mu.Lock()
		if _, ok := mg.group["h"]; ok {
			t.Error("h mutex exist after UnLockAndFree")
		}
		mg.mu.Unlock()
	}
}

func TestMutexGroupTryLockFailedAndUnLockAndFree(t *testing.T) {
	var (
		wg sync.WaitGroup
		mu = NewMutexGroup()
		mg = mu.(*mutexGroup)
	)

	for j := 1; j < 5; j++ {
		for i := 0; i < j; i++ {
			wg.Add(1)
			go func() {
				if mu.TryLock("h") {
					time.Sleep(1e6)
					mu.UnLockAndFree("h")
				}
				wg.Done()
			}()
		}
		wg.Wait()
		mg.mu.Lock()
		if _, ok := mg.group["h"]; ok {
			t.Error("h mutex exist after UnLockAndFree")
		}
		mg.mu.Unlock()
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
