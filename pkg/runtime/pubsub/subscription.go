// REF: https://github.com/dapr/dapr/blob/master/pkg/runtime/pubsub/subscription.go
package pubsub

type Subscription struct {
	PubsubName string            `json:"pubsubname"`
	Topic      string            `json:"topic"`
	Metadata   map[string]string `json:"metadata"`
}
