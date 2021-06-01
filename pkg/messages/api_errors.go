package messages

const (
	// PubSub
	ErrPubsubEmpty              = "pubsub name is empty"
	ErrPubsubNotFound           = "pubsub %s not found"
	ErrTopicEmpty               = "topic is empty in pubsub %s"
	ErrPubsubCloudEventsSer     = "error when marshalling cloud event envelope for topic %s pubsub %s: %s"
	ErrPubsubPublishMessage     = "error when publish to topic %s in pubsub %s: %s"
	ErrPubsubForbidden          = "topic %s is not allowed for app id %s"
	ErrPubsubCloudEventCreation = "cannot create cloudevent: %s"
)
