package hello

import (
	"strings"
	"testing"

	"github.com/layotto/layotto/pkg/info"
)

func TestNewRegistry(t *testing.T) {
	r := NewRegistry(info.NewRuntimeInfo())
	r.Register(NewHelloFactory("mock", func() HelloService {
		return nil
	}),
	)
	if _, err := r.Create("mock"); err != nil {
		t.Fatalf("create mock hello failed: %v", err)
	}
	if _, err := r.Create("not exists"); !strings.Contains(err.Error(), "not regsitered") {
		t.Fatalf("create mock hello failed: %v", err)
	}
}
