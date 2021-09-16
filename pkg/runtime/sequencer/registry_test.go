package sequencer

import (
	"mosn.io/layotto/components/pkg/info"
	"mosn.io/layotto/components/sequencer"
	"strings"
	"testing"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(info.NewRuntimeInfo())
	r.Register(NewFactory("mock", func() sequencer.Store {
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
