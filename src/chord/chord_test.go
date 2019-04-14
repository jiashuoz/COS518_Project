package chord

import (
	"fmt"
	"math/big"
	"testing"
)

func Test1(*testing.T) {
	server0 := MakeServer("2")  // ip == "2", id == 0
	server1 := MakeServer("31") // ip == "31", id == 1
	server3 := MakeServer("3")  // ip == "3", id == 3

	Servers["2"] = server0
	Servers["31"] = server1
	Servers["3"] = server3

	server0.node.fingerTable[0] = &FingerEntry{}
	server0.node.fingerTable[0].start = big.NewInt(1).Bytes()
	server0.node.fingerTable[0].ipAddr = "31"
	server0.node.fingerTable[0].id = intToBytes(1)

	server0.node.fingerTable[1] = &FingerEntry{}
	server0.node.fingerTable[1].start = big.NewInt(2).Bytes()
	server0.node.fingerTable[1].ipAddr = "3"
	server0.node.fingerTable[1].id = intToBytes(3)

	server0.node.fingerTable[2] = &FingerEntry{}
	server0.node.fingerTable[2].start = big.NewInt(4).Bytes()
	server0.node.fingerTable[2].ipAddr = "2"
	server0.node.fingerTable[2].id = intToBytes(0)

	server1.node.fingerTable[0] = &FingerEntry{}
	server1.node.fingerTable[0].start = big.NewInt(2).Bytes()
	server1.node.fingerTable[0].ipAddr = "3"
	server1.node.fingerTable[0].id = intToBytes(3)

	server1.node.fingerTable[1] = &FingerEntry{}
	server1.node.fingerTable[1].start = big.NewInt(3).Bytes()
	server1.node.fingerTable[1].ipAddr = "3"
	server1.node.fingerTable[1].id = intToBytes(3)

	server1.node.fingerTable[2] = &FingerEntry{}
	server1.node.fingerTable[2].start = big.NewInt(5).Bytes()
	server1.node.fingerTable[2].ipAddr = "2"
	server1.node.fingerTable[2].id = intToBytes(0)

	server3.node.fingerTable[0] = &FingerEntry{}
	server3.node.fingerTable[0].start = big.NewInt(4).Bytes()
	server3.node.fingerTable[0].ipAddr = "2"
	server3.node.fingerTable[0].id = intToBytes(0)

	server3.node.fingerTable[1] = &FingerEntry{}
	server3.node.fingerTable[1].start = big.NewInt(5).Bytes()
	server3.node.fingerTable[1].ipAddr = "2"
	server3.node.fingerTable[1].id = intToBytes(0)

	server3.node.fingerTable[2] = &FingerEntry{}
	server3.node.fingerTable[2].start = big.NewInt(7).Bytes()
	server3.node.fingerTable[2].ipAddr = "2"
	server3.node.fingerTable[2].id = intToBytes(0)

	fmt.Println(server0.String())
	fmt.Println(server1.String())
	fmt.Println(server3.String())

	// fmt.Println(server3.LookUp(intToBytes(5)))
	fmt.Println(server3.LookUp(intToBytes(1)))
}

// jace TODO: this test is not very thorough, needs more testing
// this test should test the betweenRightInclusive func in chord_server.go
func TestBetweenRightInclusive(t *testing.T) {
	fmt.Println(betweenRightInclusive(big.NewInt(5).Bytes(), big.NewInt(3).Bytes(), big.NewInt(0).Bytes()))
}

// jace TODO: this test is not complete, needs some testing
// this should test the between func in chord_server.go
func TestBetween(t *testing.T) {

}
