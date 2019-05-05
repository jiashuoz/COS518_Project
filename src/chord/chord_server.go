package chord

import (
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

const (
	// DefaultStabilizeInterval is the interval that this server will start the stabilize process
	DefaultStabilizeInterval = 1 * time.Second

	// DefaultFixFingerInterval is the interval that this server will repeat fixing its finger table
	DefaultFixFingerInterval = 1 * time.Second
)

// ChordServer is a single ChordServer
type ChordServer struct {
	node        Node
	fingerTable []Node // a table of FingerEntry pointer
	predecessor Node   // previous node on the identifier circle
	tracer      Tracer //tracer to trace node hops and latency
	running     bool   // true if running, false if stopped
	stopChan    chan bool
	sync.RWMutex
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

	chord.stopChan = make(chan bool)
	chord.running = true

	chord.routineGroup.Add(1)
	go func() {
		defer chord.routineGroup.Done()
		for chord.Running() {
			select {
			case _, ok := <-chord.stopChan:
				if !ok {
					return
				}
			case <-time.NewTimer(DefaultStabilizeInterval).C:
				err := chord.Stabilize()
				if err != nil {
					checkError(err)
					log.Printf("Stabilize error")
				}
			}
		}
	}()

	chord.routineGroup.Add(1)
	go func() {
		defer chord.routineGroup.Done()
		for chord.Running() {
			select {
			case _, ok := <-chord.stopChan:
				if !ok {
					return
				}
			case <-time.NewTimer(DefaultStabilizeInterval).C:
				chord.fixFingers()
			}
		}
	}()

	return nil
}

// SetRunning sets the running state of chord
func (chord *ChordServer) SetRunning(running bool) {
	chord.Lock()
	defer chord.Unlock()
	chord.running = running
}

func (chord *ChordServer) Running() bool {
	chord.RLock()
	defer chord.RUnlock()
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
	pred := chord.FindPredecessor(id)
	result, _ := pred.GetSuccessorRPC()
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
	chord.RLock()
	defer chord.RUnlock()

	fingerTable := chord.fingerTable
	for i := numBits - 1; i >= 0; i-- {
		if fingerTable[i].ID != nil && between(fingerTable[i].ID, chord.node.ID, id) {
			return fingerTable[i]
		}
	}
	return chord.node
}

// Join adds chord to ring based on an existing node
func (chord *ChordServer) Join(node Node) {
	chord.Lock()
	defer chord.Unlock()
	if node.ID == nil { // the only node in the ring
		chord.fingerTable[0] = chord.node
		chord.predecessor = chord.node
	} else {

		DPrintf("Calling..")
		chord.fingerTable[0], _ = node.FindSuccessorRPC(chord.node.ID)
		chord.predecessor, _ = chord.fingerTable[0].GetPredecessorRPC()
		DPrintf("Hangs..")
	}

}

// Stabilize periodically verify's chord's immediate successor
func (chord *ChordServer) Stabilize() error {
	succ := chord.GetSuccessor()
	x, _ := succ.GetPredecessorRPC()

	if between(x.ID, chord.GetID(), succ.ID) {
		chord.Lock()
		chord.fingerTable[0] = x
		chord.Unlock()
	}

	succ = chord.GetSuccessor()
	DPrintf("%v notifies %v", chord.GetID(), succ.ID)
	succ.NotifyRPC(chord.GetNode())
	return nil
}

// Notify tells chord, node thinks it might be chord's predecessor.
func (chord *ChordServer) Notify(node Node) error {
	//TODO: lock here since it is changing chordServer property
	chord.RLock()
	if idsEqual(chord.node.ID, node.ID) {
		chord.RUnlock()
		return nil
	}
	chord.RUnlock()
	chord.Lock()
	defer chord.Unlock()

	DPrintf("%v %v %v", node.ID, chord.predecessor.ID, chord.node.ID)
	// node is the only node in the ring
	if idsEqual(chord.node.ID, chord.predecessor.ID) {
		chord.predecessor = node
		chord.fingerTable[0] = node
	} else if between(node.ID, chord.predecessor.ID, chord.node.ID) {
		chord.predecessor = node
	}
	return nil
}

// periodically fresh finger table entries
func (chord *ChordServer) fixFingers() {
	i := rand.Intn(numBits-1) + 1
	fingerStart := chord.fingerStart(i)
	chord.fingerTable[i] = chord.FindSuccessor(fingerStart)
}

func (chord *ChordServer) fingerStart(i int) []byte {
	currID := new(big.Int).SetBytes(chord.node.ID)
	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
	start := new(big.Int).Add(currID, offset)
	return start.Bytes()
}

// GetNode returns chord's network information.
func (chord *ChordServer) GetNode() Node {
	chord.RLock()
	defer chord.RUnlock()
	return chord.node
}

// GetID return's chord's identifier.
func (chord *ChordServer) GetID() []byte {
	chord.RLock()
	defer chord.RUnlock()
	return chord.node.ID
}

// GetIP return's chord's identifier.
func (chord *ChordServer) GetIP() string {
	chord.RLock()
	defer chord.RUnlock()
	return chord.node.IP
}

// GetSuccessor returns chord's successor
func (chord *ChordServer) GetSuccessor() Node {
	chord.RLock()
	defer chord.RUnlock()
	return chord.fingerTable[0]
}

func (chord *ChordServer) String(printFingerTable bool) string {
	chord.RLock()
	defer chord.RUnlock()
	str := chord.node.String()
	if !printFingerTable {
		return str
	}

	str += "Finger table: \n"
	str += "ith | start | successor\n"
	for i := 0; i < numBits; i++ {
		currID := new(big.Int).SetBytes(chord.node.ID)
		offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		start := new(big.Int).Add(currID, offset)
		successor := chord.fingerTable[i].ID
		str += fmt.Sprintf("%d   | %d     | %d\n", i, start, successor)
	}
	str += fmt.Sprintf("Predecessor: %v", chord.predecessor.ID)
	return str
}
