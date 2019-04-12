package chord

import (
	"math/big"
)

type ChordServer struct {
	name   string
	ipAddr string
	node   *Node
}

// Returns the closest node preceding id: previous node in the circle to id
func (chordServer *ChordServer) closestPrecedingFinger(id []byte) *NodeInfo {
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
