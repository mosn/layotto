package s3

import (
	"fmt"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	var s map[string]string
	if len(s) == 0 {
		fmt.Println("s lenth is 0")
	}
	var n *string
	a := *n
	fmt.Println(a)
}
