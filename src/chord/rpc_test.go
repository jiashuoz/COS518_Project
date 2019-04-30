package chord

import (
	"bytes"
	"fmt"
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
	rpcServer0, _ := run(server0)
	rpcServer1, _ := run(server1)
	rpcServer3, _ := run(server3)

	fmt.Println(rpcServer0.getAddr())
	fmt.Println(rpcServer1.getAddr())
	fmt.Println(rpcServer3.getAddr())

	// TODO: test RPC functionality
}

func TestGetSuccessor(t *testing.T) {
	server0 := MakeServer("127.0.0.1:8888")  // id 0
	server1 := MakeServer("127.0.0.1:11190") // id 1
	server3 := MakeServer("127.0.0.1:10000") // id 3

	node0 := MakeNode("127.0.0.1:8888")
	node1 := MakeNode("127.0.0.1:11190")
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

	run(server0)
	run(server1)
	run(server3)

	fmt.Println(server0.String(true))
	fmt.Println(server1.String(true))
	fmt.Println(server3.String(true))

	// ba = byte array
	// ba0 := intToByteArray(0)
	// ba1 := intToByteArray(1)
	ba2 := intToByteArray(2)
	// ba3 := intToByteArray(3)
	// ba4 := intToByteArray(4)
	// ba5 := intToByteArray(5)
	// ba6 := intToByteArray(6)
	// ba7 := intToByteArray(7)

	// key0Succ := server0.FindSuccessor(ba0)
	// if !bytes.Equal(key0Succ.ID, intToByteArray(0)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 0", server0.GetID(), key0Succ.ID)
	// }

	// key0Succ = server1.FindSuccessor(ba0)
	// if !bytes.Equal(key0Succ.ID, intToByteArray(0)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 0", server1.GetID(), key0Succ.ID)
	// }

	// key0Succ = server3.FindSuccessor(ba0)
	// if !bytes.Equal(key0Succ.ID, intToByteArray(0)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key0Succ.ID)
	// }

	// key1Succ := server0.FindSuccessor(ba1)
	// if !bytes.Equal(key1Succ.ID, intToByteArray(1)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 1", server0.GetID(), key1Succ.ID)
	// }

	// key1Succ = server1.FindSuccessor(ba1)
	// if !bytes.Equal(key1Succ.ID, intToByteArray(1)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 1", server1.GetID(), key1Succ.ID)
	// }

	// key1Succ = server3.FindSuccessor(ba1)
	// if !bytes.Equal(key1Succ.ID, intToByteArray(1)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 1", server3.GetID(), key1Succ.ID)
	// }

	// key2Succ := server0.FindSuccessor(ba2)
	// if !bytes.Equal(key2Succ.ID, intToByteArray(3)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 3", server0.GetID(), key2Succ.ID)
	// }

	// key2Succ = server1.FindSuccessor(ba2)
	// if !bytes.Equal(key2Succ.ID, intToByteArray(3)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 3", server1.GetID(), key2Succ.ID)
	// }

	key2Succ := server3.FindSuccessor(ba2)
	if !bytes.Equal(key2Succ.ID, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", server3.GetID(), key2Succ.ID)
	}

	// key0Succ = server3.FindSuccessor(ba1)
	// if !bytes.Equal(key0Succ.ID, intToByteArray(1)) {
	// 	t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key1Succ.ID)
	// }
}
