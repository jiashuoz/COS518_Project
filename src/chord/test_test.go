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

	fmt.Println(server0.String())
	fmt.Println(server1.String())
	fmt.Println(server3.String())

	fmt.Println(server3.LookUp(big.NewInt(5).Bytes()))

}
