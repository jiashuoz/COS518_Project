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

// FindPredecessor returns the previous node in the circle to id
func (chordServer *Server) FindPredecessor(id []byte) *NodeInfo {
	currNode := chordServer.node
	if betweenRightInclusive(id, currNode.me.id, currNode.successor.id) {
		return currNode.me
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
	n := chordServer.node.ID()
	for i := 159; i >= 0; i-- {
		if fingerTable[i] != nil {
			entry := fingerTable[i]
			if between(entry.succ.id, n, id) {
				return entry.succ
			}
		}
	}
	return chordServer.node.me
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
