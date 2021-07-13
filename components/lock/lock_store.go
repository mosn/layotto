package lock

type LockStore interface {
	Init(metadata Metadata) error
	Features() []Feature
	TryLock(req *TryLockRequest) (*TryLockResponse, error)
	Unlock(req *UnlockRequest) (*UnlockResponse, error)
}
