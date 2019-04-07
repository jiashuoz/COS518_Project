package chord

import (
	"crypto/sha1"
	"math/big"
)

// A primitive node
type Node struct {
	id          []byte         // use sha1 to generate 160 bit id (20 bytes)
	fingerTable []*FingerEntry // a table of FingerEntry
	successor   *NodeInfo      // next node on the identifier circle
	predecessor *NodeInfo      // previous node on the identifier circle
}

// A NodeInfo contains some basic information about a node
type NodeInfo struct {
	id     []byte // id of the node
	ipAddr string // ip address of the node
}

// Each entry in fingerTable
type FingerEntry struct {
	start []byte    // start == (n + 2^k-1) mod 2^m, 1 <= k <= m
	succ  *NodeInfo // the next node >= fingerEntry.
}

// Create a new Node based on ip address
func MakeNode(ipAddr string) *Node {

	n := Node{}
	n.id = hash(ipAddr) // id
	n.fingerTable = make([]*FingerEntry, 160)
	n.successor = MakeNodeInfo(n.id, ipAddr) // initially, successor is itself
	n.predecessor = nil                      // initially, no predecessor

	return &n
}

// Create a new NodeInfo give id and ipAddr
func MakeNodeInfo(id []byte, ipAddr string) *NodeInfo {
	return &NodeInfo{id, ipAddr}
}

// Return node's ID
func (node *Node) ID() []byte {
	return node.id
}

// Return the finger table
func (node *Node) FingerTable() []*FingerEntry {
	return node.fingerTable
}

// Return the info about successor
func (node *Node) Successor() *NodeInfo {
	return node.successor
}

// Return the info about predecessor
func (node *Node) Predecessor() *NodeInfo {
	return node.predecessor
}

// Returns the hash value in []byte based an ipAddr
// Takes in a string, use sha1 to hash it to generate 160 bit hash value
// mod by 2^160 so that the value's range is [0, 2^m - 1]
func hash(ipAddr string) []byte {
	h := sha1.New()
	h.Write([]byte(ipAddr))

	idBigInt := big.NewInt(0)
	idBigInt.SetBytes(h.Sum(nil)) // Sum() returns []byte, convert it into BigInt

	maxVal := big.NewInt(0)
	maxVal.Exp(big.NewInt(2), big.NewInt(160), nil) // calculate 2^m
	idBigInt.Mod(idBigInt, maxVal)                  // mod id to make it to be [0, 2^m - 1]

	return idBigInt.Bytes()
}
