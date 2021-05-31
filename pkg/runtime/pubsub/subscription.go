// REF: https://github.com/dapr/dapr/blob/master/pkg/runtime/pubsub/subscription.go
package pubsub

type Subscription struct {
	PubsubName string            `json:"pubsubname"`
	Topic      string            `json:"topic"`
	Route      string            `json:"route"`
	Metadata   map[string]string `json:"metadata"`
	Scopes     []string          `json:"scopes"`
}
