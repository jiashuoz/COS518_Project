package chord

// chord_node_test.go should test functions in chord_node.go

import (
	"fmt"
	"testing"
)

func TestMakeNode(t *testing.T) {
	node := MakeNode("127.0.0.1")
	fmt.Println(node.String())
	fmt.Println(node.String())
}
