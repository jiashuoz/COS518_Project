package chord

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// RPCServer used to handle client RPC requests.
type RPC struct {
	mu       sync.Mutex
	chord    *ChordServer
	listener net.Listener
	rpcBase  *rpc.Server
	// running     bool
	// errChan     chan error
}

// StartRPC creates an RPCServer listening on given port.
func StartRPC(chord *ChordServer, port int) (*RPC, error) {
	return run(chord)
}

func (rpcHandler *RPC) getAddr() net.Addr {
	return rpcHandler.listener.Addr()
}

// Local start method with more control over address server listens on. Used
// for testing.
func run(chord *ChordServer) (*RPC, error) {
	rpcServer := &RPC{}
	rpcServer.chord = chord
	// rpcs.errChan = make(chan error, 1)
	rpcServer.rpcBase = rpc.NewServer()
	var err error
	if err := rpcServer.rpcBase.Register(rpcServer); err != nil {
		checkError(err)
		return nil, fmt.Errorf("rpc registration failed: %s", err)
	}

	rpcServer.listener, err = net.Listen("tcp", chord.node.IP)
	if err != nil {
		return nil, err
	}

	// Start thread for server.
	go func() {
		// rpcHandler.mu.Lock()
		// rpcHandler.running = true
		// rpcHandler.mu.Unlock()

		// Blocks until completion
		rpcServer.rpcBase.Accept(rpcServer.listener)

		// rpcs.mu.Lock()
		// rpcs.errChan <- nil
		// rpcs.running = false
		// rpcs.mu.Unlock()
	}()

	return rpcServer, nil
}

// FindSuccessorArgs holds arguments for FindSuccessor.
type FindSuccessorArgs struct{ id []byte }

// FindSuccessorReply holds reply to FindSuccessor.
type FindSuccessorReply struct{ N Node }

// FindSuccessor is an RPC call, returns the successor node based on Id
func (rpcServer *RPC) FindSuccessor(args *FindSuccessorArgs, reply *FindSuccessorReply) error {
	reply.N = rpcServer.chord.FindSuccessor(args.id)
	return nil
}

// FindClosestNodeArgs holds arguments for FindClosestNode.
type FindClosestNodeArgs struct {
	ID []byte
}

// FindClosestNodeReply holds reply to FindClosestReply.
type FindClosestNodeReply struct {
	N Node
}

// FindClosestNode is an RPC call, calls underlying equivalent function,
// finds the closest node to Id from the Chord instance on this server.
func (rpcServer *RPC) FindClosestNode(args *FindClosestNodeArgs, reply *FindClosestNodeReply) error {
	DPrintf("chord [%s]: FindClosestNode (%016x)", rpcServer.chord.GetID(), args.ID)

	tempN := rpcServer.chord.FindClosestNode(args.ID)
	reply.N = tempN

	DPrintf("ch [%s]: FindClosestNode (%016x): secceeded: %s",
		rpcServer.chord.node.String(), args.ID, tempN.String())

	return nil
}

// GetSuccessorArgs holds arguments for GetSuccessor.
type GetSuccessorArgs struct{}

// GetSuccessorReply holds reply to GetSuccessor.
type GetSuccessorReply struct{ N Node }

// GetSuccessor is an PRC call, returns the Successor of the Chord instance on this server.
func (rpcServer *RPC) GetSuccessor(args *GetSuccessorArgs, reply *GetSuccessorReply) error {
	reply.N = rpcServer.chord.fingerTable[0]
	return nil
}

// Return true if server is up. Return false otherwise.
// func (rpcHandler *RPC) isRunning() bool {
// 	rpcHandler.mu.Lock()
// 	defer rpcHandler.mu.Unlock()
// 	return rpcHandler.running
// }

// Wait blocks until server stops.
// func (rpcHandler *RPC) Wait() error {
// 	if !rpcHandler.isRunning() {
// 		return nil
// 	}

// 	return <-rpcHandler.errChan
// }

// End the server if it is running. Returns nil on success.
// func (rpcHandler *RPC) End() error {
// 	if !rpcHandler.isRunning() {
// 		return nil
// 	}

// 	if err := rpcHandler.listener.Close(); err != nil {
// 		return fmt.Errorf("error closing listener: %s", err)
// 	}

// 	return nil
// }

// GetAddr returns the address the servers is listening on.

// checkInit used at the start of every RPC to ensure that the server does
// not respond to requests until the server is fully initialized.
// func (rpcHandler *RPC) checkInit() {
// 	rpcs.initBarrier.RLock()
// 	rpcs.initBarrier.RUnlock()
// }

// // ChordLookupArgs holds arguments to ChordLookup RPC.
// type ChordLookupArgs struct {
// 	id []byte
// }

// // ChordLookupReply holds reply to ChordLookup RPC.
// type ChordLookupReply struct {
// 	ip string
// }

// // ChordLookup returns result of lookup performed by the Chord instance running
// // on this RPCServer.
// func (rpcs *RPCServer) ChordLookup(args *ChordLookupArgs, reply *ChordLookupReply) error {
// 	rpcs.checkInit()
// 	DPrintf("ch [%016x]: ChordLookup for %016x", rpcs.chordServer.node.ipAddr, args.id)
// 	result := rpcs.chordServer.LookUp(args.id)
// 	reply.ip = result
// 	return nil
// }
