package chord

// import (
// 	"crypto/sha1"
// 	"encoding/hex"
// 	"log"
// 	"math/big"
// 	"net"
// 	"net/rpc"
// 	"time"
// )

// // Node is a primitive structure that contains information about a chord
// type Node struct {
// 	IP          string
// 	ID          []byte  // use sha1 to generate 160 bit id (20 bytes)
// }

// // FingerEntry in fingerTable
// type FingerEntry struct {
// 	start []byte // start == (n + 2^k-1) mod 2^m, 1 <= k <= m
// 	node  *Node
// }

// for i := range n.FingerTable {
// 	n.FingerTable[i] = &Node{}
// }
// n.FingerTable = make([]*Node, numBits)

// // MakeNode creates a new Node based on ip address and returns a pointer to it
// func MakeNode(ipAddr string) *Node {
// 	n := Node{ipAddr, hash(ipAddr)}
// 	DPrintf("Initialized node ----> ip: %v | id: %v", n.IP, n.ID)
// 	return &n
// }

// // ID returns node's ID in []byte
// func (node *Node) GetID() []byte {
// 	return node.ID
// }

// // IP returns node's IP in string
// func (node *Node) GetIP() string {
// 	return node.IP
// }

// // FingerTable returns a pointer to an array of table entry pointers
// func (node *Node) GetFingerTable() []*Node {
// 	return node.FingerTable
// }

// // Successor returns a pointer to a NodeInfo struct about successor
// func (node *Node) GetSuccessor() *Node {
// 	return node.FingerTable[0]
// }

// // SetSuccessor sets successor field
// func (node *Node) SetSuccessor(newSucc *Node) {
// 	node.FingerTable[0] = newSucc
// }

// // Predecessor returns a pointer to a NodeInfo struct about predecessor
// func (node *Node) GetPredecessor() *Node {
// 	return node.Predecessor
// }

// // SetPredecessor sets successor field
// func (node *Node) SetPredecessor(newPred *Node) {
// 	node.Predecessor = newPred
// }

// // String returns the string representation of a Node
// func (node *Node) String() string {
// 	str := "Node: \n"
// 	str += "id: " + hex.EncodeToString(node.ID) + "\n"
// 	str += "fingerTable: \n"
// 	for _, entry := range node.FingerTable {
// 		if entry != nil {
// 			str += "    " + entry.string()
// 		}
// 	}
// 	if node.GetSuccessor() != nil {
// 		str += "successor: " + node.GetSuccessor().string()
// 	} else {
// 		str += "successor: nil"
// 	}

// 	if node.GetPredecessor() != nil {
// 		str += "predecessor: " + node.GetPredecessor().string()
// 	} else {
// 		str += "predecessor: nil"
// 	}
// 	return str
// }

// // Returns the hash value in []byte based an ipAddr
// // Takes in a string, use sha1 to hash it to generate 160 bit hash value
// // mod by 2^160 so that the value's range is [0, 2^m - 1]
// func hash(ipAddr string) []byte {
// 	h := sha1.New()
// 	h.Write([]byte(ipAddr))

// 	idInt := big.NewInt(0)
// 	idInt.SetBytes(h.Sum(nil)) // Sum() returns []byte, convert it into BigInt

// 	maxVal := big.NewInt(0)
// 	maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil) // calculate 2^m
// 	idInt.Mod(idInt, maxVal)                            // mod id to make it to be [0, 2^m - 1]
// 	if idInt.Cmp(big.NewInt(0)) == 0 {
// 		return []byte{0}
// 	}
// 	return idInt.Bytes()
// }

// func (node *Node) string() string {
// 	return "Node ---> " + "id: " + hex.EncodeToString(node.ID) + " ip: " + node.IP + "\n"
// }

// func (node *Node) openConn() (*rpc.Client, error) {
// 	conn, err := net.DialTimeout("tcp", node.GetIP(), 1*time.Second)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return rpc.NewClient(conn), nil
// }

// // RemoteGetID launches an RPC call to node
// func (node *Node) RemoteGetSuccessor() (succ *Node) {
// 	client, err := node.openConn()
// 	if err != nil {

// 		return nil
// 	}
// 	defer client.Close()

// 	args := GetSuccArgs{}
// 	var reply GetSuccReply
// 	DPrintf("calling....")
// 	err = client.Call("ChordServer.GetSucc", args, &reply)
// 	DPrintf("finish calling....")
// 	if err != nil {
// 		return nil
// 	}

// 	return &reply.N
// }

// // RemoteGetID launches an RPC call to node
// func (node *Node) RemoteGetPredecessor() (pred *Node) {
// 	client, err := node.openConn()
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	defer client.Close()

// 	args := GetPredArgs{}
// 	var reply GetPredReply
// 	err = client.Call("ChordServer.GetPredecessor", args, &reply)
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	return &reply.N
// }

// // RemoteGetID launches an RPC call to node
// func (node *Node) RemoteGetID() ([]byte, error) {
// 	client, err := node.openConn()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer client.Close()

// 	args := GetInfoArgs{}
// 	var reply GetInfoReply
// 	err = client.Call("ChordServer.GetID", args, &reply)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil, nil
// }

// func (node *Node) RemoteGetInfo() (me *Node, succ *Node, pred *Node) {
// 	client, err := node.openConn()
// 	if err != nil {
// 		return nil, nil, nil
// 	}
// 	defer client.Close()

// 	args := GetInfoArgs{}
// 	var reply GetInfoReply
// 	err = client.Call("ChordServer.GetInfo", args, &reply)
// 	if err != nil {
// 		return nil, nil, nil
// 	}
// 	DPrintf("%v", reply.Me)
// 	DPrintf("%v", reply.Succ)
// 	DPrintf("%v", reply.Pred)
// 	return &reply.Me, &reply.Succ, &reply.Pred
// }

// func (node *Node) RemoteClosestPreFinger(id []byte) *Node {
// 	client, err := node.openConn()
// 	if err != nil {
// 		return nil
// 	}
// 	defer client.Close()

// 	args := ClosestPreNodeArgs{}
// 	var reply ClosestPreNodeReply
// 	err = client.Call("ChordServer.ClosestPreFinger", args, &reply)
// 	if err != nil {
// 		return nil
// 	}
// 	return &reply.ClosestPreNode
// }
