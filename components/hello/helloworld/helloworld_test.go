package helloworld

import (
	"testing"

	"github.com/layotto/layotto/components/hello"
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
