package pubsub

import (
	"github.com/dapr/components-contrib/pubsub"
	"github.com/layotto/L8-components/pkg/info"
	"strings"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(info.NewRuntimeInfo())
	r.Register(NewFactory("mock", func() pubsub.PubSub {
		return nil
	}),
	)
	if _, err := r.Create("mock"); err != nil {
		t.Fatalf("create mock store failed: %v", err)
	}
	if _, err := r.Create("not exists"); !strings.Contains(err.Error(), "not regsitered") {
		t.Fatalf("create mock store failed: %v", err)
	}
}
