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

func (chordServer *Server) FindSuccessor(id []byte) string {
	return ""
}

// FindPredecessor returns the previous node in the circle to id
func (chordServer *Server) FindPredecessor(id []byte) *NodeInfo {
	currNode := chordServer.node
	if betweenRightInclusive(id, currNode.id, currNode.successor.id) {
		return &NodeInfo{currNode.id, chordServer.ipAddr}
	}
	return chordServer.closestPrecedingFinger(id)
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
