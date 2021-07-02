package lock

type Feature string

type Config struct {
	Metadata map[string]string `json:"metadata"`
}

type Metadata struct {
	Properties map[string]string `json:"properties"`
}

type TryLockRequest struct {
	ResourceId string
	LockOwner  string
	Expire     int32
}

type TryLockResponse struct {
	Success   bool
	LockOwner string
}

type UnlockRequest struct {
	ResourceId string
	LockOwner  string
}

type UnlockResponse struct {
	Status LockStatus
}

type LockStatus int32

const (
	SUCCESS               LockStatus = 0
	LOCK_UNEXIST          LockStatus = 1
	LOCK_BELONG_TO_OTHERS LockStatus = 2
	INTERNAL_ERROR        LockStatus = 3
)
