# A mutex implementation with TryLock

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
