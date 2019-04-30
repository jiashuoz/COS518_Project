package chord

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
)

// "fmt"
// "bytes"
// "math/big"

// ChordServer is a single ChordServer
type ChordServer struct {
	node        Node
	fingerTable []Node // a table of FingerEntry pointer
	predecessor Node   // previous node on the identifier circle
}

// MakeServer returns a pointer to a  server
func MakeServer(ip string) *ChordServer {
	server := &ChordServer{}
	server.node = MakeNode(ip)
	server.fingerTable = make([]Node, numBits)
	return server
}

func (chord *ChordServer) FindSuccessor(id []byte) Node {
	pred := chord.FindPredecessor(id)
	predSucc, _ := pred.GetSuccessorRPC()
	return predSucc
}

func (chord *ChordServer) FindPredecessor(id []byte) Node {
	closest := chord.FindClosestNode(id)

	if idsEqual(closest.ID, chord.node.ID) {
		return closest
	}

	closestSucc, _ := closest.GetSuccessorRPC()

	for !betweenRightInclusive(id, closest.ID, closestSucc.ID) {
		closest, _ := closest.FindClosestNodeRPC(id)
		closestSucc, _ = closest.GetSuccessorRPC()
	}
	return closest
}

func (chord *ChordServer) FindClosestNode(id []byte) Node {
	fingerTable := chord.fingerTable
	for i := numBits - 1; i >= 0; i-- {
		if fingerTable[i].ID != nil && between(fingerTable[i].ID, chord.GetID(), id) {
			return fingerTable[i]
		}
	}
	return fingerTable[numBits-1]
}

// GetNode returns ch's network information.
func (chord *ChordServer) GetNode() Node {
	return chord.node
}

// GetID return's ch's identifier.
func (chord *ChordServer) GetID() []byte {
	return chord.node.ID
}

// node thinks it might be chord's predecessor.
func (chord *ChordServer) Notify(node Node) error {
	//TODO: lock here since it is changing chordServer property
	// notify itself
	if bytes.Equal(chord.node.ID, node.ID) {
		return nil
	}
	// node is the only node in the ring
	if bytes.Equal(chord.node.ID, chord.predecessor.ID) {
		chord.predecessor = node
		chord.fingerTable[0] = node
	} else if betweenRightInclusive(node.ID, chord.predecessor.ID, chord.node.ID) {
		chord.predecessor = node
	}
	return nil
}

func (chord *ChordServer) fingerStart(i int) []byte {
	currID := new(big.Int).SetBytes(chord.node.ID)
	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
	start := new(big.Int).Add(currID, offset)
	return start.Bytes()
}

// periodically fresh finger table entries
func (chord *ChordServer) fixFingers() {
	i := rand.Intn(numBits-1) + 1
	fingerStart := chord.fingerStart(i)
	chord.fingerTable[i] = chord.FindSuccessor(fingerStart)
}

// // Join adds chordServer to the network
// func (chordServer *ChordServer) Join(exisitingServer *ChordServer) {
// 	if exisitingServer == nil { // the only node in the network
// 		for _, entry := range chordServer.node.FingerTable() {
// 			entry.id = chordServer.node.ID()
// 			entry.ipAddr = chordServer.node.IP()
// 		}
// 		chordServer.node.predecessor = &Node{id: chordServer.node.ID(), ipAddr: chordServer.node.IP()}
// 	} else { // update other nodes' fingers
// 		chordServer.InitFingerTable(exisitingServer)
// 		chordServer.UpdateOthers()
// 	}
// }

// func (chordServer *ChordServer) InitFingerTable(existingServer *ChordServer) {
// 	DPrintf("Node%v: Start Initializing Finger Table...", chordServer.node.ID())
// 	currNode := chordServer.node
// 	fingerTable := currNode.FingerTable()
// 	successor := existingServer.FindSuccessor(fingerTable[0].start)
// 	successorServer := Servers[successor.ipAddr]
// 	DPrintf("successor should be 0: %v\n", successorServer.node.ID())
// 	// curr.prev = node.prev
// 	currNode.predecessor = &Node{id: successorServer.node.predecessor.id, ipAddr: successorServer.node.predecessor.ipAddr}
// 	currNode.SetSuccessor(&Node{id: successorServer.node.ID(), ipAddr: successorServer.node.IP()})
// 	// node.prev = curr
// 	//successorServer.node.predecessor = &Node{id: currNode.ID(), ipAddr: chordServer.node.ipAddr}
// 	successorServer.node.SetPredecessor(currNode)
// 	ChangeServer(currNode.Predecessor().IP()).node.SetSuccessor(currNode)
// 	DPrintf("first finger is: id: %v  ip: %v", currNode.Successor().ID(), currNode.Successor().IP())
// 	for i := 1; i < numBits; i++ {
// 		DPrintf("initializing %vth finger", i)
// 		if betweenLeftInclusive(fingerTable[i].start, currNode.ID(), fingerTable[i-1].id) {
// 			DPrintf("%v should be between %v and %v\n",
// 				fingerTable[i].start,
// 				currNode.ID(),
// 				fingerTable[i-1].id)
// 			fingerTable[i].id = fingerTable[i-1].id
// 			fingerTable[i].ipAddr = fingerTable[i-1].ipAddr
// 		} else {
// 			DPrintf("else, find successor based on fingerTable")
// 			fartherNode := existingServer.FindSuccessor(fingerTable[i].start)
// 			fingerTable[i].id = fartherNode.id
// 			fingerTable[i].ipAddr = fartherNode.ipAddr
// 		}
// 	}
// 	DPrintf("Done initializing:\n")
// 	DPrintf(chordServer.node.String())

// }

// // Update all nodes whose finger tables should refer to chordServer
// func (chordServer *ChordServer) UpdateOthers() {
// for i := 0; i < numBits; i++ {
// 	// p = find_predecessor(n - 2^i)
// 	DPrintf("n-2^i = %v", chordServer.nodeIdToUpdateFinger(i))
// 	p := chordServer.FindPredecessor(chordServer.nodeIdToUpdateFinger(i))
// 	DPrintf("predecessor: %v", p.ID())

// 	if bytes.Compare(p.ID(), chordServer.node.ID()) == 0 {
// 		DPrintf("reached new node itself")
// 		return
// 	}

// 	if bytes.Compare(p.Successor().ID(), chordServer.nodeIdToUpdateFinger(i)) == 0 {
// 		p = p.Successor()
// 	}
// 	pServer := ChangeServer(p.IP())
// 	pServer.UpdateFingerTable(chordServer.node, i)
// }
// }

// // Returns the id of node whose ith finger might be chordServer
// func (chordServer *ChordServer) nodeIdToUpdateFinger(i int) []byte {
// 	n := new(big.Int).SetBytes(chordServer.node.ID())
// 	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
// 	diff := new(big.Int).Sub(n, offset)
// 	// diff.Add(diff, big.NewInt(1))

// 	if diff.Sign() < 0 {
// 		diff = diff.Add(diff, new(big.Int).Exp(big.NewInt(2), big.NewInt(numBits), nil))
// 	}

// 	if diff.Cmp(big.NewInt(0)) == 0 {
// 		return []byte{0}
// 	}

// 	return diff.Bytes()
// }

// // Update chordServer's finger if s should be the ith finger
// func (chordServer *ChordServer) UpdateFingerTable(s *Node, i int) {
// 	if bytes.Compare(s.ID(), chordServer.node.ID()) == 0 {
// 		DPrintf("reached new node itself")
// 		return
// 	}
// 	DPrintf("update %v's finger", chordServer.node.id)
// 	fingerTable := chordServer.node.fingerTable
// 	if betweenLeftInclusive(s.ID(), chordServer.node.ID(), fingerTable[i].id) {
// 		DPrintf("yes")
// 		fingerTable[i].id = s.ID()
// 		fingerTable[i].ipAddr = s.IP()
// 		DPrintf("fingerTable: %v", chordServer.node.String())
// 		p := chordServer.node.Predecessor()
// 		pServer := ChangeServer(p.IP())
// 		pServer.UpdateFingerTable(s, i)
// 	}

// }

// // LookUp returns the ip addr of the successor node of id
// func (chordServer *ChordServer) LookUp(id []byte) string {
// 	ipAddr := chordServer.FindSuccessor(id).ipAddr
// 	DPrintf("lookup result: " + ipAddr)
// 	return ipAddr
// }

// // FindSuccessor returns the successor node of id
// func (chordServer *ChordServer) FindSuccessor(id []byte) *Node {
// 	predecessor := chordServer.FindPredecessor(id)
// 	return Servers[predecessor.ipAddr].node.Successor()
// }

func (chord *ChordServer) String(printFingerTable bool) string {
	str := chord.node.String()
	if !printFingerTable {
		return str
	}

	str += "Finger table: \n"
	str += "ith | start | successor\n"
	for i := 0; i < numBits; i++ {
		currID := new(big.Int).SetBytes(chord.node.ID)
		offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		start := new(big.Int).Add(currID, offset)
		successor := chord.fingerTable[i].ID
		str += fmt.Sprintf("%d   | %d     | %d\n", i, start, successor)
	}
	return str
}
