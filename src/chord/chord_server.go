package chord

import (
	"math/big"
)

// Server is a single ChordServer
type Server struct {
	name   string
	ipAddr string
	node   *Node
}

// MakeServer returns a pointer to a  server
func MakeServer(ip string) *Server {
	server := &Server{
		name:   "server 1",
		ipAddr: ip,
		node:   MakeNode(ip),
	}
	return server
}

// Join adds chordServer to the network
func (chordServer *Server) Join(exisitingServer *Server) {
	if exisitingServer == nil { // the only node in the network
		for _, entry := range chordServer.node.FingerTable() {
			entry.id = chordServer.node.ID()
			entry.ipAddr = chordServer.ipAddr
		}
		chordServer.node.predecessor = &NodeInfo{chordServer.node.ID(), chordServer.ipAddr}
	} else {

	}
}

func (chordServer *Server) InitFingerTable(existingServer *Server) {
	currNode := chordServer.node
	fingerTable := currNode.FingerTable()
	successor := existingServer.FindSuccessor(fingerTable[0].start)
	successorServer := Servers[successor.ipAddr]
	// curr.prev = node.prev
	currNode.predecessor = &NodeInfo{successorServer.node.predecessor.id, successorServer.node.predecessor.ipAddr}
	// node.prev = curr
	successorServer.node.predecessor = &NodeInfo{currNode.ID(), chordServer.ipAddr}

	for i := 1; i < numBits; i++ {
		if betweenLeftInclusive(fingerTable[i].start, currNode.ID(), fingerTable[i-1].id) {
			fingerTable[i].id = fingerTable[i-1].id
			fingerTable[i].ipAddr = fingerTable[i-1].ipAddr
		} else {
			fartherNodeInfo := existingServer.FindSuccessor(fingerTable[i+1].start)
			fingerTable[i].id = fartherNodeInfo.id
			fingerTable[i].ipAddr = fartherNodeInfo.ipAddr
		}
	}
}

// LookUp returns the ip addr of the successor node of id
func (chordServer *Server) LookUp(id []byte) string {
	return chordServer.FindSuccessor(id).ipAddr
}

// FindSuccessor returns the successor node of id
func (chordServer *Server) FindSuccessor(id []byte) *NodeInfo {
	predecessor := chordServer.FindPredecessor(id)
	return Servers[predecessor.ipAddr].node.Successor()
}

// FindPredecessor returns the previous node in the circle to id
func (chordServer *Server) FindPredecessor(id []byte) *NodeInfo {
	currServer := chordServer
	DPrintf("looking for : %d\n", id)
	DPrintf("curr node id: %d\n", currServer.node.id)
	DPrintf("curr node successor id: %d\n", currServer.node.Successor().id)
	for !betweenRightInclusive(id, currServer.node.id, currServer.node.Successor().id) {
		closerNodeInfo := currServer.closestPrecedingFinger(id)
		currServer = ChangeServer(closerNodeInfo.ipAddr)
	}
	DPrintf("predecessor of id: %d\n", currServer.node.id)
	return &NodeInfo{currServer.node.id, currServer.ipAddr}
}

// Returns the closest node preceding id: previous node in the circle to id
func (chordServer *Server) closestPrecedingFinger(id []byte) *NodeInfo {
	// slice contains pointers to TableEntry
	fingerTable := chordServer.node.FingerTable()
	currID := chordServer.node.ID()
	for i := numBits - 1; i >= 0; i-- {
		if fingerTable[i].id != nil {
			finger := fingerTable[i]
			if between(finger.id, currID, id) {
				return &NodeInfo{finger.id, finger.ipAddr}
			}
		}
	}
	return nil
}

func (chordServer *Server) String() string {
	// str := "server name: " + chordServer.name + "\n"
	str := "Server IP: " + chordServer.ipAddr + "\n"
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
	// begin < end, (2, 3]
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

	// [2, 3)
	return targetBigInt.Cmp(beginBigInt) == 0 ||
		targetBigInt.Cmp(beginBigInt) == 1 ||
		targetBigInt.Cmp(endBigInt) == -1
}
