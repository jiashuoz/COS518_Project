package chord

import (
	"fmt"
	"testing"
)

func TestRPC(t *testing.T) {
	server0 := MakeServer("127.0.0.1:8888")  // id 0
	server1 := MakeServer("127.0.0.1:11190") // id 1
	server3 := MakeServer("127.0.0.1:10000") // id 3

	fmt.Println(server0.GetID())
	fmt.Println(server1.GetID())
	fmt.Println(server3.GetID())

	// TODO: manually populate the finger tables of three servers...

	// Initialize RPC servers on top of chord servers
	rpcServer0, _ := run(server0, "127.0.0.1:8888")
	rpcServer1, _ := run(server1, "127.0.0.1:11190")
	rpcServer3, _ := run(server3, "127.0.0.1:10000")

	fmt.Println(rpcServer0.getAddr())
	fmt.Println(rpcServer1.getAddr())
	fmt.Println(rpcServer3.getAddr())

	// TODO: test RPC functionality
}
