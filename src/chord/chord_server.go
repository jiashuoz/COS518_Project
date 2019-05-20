package chord

import (
	"fmt"
	// "log"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

const (
	// DefaultStabilizeInterval is the interval that this server will start the stabilize process
	DefaultStabilizeInterval = 600 * time.Millisecond

	// DefaultFixFingerInterval is the interval that this server will repeat fixing its finger table
	DefaultFixFingerInterval = 500 * time.Millisecond
)

// ChordServer is a single ChordServer
type ChordServer struct {
	node         Node
	fingerTable  []Node // a table of FingerEntry pointer
	predecessor  Node   // previous node on the identifier circle
	tracer       Tracer //tracer to trace node hops and latency
	running      bool   // true if running, false if stopped
	stopChan     chan bool
	rwmu         sync.RWMutex
	routineGroup sync.WaitGroup
}

// MakeServer returns a pointer to a  server
func MakeServer(ip string) *ChordServer {
	server := &ChordServer{}
	server.node = MakeNode(ip)
	server.fingerTable = make([]Node, numBits)
	server.tracer = MakeTracer()
	server.running = false
	DPrintf("Initialized ChordServer ----> ip: %v | id: %v", server.node.IP, server.node.ID)
	return server
}

func (chord *ChordServer) Start() error {
	if chord.Running() {
		return fmt.Errorf("Start() failed: %v already running", chord.GetID())
	}

	rpcRun(chord) // launch RPC service

	// chord.stopChan = make(chan bool)
	// chord.running = true

	// chord.routineGroup.Add(1)
	// go func() {
	// 	// defer chord.routineGroup.Done()
	// 	for {
	// 		select {
	// 		case _, ok := <-chord.stopChan:
	// 			if !ok {
	// 				DPrintf("stopping stabilize process...")
	// 				return
	// 			}
	// 		case <-time.NewTimer(DefaultStabilizeInterval).C:
	// 			DPrintf("stabilize...")
	// 			err := chord.Stabilize()
	// 			if err != nil {
	// 				checkError(err)
	// 				log.Printf("Stabilize error")
	// 			}
	// 		}
	// 	}
	// }()

	// chord.routineGroup.Add(1)
	go func() {
		// defer chord.routineGroup.Done()
		for {
			select {
			case _, ok := <-chord.stopChan:
				if !ok {
					DPrintf("stopping fix finger process...")
					return
				}
			case <-time.NewTimer(DefaultStabilizeInterval).C:
				DPrintf("fix fingers...")
				chord.fixFingers()
			}
		}
	}()

	return nil
}

// SetRunning sets the running state of chord
func (chord *ChordServer) SetRunning(running bool) {
	chord.rwmu.Lock()
	defer chord.rwmu.Unlock()
	chord.running = running
}

func (chord *ChordServer) Running() bool {
	chord.rwmu.Lock()
	defer chord.rwmu.Unlock()
	return chord.running
}

// Lookup returns the ip of node that stores key id
func (chord *ChordServer) Lookup(id []byte) string {
	succ := chord.FindSuccessor(id)
	return succ.IP
}

// FindSuccessor returns Node that's the successor of key id
func (chord *ChordServer) FindSuccessor(id []byte) Node {
	// chord.tracer.startTracer(chord.GetID(), id)
	DPrintf("FindSuccessor: 1")
	pred := chord.FindPredecessor(id)
	DPrintf("FindSuccessor: 2")
	result, _ := pred.GetSuccessorRPC()
	DPrintf("FindSuccessor: return")
	// chord.tracer.endTracer(result.ID)
	return result
}

// FindPredecessor returns the Node that's preceding id
func (chord *ChordServer) FindPredecessor(id []byte) Node {
	closest := chord.FindClosestNode(id)
	if idsEqual(closest.ID, chord.GetID()) {
		return closest
	}

	closestSucc, _ := closest.GetSuccessorRPC()

	// chord.tracer.traceNode(closest.ID)

	for !betweenRightInclusive(id, closest.ID, closestSucc.ID) {
		closest, _ = closest.FindClosestNodeRPC(id)
		closestSucc, _ = closest.GetSuccessorRPC()
		// chord.tracer.traceNode(closest.ID)
	}
	return closest
}

// FindClosestNode returns the closest node to id based on fingerTable
func (chord *ChordServer) FindClosestNode(id []byte) Node {
	chord.rwmu.RLock()
	fingerTable := chord.fingerTable
	chord.rwmu.RUnlock()
	for i := numBits - 1; i >= 0; i-- {
		chord.rwmu.RLock()
		if fingerTable[i].ID != nil && between(fingerTable[i].ID, chord.node.ID, id) {
			chord.rwmu.RUnlock()
			return fingerTable[i]
		}
		chord.rwmu.RUnlock()
	}
	return chord.node
}

// Join adds chord to ring based on an existing node
func (chord *ChordServer) Join(node Node) {
	chord.rwmu.Lock()
	defer chord.rwmu.Unlock()
	if node.ID == nil { // the only node in the ring
		chord.fingerTable[0] = chord.node
		chord.predecessor = chord.node
	} else {
		chord.fingerTable[0], _ = node.FindSuccessorRPC(chord.node.ID)
		chord.predecessor, _ = chord.fingerTable[0].GetPredecessorRPC()
	}

}

// Stabilize periodically verify's chord's immediate successor
func (chord *ChordServer) Stabilize() error {
	DPrintf("Stabilize: 1....")
	succ := chord.GetSuccessor()
	DPrintf("Stabilize: 1a....")
	if idsEqual(succ.ID, chord.GetID()) {
		DPrintf("Stabilize: 1b....")
		return nil
	}
	DPrintf("Stabilize: 2....")
	x, _ := succ.GetPredecessorRPC()

	DPrintf("Stabilize: 3....")
	if between(x.ID, chord.GetID(), succ.ID) {
		DPrintf("%v %v %v", x.ID, chord.GetID(), succ.ID)
		chord.rwmu.Lock()
		chord.fingerTable[0] = x
		chord.rwmu.Unlock()
	}

	DPrintf("Stabilize: 4....")
	succ = chord.GetSuccessor()
	DPrintf("Stabilize: 5....")
	succ.NotifyRPC(chord.GetNode())
	DPrintf("Stabilize: 6....")
	return nil
}

// Notify tells chord, node thinks it might be chord's predecessor.
func (chord *ChordServer) Notify(node Node) error {
	//TODO: lock here since it is changing chordServer property
	if idsEqual(chord.GetID(), node.ID) {
		chord.rwmu.RUnlock()
		return nil
	}

	// node is the only node in the ring
	if idsEqual(chord.GetID(), chord.predecessor.ID) {
		chord.rwmu.Lock()
		chord.predecessor = node
		chord.fingerTable[0] = node
		chord.rwmu.Unlock()
	} else if between(node.ID, chord.predecessor.ID, chord.node.ID) {
		chord.rwmu.Lock()
		chord.predecessor = node
		chord.rwmu.Unlock()
	}
	return nil
}

// periodically fresh finger table entries
func (chord *ChordServer) fixFingers() {
	i := rand.Intn(numBits-1) + 1
	DPrintf("finger start")
	fingerStart := chord.fingerStart(i)
	DPrintf("finger findsuccessor")
	finger := chord.FindSuccessor(fingerStart)
	DPrintf("finger set finger %v", fingerStart)
	chord.rwmu.Lock()
	chord.fingerTable[i] = finger
	chord.rwmu.Unlock()
}

func (chord *ChordServer) fingerStart(i int) []byte {
	currID := new(big.Int).SetBytes(chord.GetID())
	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
	maxVal := big.NewInt(0)
	maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil)
	start := new(big.Int).Add(currID, offset)
	start.Mod(start, maxVal)
	return start.Bytes()
}

// GetNode returns chord's network information.
func (chord *ChordServer) GetNode() Node {
	chord.rwmu.RLock()
	defer chord.rwmu.RUnlock()
	return chord.node
}

// GetID return's chord's identifier.
func (chord *ChordServer) GetID() []byte {
	chord.rwmu.RLock()
	defer chord.rwmu.RUnlock()
	return chord.node.ID
}

// GetIP return's chord's identifier.
func (chord *ChordServer) GetIP() string {
	chord.rwmu.RLock()
	defer chord.rwmu.RUnlock()
	return chord.node.IP
}

// GetSuccessor returns chord's successor
func (chord *ChordServer) GetSuccessor() Node {
	chord.rwmu.RLock()
	defer chord.rwmu.RUnlock()
	return chord.fingerTable[0]
}

// GetPredecessor returns chord's predecessor
func (chord *ChordServer) GetPredecessor() Node {
	chord.rwmu.RLock()
	defer chord.rwmu.RUnlock()
	return chord.predecessor
}

func (chord *ChordServer) String(printFingerTable bool) string {
	chord.rwmu.RLock()
	str := chord.node.String()
	chord.rwmu.RUnlock()
	if !printFingerTable {
		return str
	}

	str += "Finger table: \n"
	str += "ith | start | successor\n"
	for i := 0; i < numBits; i++ {
		currID := new(big.Int).SetBytes(chord.GetID())
		offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		maxVal := big.NewInt(0)
		maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil)
		start := new(big.Int).Add(currID, offset)
		start.Mod(start, maxVal)
		chord.rwmu.Lock()
		successor := chord.fingerTable[i].ID
		chord.rwmu.Unlock()
		str += fmt.Sprintf("%d   | %d     | %d\n", i, start, successor)
	}
	chord.rwmu.Lock()
	str += fmt.Sprintf("Predecessor: %v", chord.predecessor.ID)
	chord.rwmu.Unlock()
	return str
}
