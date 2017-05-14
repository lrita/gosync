# A mutex implementation with TryLock
![gosync.svg](https://travis-ci.org/lrita/gosync.svg?branch=master)

## safe guaranteed
Here is a general situations:

```
mutex.Lock()
A()                       mutex.Lock()
mutex.UnLock()            B()
                          mutex.UnLock()
```

When we synchronize using mutex, we should guarantee unlock() `happens before` lock()
to make no data race.

[_A send on a channel happens before the corresponding receive from that channel completes_](https://golang.org/ref/mem#tmp_7), 
so we use a buffered channel to guarantee this.

## interface
```
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

type MutexGroup interface {
	Lock(i interface{})
	UnLock(i interface{})
	// TryLock return true if it fetch mutex
	TryLock(i interface{}) bool
	// TryLockTimeout return true if it fetch mutex, return false if timeout
	TryLockTimeout(i interface{}, timeout time.Duration) bool
	// TryLockTimeout return true if it fetch mutex, return false if context done
	TryLockContext(i interface{}, ctx context.Context) bool
}
```
