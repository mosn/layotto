package helloworld

import (
	"testing"

	"github.com/layotto/layotto/pkg/services/hello"
)

func TestHelloWorld(t *testing.T) {
	hs := NewHelloWorld()
	hs.Init(&hello.HelloConfig{
		HelloString: "test",
	})
	resp, _ := hs.Hello(nil)
	if resp.HelloString != "test" {
		t.Fatalf("hello output failed")
	}
}
