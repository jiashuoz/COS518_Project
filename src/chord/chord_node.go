package chord

import (
	"crypto/sha1"
	"encoding/hex"
	"math/big"
)

// Node is a primitive structure that contains information about a chord
type Node struct {
	id          []byte         // use sha1 to generate 160 bit id (20 bytes)
	fingerTable []*FingerEntry // a table of FingerEntry pointer
	successor   *NodeInfo      // next node on the identifier circle
	predecessor *NodeInfo      // previous node on the identifier circle
}

// NodeInfo contains some basic information about a node
type NodeInfo struct {
	id     []byte // id of the node, sha1 generated 160 bit id (20 bytes)
	ipAddr string // ip address of the node
}

// FingerEntry in fingerTable
type FingerEntry struct {
	start  []byte // start == (n + 2^k-1) mod 2^m, 1 <= k <= m
	id     []byte
	ipAddr string
}

// MakeNode creates a new Node based on ip address and returns a pointer to it
func MakeNode(ipAddr string) *Node {

	n := Node{}
	n.id = hash(ipAddr) // id
	n.fingerTable = make([]*FingerEntry, numBits)
	n.successor = MakeNodeInfo(n.id, ipAddr) // initially, successor is itself
	n.predecessor = nil                      // initially, no predecessor

	return &n
}

// MakeNodeInfo creates a new NodeInfo give id and ipAddr and returns a pointer to it
func MakeNodeInfo(id []byte, ipAddr string) *NodeInfo {
	return &NodeInfo{id, ipAddr}
}

// ID returns node's ID in []byte
func (node *Node) ID() []byte {
	return node.id
}

// FingerTable returns a pointer to an array of table entry pointers
func (node *Node) FingerTable() []*FingerEntry {
	return node.fingerTable
}

// Successor returns a pointer to a NodeInfo struct about successor
func (node *Node) Successor() *NodeInfo {
	return node.successor
}

// SetSuccessor sets successor field
func (node *Node) SetSuccessor(newSucc *NodeInfo) {
	node.successor = newSucc
}

// Predecessor returns a pointer to a NodeInfo struct about predecessor
func (node *Node) Predecessor() *NodeInfo {
	return node.predecessor
}

// SetPredecessor sets successor field
func (node *Node) SetPredecessor(newPred *NodeInfo) {
	node.predecessor = newPred
}

// String returns the string representation of a Node
func (node *Node) String() string {
	str := "Node: \n"
	str += "id: " + hex.EncodeToString(node.id)
	str += "fingerTable: \n"
	for _, entry := range node.fingerTable {
		if entry != nil {
			str += "   " + entry.string()
		}
	}

	if node.successor != nil {
		str += "successor: " + node.successor.string()
	} else {
		str += "successor: nil"
	}

	if node.predecessor != nil {
		str += "predecessor: " + node.predecessor.string()
	} else {
		str += "predecessor: nil"
	}
	return str
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
	maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil) // calculate 2^m
	idBigInt.Mod(idBigInt, maxVal)                      // mod id to make it to be [0, 2^m - 1]

	return idBigInt.Bytes()
}

func (fingerEntry *FingerEntry) string() string {
	str := "finger entry: "
	str += hex.EncodeToString(fingerEntry.start) + " "
	str += hex.EncodeToString(fingerEntry.id)
	str += string('\n')
	return str
}

func (nodeInfo *NodeInfo) string() string {
	str := "NodeInfo: "
	str += hex.EncodeToString(nodeInfo.id) + " "
	str += nodeInfo.ipAddr + "\n"
	return str
}
