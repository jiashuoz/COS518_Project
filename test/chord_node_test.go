package test

import "testing"
import "chord"
import "fmt"

func TestChordNode(t *testing.T) {
	node := chord.MakeNode("0.0.0.0")
	fmt.Println(node.String())
}
