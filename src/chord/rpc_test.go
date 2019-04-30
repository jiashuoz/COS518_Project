package chord

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

func TestRunRPC(t *testing.T) {
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

func TestGetSuccessor(t *testing.T) {
	server0 := MakeServer("127.0.0.1:8888")  // id 0
	server1 := MakeServer("127.0.0.1:11190") // id 1
	server3 := MakeServer("127.0.0.1:10000") // id 3

	node1 := MakeNode("127.0.0.1:11190")
	node0 := MakeNode("127.0.0.1:8888")
	node3 := MakeNode("127.0.0.1:10000")

	server0.fingerTable[0] = node1
	server0.fingerTable[1] = node3
	server0.fingerTable[2] = node0

	server1.fingerTable[0] = node3
	server1.fingerTable[1] = node3
	server1.fingerTable[2] = node0

	server3.fingerTable[0] = node0
	server3.fingerTable[1] = node0
	server3.fingerTable[2] = node0

	run(server0, "127.0.0.1:8888")
	run(server1, "127.0.0.1:11190")
	run(server3, "127.0.0.1:10000")

	fmt.Println(server0.String(true))
	fmt.Println(server1.String(true))
	fmt.Println(server3.String(true))

	test1 := server0.FindSuccessor(big.NewInt(3).Bytes())
	if !bytes.Equal(test1.ID, big.NewInt(3).Bytes()) {
		t.Errorf("Find successor for key 3 = %d; want 3", test1.ID)
	}

	test2 := server0.FindSuccessor(big.NewInt(1).Bytes())
	if !bytes.Equal(test2.ID, big.NewInt(1).Bytes()) {
		t.Errorf("Find successor for key 1 = %d; want 1", test2.ID)
	}

	test3 := server0.FindSuccessor(big.NewInt(5).Bytes())
	if !bytes.Equal(test3.ID, big.NewInt(0).Bytes()) {
		t.Errorf("Find successor for key 5 = %d; want 0", test3.ID)
	}

	test4 := server0.FindSuccessor(big.NewInt(7).Bytes())
	if !bytes.Equal(test4.ID, big.NewInt(0).Bytes()) {
		t.Errorf("Find successor for key 7 = %d; want 0", test4.ID)
	}

	test5 := server0.FindSuccessor(big.NewInt(2).Bytes())
	if !bytes.Equal(test5.ID, big.NewInt(3).Bytes()) {
		t.Errorf("Find successor for key 2 = %d; want 3", test5.ID)
	}

}
