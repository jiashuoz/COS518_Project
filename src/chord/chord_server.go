package chord

import (
	"fmt"
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

func (chordServer *Server) LookUp(id []byte) string {
	return chordServer.FindSuccessor(id).ipAddr
}

func (chordServer *Server) FindSuccessor(id []byte) *NodeInfo {
	predecessor := chordServer.FindPredecessor(id)

	return Servers[predecessor.ipAddr].node.successor
}

// FindPredecessor returns the previous node in the circle to id
func (chordServer *Server) FindPredecessor(id []byte) *NodeInfo {
	currServer := chordServer

	for !betweenRightInclusive(id, currServer.node.id, currServer.node.successor.id) {
		closerNodeInfo := currServer.closestPrecedingFinger(id)
		fmt.Println(closerNodeInfo)
		currServer = ChangeServer(closerNodeInfo.ipAddr)
	}

	return &NodeInfo{currServer.node.id, currServer.ipAddr}
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

// Returns the closest node preceding id: previous node in the circle to id
func (chordServer *Server) closestPrecedingFinger(id []byte) *NodeInfo {
	// slice contains pointers to TableEntry
	fingerTable := chordServer.node.FingerTable()
	currID := chordServer.node.ID()
	for i := numBits - 1; i >= 0; i-- {
		if fingerTable[i] != nil {
			finger := fingerTable[i]
			if between(finger.id, currID, id) {
				return &NodeInfo{finger.id, finger.ipAddr}
			}
		}
	}
	return nil
}

// Returns true if begin < target < end
func between(target []byte, begin []byte, end []byte) bool {
	targetBigInt := big.NewInt(0).SetBytes(target)
	beginBigInt := big.NewInt(0).SetBytes(begin)
	endBigInt := big.NewInt(0).SetBytes(end)

	if beginBigInt.Cmp(endBigInt) == 1 { // (3, 2)
		return targetBigInt.Cmp(beginBigInt) == 1 || targetBigInt.Cmp(endBigInt) == -1
	}
	// (2, 3) or (3, 3)
	return targetBigInt.Cmp(beginBigInt) == 1 && targetBigInt.Cmp(endBigInt) == -1
}

func (chordServer *Server) String() string {
	// str := "server name: " + chordServer.name + "\n"
	str := "Server IP: " + chordServer.ipAddr + "\n"
	str += chordServer.node.String() + "\n"
	return str
}
