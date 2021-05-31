package helloworld

import (
	"github.com/layotto/layotto/pkg/services/hello"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	hs := NewHelloWorld()
	hs.Init(&hello.HelloConfig{
		HelloString: "Hi",
	})

	req := &hello.HelloRequest{
		Name: "Layotto",
	}

	resp, _ := hs.Hello(req)
	if resp.HelloString != "Hi, Layotto" {
		t.Fatalf("hello output failed")
	}
}
