package chord

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// RPC is used to handle client RPC requests.
type RPC struct {
	mu       sync.Mutex
	chord    *ChordServer
	listener net.Listener
	rpcBase  *rpc.Server
	// running     bool
	// errChan     chan error
}

func (rpcServer *RPC) getAddr() net.Addr {
	return rpcServer.listener.Addr()
}

// Local start method with more control over address server listens on. Used
// for testing.
func rpcRun(chord *ChordServer) (*RPC, error) {
	rpcServer := &RPC{}
	rpcServer.chord = chord
	// rpcs.errChan = make(chan error, 1)
	rpcServer.rpcBase = rpc.NewServer()
	var err error
	if err := rpcServer.rpcBase.Register(rpcServer); err != nil {
		checkError(err)
		return nil, fmt.Errorf("rpc registration failed: %s", err)
	}

	rpcServer.listener, err = net.Listen("tcp", chord.GetIP())
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
type FindSuccessorArgs struct{ Id []byte }

// FindSuccessorReply holds reply to FindSuccessor.
type FindSuccessorReply struct{ N Node }

// FindSuccessor is an RPC call, returns the successor node based on Id
func (rpcServer *RPC) FindSuccessor(args FindSuccessorArgs, reply *FindSuccessorReply) error {
	DPrintf("rpc FindSuccessor args: %v", args.Id)
	reply.N = rpcServer.chord.FindSuccessor(args.Id)
	return nil
}

// FindClosestNodeArgs holds arguments for FindClosestNode.
type FindClosestNodeArgs struct {
	Id []byte
}

// FindClosestNodeReply holds reply to FindClosestReply.
type FindClosestNodeReply struct {
	N Node
}

// FindClosestNode is an RPC call, calls underlying equivalent function,
// finds the closest node to Id from the Chord instance on this server.
func (rpcServer *RPC) FindClosestNode(args FindClosestNodeArgs, reply *FindClosestNodeReply) error {
	// DPrintf("server id(%d) received FindClosestNode RPC call", rpcServer.chord.GetID())
	tempN := rpcServer.chord.FindClosestNode(args.Id)
	reply.N = tempN
	// DPrintf("server id(%d) RPC result is: "+tempN.String(), rpcServer.chord.GetID())
	return nil
}

// GetSuccessorArgs holds arguments for GetSuccessor.
type GetSuccessorArgs struct{}

// GetSuccessorReply holds reply to GetSuccessor.
type GetSuccessorReply struct{ N Node }

// GetSuccessor is an PRC call, returns the Successor of the Chord instance on this server.
func (rpcServer *RPC) GetSuccessor(args GetSuccessorArgs, reply *GetSuccessorReply) error {
	reply.N = rpcServer.chord.fingerTable[0]
	return nil
}

// GetPredecessorArgs holds arguments for GetPredecessor.
type GetPredecessorArgs struct{}

// GetPredecessorReply holds reply to GetPredecessor.
type GetPredecessorReply struct{ N Node }

// GetPredecessor is an PRC call, returns the Predecessor of the Chord instance on this server.
func (rpcServer *RPC) GetPredecessor(args GetPredecessorArgs, reply *GetPredecessorReply) error {
	reply.N = rpcServer.chord.predecessor
	return nil
}

// NotifyArgs holds arguments for GetPredecessor.
type NotifyArgs struct{ N Node }

// NotifyReply holds reply to GetPredecessor.
type NotifyReply struct{}

// Notify is an RPC call, a remote node notifies us it might be our predecessor
func (rpcServer *RPC) Notify(args NotifyArgs, reply *NotifyArgs) error {
	err := rpcServer.chord.Notify(args.N)
	return err
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
