package chord

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

// Node is a primitive structure that contains information about a chord
type Node struct {
	IP string
	ID []byte // use sha1 to generate 160 bit id (20 bytes)
}

// MakeNode creates a new Node based on ip address and returns a pointer to it
func MakeNode(ipAddr string) Node {
	n := Node{ipAddr, hash(ipAddr)}
	return n
}

// String returns string representation of n.
func (n *Node) String() string {
	return fmt.Sprintf("Server IP: %s, ID: %d\n", n.IP, n.ID)
}

// // Returns the hash value in []byte based an ipAddr
// // Takes in a string, use sha1 to hash it to generate 160 bit hash value
// // mod by 2^160 so that the value's range is [0, 2^m - 1]
func hash(ipAddr string) []byte {
	h := sha1.New()
	h.Write([]byte(ipAddr))

	idInt := big.NewInt(0)
	idInt.SetBytes(h.Sum(nil)) // Sum() returns []byte, convert it into BigInt

	maxVal := big.NewInt(0)
	maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil) // calculate 2^m
	idInt.Mod(idInt, maxVal)                            // mod id to make it to be [0, 2^m - 1]
	if idInt.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}
	return idInt.Bytes()
}
