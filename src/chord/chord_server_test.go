package chord

// // chord_server_test.go should test functions in chord_server.go

// import (
// 	"fmt"
// 	// "log"
// 	"math/big"
// 	// "net/rpc"
// 	"testing"
// )

// // import (
// // 	"net"
// // 	"net/http"
// // )

// func Test(t *testing.T) {
// 	fmt.Println("Testing make server and finger table:")
// 	server0 := MakeServer("2")  // ip == "2", id == 0
// 	server1 := MakeServer("31") // ip == "31", id == 1
// 	server3 := MakeServer("3")  // ip == "3", id == 3

// 	Servers["2"] = server0
// 	Servers["31"] = server1
// 	Servers["3"] = server3

// 	server0.node.fingerTable[0] = &FingerEntry{}
// 	server0.node.fingerTable[0].start = big.NewInt(1).Bytes()
// 	server0.node.fingerTable[0].ipAddr = "31"
// 	server0.node.fingerTable[0].id = intToBytes(1)

// 	server0.node.fingerTable[1] = &FingerEntry{}
// 	server0.node.fingerTable[1].start = big.NewInt(2).Bytes()
// 	server0.node.fingerTable[1].ipAddr = "3"
// 	server0.node.fingerTable[1].id = intToBytes(3)

// 	server0.node.fingerTable[2] = &FingerEntry{}
// 	server0.node.fingerTable[2].start = big.NewInt(4).Bytes()
// 	server0.node.fingerTable[2].ipAddr = "2"
// 	server0.node.fingerTable[2].id = intToBytes(0)

// 	server1.node.fingerTable[0] = &FingerEntry{}
// 	server1.node.fingerTable[0].start = big.NewInt(2).Bytes()
// 	server1.node.fingerTable[0].ipAddr = "3"
// 	server1.node.fingerTable[0].id = intToBytes(3)

// 	server1.node.fingerTable[1] = &FingerEntry{}
// 	server1.node.fingerTable[1].start = big.NewInt(3).Bytes()
// 	server1.node.fingerTable[1].ipAddr = "3"
// 	server1.node.fingerTable[1].id = intToBytes(3)

// 	server1.node.fingerTable[2] = &FingerEntry{}
// 	server1.node.fingerTable[2].start = big.NewInt(5).Bytes()
// 	server1.node.fingerTable[2].ipAddr = "2"
// 	server1.node.fingerTable[2].id = intToBytes(0)

// 	server3.node.fingerTable[0] = &FingerEntry{}
// 	server3.node.fingerTable[0].start = big.NewInt(4).Bytes()
// 	server3.node.fingerTable[0].ipAddr = "2"
// 	server3.node.fingerTable[0].id = intToBytes(0)

// 	server3.node.fingerTable[1] = &FingerEntry{}
// 	server3.node.fingerTable[1].start = big.NewInt(5).Bytes()
// 	server3.node.fingerTable[1].ipAddr = "2"
// 	server3.node.fingerTable[1].id = intToBytes(0)

// 	server3.node.fingerTable[2] = &FingerEntry{}
// 	server3.node.fingerTable[2].start = big.NewInt(7).Bytes()
// 	server3.node.fingerTable[2].ipAddr = "2"
// 	server3.node.fingerTable[2].id = intToBytes(0)

// 	fmt.Println(server0.String())
// 	fmt.Println(server1.String())
// 	fmt.Println(server3.String())

// 	server0.lookup(intToBytes(0))
// }

// func TestJoinOtherFigure3(t *testing.T) {
// 	server1 := MakeServer("31") // id == 1
// 	Servers["31"] = server1
// 	server1.Join(nil)

// 	fmt.Println("joined: " + server1.String())

// 	server3 := MakeServer("3")
// 	Servers["3"] = server3
// 	server3.Join(server1)

// 	fmt.Println("joined: " + server1.String())

// 	server0 := MakeServer("2") // id == 0
// 	Servers["2"] = server0
// 	server0.Join(server1)

// 	fmt.Println(server0.String())
// 	fmt.Println(server1.String())
// 	fmt.Println(server3.String())

// 	checkFingerTable(server0, server1, server3, t)
// }

// func checkFingerTable(server0 *Server, server1 *Server, server3 *Server, t *testing.T) {
// 	fingerTable0 := server0.node.FingerTable()
// 	fingerTable1 := server1.node.FingerTable()
// 	fingerTable3 := server3.node.FingerTable()

// 	if fingerTable0[0].id[0] != 1 || fingerTable0[1].id[0] != 3 || fingerTable0[2].id[0] != 0 {
// 		t.Errorf("fingerTable0 is incorrect.")
// 	}

// 	if fingerTable1[0].id[0] != 3 || fingerTable1[1].id[0] != 3 || fingerTable1[2].id[0] != 0 {
// 		t.Errorf("fingerTable1 is incorrect.")
// 	}

// 	if fingerTable3[0].id[0] != 0 || fingerTable3[1].id[0] != 0 || fingerTable3[2].id[0] != 0 {
// 		t.Errorf("fingerTable3 is incorrect.")
// 	}
// }

// func TestJoinOtherFigure5(t *testing.T) {
// 	server0 := MakeServer("2") // id == 0
// 	Servers["2"] = server0
// 	server0.Join(nil)
// 	fmt.Println(server0.String())

// 	server1 := MakeServer("31") // id == 1
// 	fmt.Println(server1.String())
// 	Servers["31"] = server1
// 	server1.Join(server0)

// 	fmt.Println(server0.String())
// 	fmt.Println(server1.String())

// 	fmt.Println("Adding node3...")

// 	server3 := MakeServer("3")
// 	Servers["3"] = server3
// 	server3.Join(server0)

// 	fmt.Println(server0.String())
// 	fmt.Println(server1.String())
// 	fmt.Println(server3.String())

// 	fmt.Println("Adding node6...")

// 	server6 := MakeServer("6666")
// 	Servers["6666"] = server6

// 	server6.Join(server0)

// 	fmt.Println(server0.String())
// 	fmt.Println(server1.String())
// 	fmt.Println(server3.String())
// 	fmt.Println(server6.String())
// }

// func TestJoinNone(t *testing.T) {
// 	server := MakeServer("2")
// 	server.Join(nil)
// 	fmt.Println(server.String())
// }

// // func Test1(t *testing.T) {
// // 	fmt.Println("Testing make server and finger table:")
// // 	server0 := MakeServer("2")  // ip == "2", id == 0
// // 	server1 := MakeServer("31") // ip == "31", id == 1
// // 	server3 := MakeServer("3")  // ip == "3", id == 3

// // 	Servers["2"] = server0
// // 	Servers["31"] = server1
// // 	Servers["3"] = server3

// // 	server0.node.fingerTable[0] = &FingerEntry{}
// // 	server0.node.fingerTable[0].start = big.NewInt(1).Bytes()
// // 	server0.node.fingerTable[0].ipAddr = "31"
// // 	server0.node.fingerTable[0].id = intToBytes(1)

// // 	server0.node.fingerTable[1] = &FingerEntry{}
// // 	server0.node.fingerTable[1].start = big.NewInt(2).Bytes()
// // 	server0.node.fingerTable[1].ipAddr = "3"
// // 	server0.node.fingerTable[1].id = intToBytes(3)

// // 	server0.node.fingerTable[2] = &FingerEntry{}
// // 	server0.node.fingerTable[2].start = big.NewInt(4).Bytes()
// // 	server0.node.fingerTable[2].ipAddr = "2"
// // 	server0.node.fingerTable[2].id = intToBytes(0)

// // 	server1.node.fingerTable[0] = &FingerEntry{}
// // 	server1.node.fingerTable[0].start = big.NewInt(2).Bytes()
// // 	server1.node.fingerTable[0].ipAddr = "3"
// // 	server1.node.fingerTable[0].id = intToBytes(3)

// // 	server1.node.fingerTable[1] = &FingerEntry{}
// // 	server1.node.fingerTable[1].start = big.NewInt(3).Bytes()
// // 	server1.node.fingerTable[1].ipAddr = "3"
// // 	server1.node.fingerTable[1].id = intToBytes(3)

// // 	server1.node.fingerTable[2] = &FingerEntry{}
// // 	server1.node.fingerTable[2].start = big.NewInt(5).Bytes()
// // 	server1.node.fingerTable[2].ipAddr = "2"
// // 	server1.node.fingerTable[2].id = intToBytes(0)

// // 	server3.node.fingerTable[0] = &FingerEntry{}
// // 	server3.node.fingerTable[0].start = big.NewInt(4).Bytes()
// // 	server3.node.fingerTable[0].ipAddr = "2"
// // 	server3.node.fingerTable[0].id = intToBytes(0)

// // 	server3.node.fingerTable[1] = &FingerEntry{}
// // 	server3.node.fingerTable[1].start = big.NewInt(5).Bytes()
// // 	server3.node.fingerTable[1].ipAddr = "2"
// // 	server3.node.fingerTable[1].id = intToBytes(0)

// // 	server3.node.fingerTable[2] = &FingerEntry{}
// // 	server3.node.fingerTable[2].start = big.NewInt(7).Bytes()
// // 	server3.node.fingerTable[2].ipAddr = "2"
// // 	server3.node.fingerTable[2].id = intToBytes(0)

// // 	fmt.Println(server0.String())
// // 	fmt.Println(server1.String())
// // 	fmt.Println(server3.String())

// // 	ipForKey0 := server3.LookUp(intToBytes(0))
// // 	if ipForKey0 != "2" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey0, "2")
// // 	}
// // 	ipForKey1 := server3.LookUp(intToBytes(1))
// // 	if ipForKey1 != "31" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey1, "31")
// // 	}
// // 	ipForKey2 := server3.LookUp(intToBytes(2))
// // 	if ipForKey2 != "3" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey2, "3")
// // 	}
// // 	ipForKey3 := server3.LookUp(intToBytes(3))
// // 	if ipForKey3 != "3" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey3, "3")
// // 	}
// // 	ipForKey4 := server3.LookUp(intToBytes(4))
// // 	if ipForKey4 != "2" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey4, "2")
// // 	}
// // 	ipForKey5 := server3.LookUp(intToBytes(5))
// // 	if ipForKey5 != "2" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey5, "2")
// // 	}
// // 	ipForKey6 := server3.LookUp(intToBytes(6))
// // 	if ipForKey6 != "2" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey6, "2")
// // 	}
// // 	ipForKey7 := server3.LookUp(intToBytes(7))
// // 	if ipForKey7 != "2" {
// // 		t.Errorf("LookUp was incorrect, got: %s, want: %s.", ipForKey7, "2")
// // 	}
// // }

// // this test should test the betweenRightInclusive func in chord_server.go
// func TestBetweenRightInclusive(t *testing.T) {
// 	fmt.Println("Testing betweenRightInclusive:")
// 	fail := false
// 	// test (3, 0]
// 	//5
// 	got1 := betweenRightInclusive(big.NewInt(5).Bytes(), big.NewInt(3).Bytes(), big.NewInt(0).Bytes())
// 	if got1 != true {
// 		t.Errorf("betweenRightInclusive(5, 3, 0) = %t; want true", got1)
// 		fail = true
// 	}
// 	//0
// 	got2 := betweenRightInclusive(big.NewInt(0).Bytes(), big.NewInt(3).Bytes(), big.NewInt(0).Bytes())
// 	if got2 != true {
// 		t.Errorf("betweenRightInclusive(0, 3, 0) = %t; want true", got2)
// 		fail = true
// 	}
// 	//3
// 	got4 := betweenRightInclusive(big.NewInt(3).Bytes(), big.NewInt(3).Bytes(), big.NewInt(0).Bytes())
// 	if got4 != false {
// 		t.Errorf("betweenRightInclusive(3, 3, 0) = %t; want false", got4)
// 		fail = true
// 	}
// 	//1
// 	got3 := betweenRightInclusive(big.NewInt(1).Bytes(), big.NewInt(3).Bytes(), big.NewInt(0).Bytes())
// 	if got3 != false {
// 		t.Errorf("betweenRightInclusive(1, 3, 0) = %t; want false", got3)
// 		fail = true
// 	}

// 	// test (1, 3]
// 	//2
// 	got5 := betweenRightInclusive(big.NewInt(2).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got5 != true {
// 		t.Errorf("betweenRightInclusive(2, 1, 3) = %t; want true", got5)
// 		fail = true
// 	}
// 	//1
// 	got6 := betweenRightInclusive(big.NewInt(1).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got6 != false {
// 		t.Errorf("betweenRightInclusive(1, 1, 3) = %t; want false", got6)
// 		fail = true
// 	}
// 	//3
// 	got7 := betweenRightInclusive(big.NewInt(3).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got7 != true {
// 		t.Errorf("betweenRightInclusive(3, 1, 3) = %t; want true", got7)
// 		fail = true
// 	}
// 	//4
// 	got8 := betweenRightInclusive(big.NewInt(4).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got8 != false {
// 		t.Errorf("betweenRightInclusive(4, 1, 3) = %t; want false", got8)
// 		fail = true
// 	}
// 	//0
// 	got9 := betweenRightInclusive(big.NewInt(0).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got9 != false {
// 		t.Errorf("betweenRightInclusive(0, 1, 3) = %t; want false", got8)
// 		fail = true
// 	}

// 	if !fail {
// 		fmt.Println("PASS")
// 	}

// }

// // this should test the between func in chord_server.go
// func TestBetween(t *testing.T) {
// 	fmt.Println("Testing between:")
// 	fail := false

// 	// test (1, 3)
// 	//2
// 	got5 := between(big.NewInt(2).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got5 != true {
// 		t.Errorf("between(2, 1, 3) = %t; want true", got5)
// 		fail = true
// 	}
// 	//1
// 	got6 := between(big.NewInt(1).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got6 != false {
// 		t.Errorf("between(1, 1, 3) = %t; want false", got6)
// 		fail = true
// 	}
// 	//3
// 	got7 := between(big.NewInt(3).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got7 != false {
// 		t.Errorf("between(3, 1, 3) = %t; want false", got7)
// 		fail = true
// 	}
// 	//4
// 	got8 := between(big.NewInt(4).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got8 != false {
// 		t.Errorf("between(4, 1, 3) = %t; want false", got8)
// 		fail = true
// 	}
// 	//0
// 	got9 := between(big.NewInt(0).Bytes(), big.NewInt(1).Bytes(), big.NewInt(3).Bytes())
// 	if got9 != false {
// 		t.Errorf("between(0, 1, 3) = %t; want false", got8)
// 		fail = true
// 	}

// 	if !fail {
// 		fmt.Println("PASS")
// 	}
// }

// // TODO: this should test betweenLeftInclusive func in chord_server.go
// func TestBetweenLeftInclusive(t *testing.T) {

// }
