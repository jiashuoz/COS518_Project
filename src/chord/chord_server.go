package chord

import (
	"bytes"
	"math/big"
)

// ChordServer is a single ChordServer
type ChordServer struct {
	name string
	node *Node
}

// MakeServer returns a pointer to a  server
func MakeServer(ip string) *ChordServer {
	server := &ChordServer{
		name: "server 1",
		node: MakeNode(ip),
	}
	return server
}

// Join adds chordServer to the network
func (chordServer *ChordServer) Join(exisitingServer *ChordServer) {
	if exisitingServer == nil { // the only node in the network
		for _, entry := range chordServer.node.FingerTable() {
			entry.id = chordServer.node.ID()
			entry.ipAddr = chordServer.node.IP()
		}
		chordServer.node.predecessor = &Node{id: chordServer.node.ID(), ipAddr: chordServer.node.IP()}
	} else { // update other nodes' fingers
		chordServer.InitFingerTable(exisitingServer)
		chordServer.UpdateOthers()
	}
}

func (chordServer *ChordServer) InitFingerTable(existingServer *ChordServer) {
	DPrintf("Node%v: Start Initializing Finger Table...", chordServer.node.ID())
	currNode := chordServer.node
	fingerTable := currNode.FingerTable()
	successor := existingServer.FindSuccessor(fingerTable[0].start)
	successorServer := Servers[successor.ipAddr]
	DPrintf("successor should be 0: %v\n", successorServer.node.ID())
	// curr.prev = node.prev
	currNode.predecessor = &Node{id: successorServer.node.predecessor.id, ipAddr: successorServer.node.predecessor.ipAddr}
	currNode.SetSuccessor(&Node{id: successorServer.node.ID(), ipAddr: successorServer.node.IP()})
	// node.prev = curr
	//successorServer.node.predecessor = &Node{id: currNode.ID(), ipAddr: chordServer.node.ipAddr}
	successorServer.node.SetPredecessor(currNode)
	ChangeServer(currNode.Predecessor().IP()).node.SetSuccessor(currNode)
	DPrintf("first finger is: id: %v  ip: %v", currNode.Successor().ID(), currNode.Successor().IP())
	for i := 1; i < numBits; i++ {
		DPrintf("initializing %vth finger", i)
		if betweenLeftInclusive(fingerTable[i].start, currNode.ID(), fingerTable[i-1].id) {
			DPrintf("%v should be between %v and %v\n",
				fingerTable[i].start,
				currNode.ID(),
				fingerTable[i-1].id)
			fingerTable[i].id = fingerTable[i-1].id
			fingerTable[i].ipAddr = fingerTable[i-1].ipAddr
		} else {
			DPrintf("else, find successor based on fingerTable")
			fartherNode := existingServer.FindSuccessor(fingerTable[i].start)
			fingerTable[i].id = fartherNode.id
			fingerTable[i].ipAddr = fartherNode.ipAddr
		}
	}
	DPrintf("Done initializing:\n")
	DPrintf(chordServer.node.String())

}

// Update all nodes whose finger tables should refer to chordServer
func (chordServer *ChordServer) UpdateOthers() {
	for i := 0; i < numBits; i++ {
		// p = find_predecessor(n - 2^i)
		DPrintf("n-2^i = %v", chordServer.nodeIdToUpdateFinger(i))
		p := chordServer.FindPredecessor(chordServer.nodeIdToUpdateFinger(i))
		DPrintf("predecessor: %v", p.ID())

		if bytes.Compare(p.ID(), chordServer.node.ID()) == 0 {
			DPrintf("reached new node itself")
			return
		}

		if bytes.Compare(p.Successor().ID(), chordServer.nodeIdToUpdateFinger(i)) == 0 {
			p = p.Successor()
		}
		pServer := ChangeServer(p.IP())
		pServer.UpdateFingerTable(chordServer.node, i)
	}
}

// Returns the id of node whose ith finger might be chordServer
func (chordServer *ChordServer) nodeIdToUpdateFinger(i int) []byte {
	n := new(big.Int).SetBytes(chordServer.node.ID())
	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
	diff := new(big.Int).Sub(n, offset)
	// diff.Add(diff, big.NewInt(1))

	if diff.Sign() < 0 {
		diff = diff.Add(diff, new(big.Int).Exp(big.NewInt(2), big.NewInt(numBits), nil))
	}

	if diff.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}

	return diff.Bytes()
}

// Update chordServer's finger if s should be the ith finger
func (chordServer *ChordServer) UpdateFingerTable(s *Node, i int) {
	if bytes.Compare(s.ID(), chordServer.node.ID()) == 0 {
		DPrintf("reached new node itself")
		return
	}
	DPrintf("update %v's finger", chordServer.node.id)
	fingerTable := chordServer.node.fingerTable
	if betweenLeftInclusive(s.ID(), chordServer.node.ID(), fingerTable[i].id) {
		DPrintf("yes")
		fingerTable[i].id = s.ID()
		fingerTable[i].ipAddr = s.IP()
		DPrintf("fingerTable: %v", chordServer.node.String())
		p := chordServer.node.Predecessor()
		pServer := ChangeServer(p.IP())
		pServer.UpdateFingerTable(s, i)
	}

}

// LookUp returns the ip addr of the successor node of id
func (chordServer *ChordServer) LookUp(id []byte) string {
	ipAddr := chordServer.FindSuccessor(id).ipAddr
	DPrintf("lookup result: " + ipAddr)
	return ipAddr
}

// FindSuccessor returns the successor node of id
func (chordServer *ChordServer) FindSuccessor(id []byte) *Node {
	predecessor := chordServer.FindPredecessor(id)
	return Servers[predecessor.ipAddr].node.Successor()
}

// FindPredecessor returns the previous node in the circle to id
func (chordServer *ChordServer) FindPredecessor(id []byte) *Node {
	currServer := chordServer
	for !betweenRightInclusive(id, currServer.node.id, currServer.node.Successor().id) {
		closerNode := currServer.closestPrecedingFinger(id)
		currServer = ChangeServer(closerNode.ipAddr)
	}
	return currServer.node
}

// Returns the closest node preceding id: previous node in the circle to id
func (chordServer *ChordServer) closestPrecedingFinger(id []byte) *Node {
	// slice contains pointers to TableEntry
	fingerTable := chordServer.node.FingerTable()
	currID := chordServer.node.ID()
	for i := numBits - 1; i >= 0; i-- {
		if fingerTable[i].id != nil {
			finger := fingerTable[i]
			if between(finger.id, currID, id) {
				return &Node{id: finger.id, ipAddr: finger.ipAddr}
			}
		}
	}
	return nil
}

func (chordServer *ChordServer) String() string {
	// str := "server name: " + chordServer.name + "\n"
	str := "ChordServer IP: " + chordServer.node.ipAddr + "\n"
	str += chordServer.node.String() + "\n"
	return str
}

func betweenRightInclusive(target []byte, begin []byte, end []byte) bool {
	targetBigInt := big.NewInt(0).SetBytes(target)
	beginBigInt := big.NewInt(0).SetBytes(begin)
	endBigInt := big.NewInt(0).SetBytes(end)

	if beginBigInt.Cmp(endBigInt) == 1 { // begin > end, (3, 0]
		return targetBigInt.Cmp(beginBigInt) == 1 || targetBigInt.Cmp(endBigInt) == -1 || targetBigInt.Cmp(endBigInt) == 0
	}

	if beginBigInt.Cmp(endBigInt) == 0 {
		return true
	}
	// begin < end, (2, 3] or begin == end (0, 0]
	return targetBigInt.Cmp(beginBigInt) == 1 && (targetBigInt.Cmp(endBigInt) == -1 || targetBigInt.Cmp(endBigInt) == 0)
}

// Returns true if begin < target < end
func between(target []byte, begin []byte, end []byte) bool {
	targetBigInt := big.NewInt(0).SetBytes(target)
	beginBigInt := big.NewInt(0).SetBytes(begin)
	endBigInt := big.NewInt(0).SetBytes(end)

	if beginBigInt.Cmp(endBigInt) == 1 || beginBigInt.Cmp(endBigInt) == 0 { // (3, 2), or (3, 3)
		return targetBigInt.Cmp(beginBigInt) == 1 || targetBigInt.Cmp(endBigInt) == -1
	}
	// (2, 3)
	return targetBigInt.Cmp(beginBigInt) == 1 && targetBigInt.Cmp(endBigInt) == -1
}

// Returns true if begin <= target < end, in the ring
func betweenLeftInclusive(target []byte, begin []byte, end []byte) bool {
	targetBigInt := big.NewInt(0).SetBytes(target)
	beginBigInt := big.NewInt(0).SetBytes(begin)
	endBigInt := big.NewInt(0).SetBytes(end)

	// [2, 3)
	if beginBigInt.Cmp(endBigInt) == -1 {
		return (targetBigInt.Cmp(beginBigInt) == 1 || targetBigInt.Cmp(beginBigInt) == 0) &&
			targetBigInt.Cmp(endBigInt) == -1
	}

	if beginBigInt.Cmp(endBigInt) == 0 {
		return true
	}

	// [3, 2) or [0, 0)
	return targetBigInt.Cmp(beginBigInt) == 0 || // target == begin
		targetBigInt.Cmp(beginBigInt) == 1 || // target > begin
		targetBigInt.Cmp(endBigInt) == -1 // target < end
}
