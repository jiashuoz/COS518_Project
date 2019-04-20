package chord

import (
	"crypto/sha1"
	"encoding/hex"
	"math/big"
)

// Node is a primitive structure that contains information about a chord
type Node struct {
	ipAddr      string
	id          []byte         // use sha1 to generate 160 bit id (20 bytes)
	fingerTable []*FingerEntry // a table of FingerEntry pointer
	predecessor *Node          // previous node on the identifier circle
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
	n.ipAddr = ipAddr
	n.id = hash(ipAddr)
	n.fingerTable = make([]*FingerEntry, numBits)

	for i := range n.fingerTable {
		n.fingerTable[i] = &FingerEntry{}
		// n.id + 2^i
		n.fingerTable[i].start = fingerStart(n.id, i)
		DPrintf("start: %d\n", n.fingerTable[i].start)
	}
	n.predecessor = nil // initially, no predecessor

	return &n
}

func fingerStart(id []byte, i int) []byte {
	nodeID := new(big.Int).SetBytes(id)                                  // id
	addend := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil) // 2^i

	nodeID.Add(nodeID, addend)
	maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(numBits), nil)

	nodeID.Mod(nodeID, maxVal)

	if nodeID.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}
	return nodeID.Bytes()
}

// MakeNodeInfo creates a new NodeInfo give id and ipAddr and returns a pointer to it
func MakeNodeInfo(id []byte, ipAddr string) *NodeInfo {
	return &NodeInfo{id, ipAddr}
}

// ID returns node's ID in []byte
func (node *Node) ID() []byte {
	return node.id
}

// IP returns node's IP in string
func (node *Node) IP() string {
	return node.ipAddr
}

// FingerTable returns a pointer to an array of table entry pointers
func (node *Node) FingerTable() []*FingerEntry {
	return node.fingerTable
}

// Successor returns a pointer to a NodeInfo struct about successor
func (node *Node) Successor() *Node {
	return &Node{id: node.fingerTable[0].id, ipAddr: node.fingerTable[0].ipAddr}
}

// SetSuccessor sets successor field
func (node *Node) SetSuccessor(newSucc *Node) {
	node.fingerTable[0].id = newSucc.id
	node.fingerTable[0].ipAddr = newSucc.ipAddr
}

// Predecessor returns a pointer to a NodeInfo struct about predecessor
func (node *Node) Predecessor() *Node {
	return node.predecessor
}

// SetPredecessor sets successor field
func (node *Node) SetPredecessor(newPred *Node) {
	node.predecessor = newPred
}

// String returns the string representation of a Node
func (node *Node) String() string {
	str := "Node: \n"
	str += "id: " + hex.EncodeToString(node.id) + "\n"
	str += "fingerTable: \n"
	for _, entry := range node.fingerTable {
		if entry != nil {
			str += "   " + entry.string()
		}
	}

	if node.Successor() != nil {
		str += "successor: " + node.Successor().string()
	} else {
		str += "successor: nil"
	}

	if node.Predecessor() != nil {
		str += "predecessor: " + node.Predecessor().string()
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

	maxVal := new(big.Int).Exp(big.NewInt(2), big.NewInt(numBits), nil)
	idBigInt.Mod(idBigInt, maxVal) // mod id to make it to be [0, 2^m - 1]
	if idBigInt.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}
	return idBigInt.Bytes()
}

func (fingerEntry *FingerEntry) string() string {
	str := "finger entry: "
	str += hex.EncodeToString(fingerEntry.start) + " "
	str += "\"" + fingerEntry.ipAddr + "\" "
	str += hex.EncodeToString(fingerEntry.id) + "\n"
	return str
}

func (node *Node) string() string {
	str := "Node: "
	str += hex.EncodeToString(node.id) + " "
	str += node.ipAddr + "\n"
	return str
}
