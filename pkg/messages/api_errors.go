package messages

const (
	// PubSub
	ErrPubsubEmpty              = "pubsub name is empty"
	ErrPubsubNotFound           = "pubsub %s not found"
	ErrTopicEmpty               = "topic is empty in pubsub %s"
	ErrPubsubCloudEventsSer     = "error when marshalling cloud event envelope for topic %s pubsub %s: %s"
	ErrPubsubPublishMessage     = "error when publish to topic %s in pubsub %s: %s"
	ErrPubsubCloudEventCreation = "cannot create cloudevent: %s"
	// State
	ErrStateStoresNotConfigured = "state store is not configured"
	ErrStateStoreNotFound       = "state store %s is not found"
	ErrStateGet                 = "fail to get %s from state store %s: %s"
	ErrStateDelete              = "failed deleting state with key %s: %s"
	ErrStateSave                = "failed saving state in state store %s: %s"
	// StateTransaction
	ErrStateStoreNotSupported     = "state store %s doesn't support transaction"
	ErrNotSupportedStateOperation = "operation type %s not supported"
	ErrStateTransaction           = "error while executing state transaction: %s"
	//	Lock
	ErrLockStoresNotConfigured = "lock store is not configured"
	ErrResourceIdEmpty         = "ResourceId is empty in lock store %s"
	ErrLockOwnerEmpty          = "LockOwner is empty in lock store %s"
	ErrExpireNotPositive       = "Expire is not positive in lock store %s"
	ErrLockStoreNotFound       = "lock store %s not found"
)
