package chord

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// RPCServer used to handle client RPC requests.
type RPCServer struct {
	mu           sync.Mutex
	initBarrier  sync.RWMutex
	chordServer          *ChordServer
	kvServer       *KVServer	//don't have KVServer yet
	servListener net.Listener
	baseServ     *rpc.Server
	running      bool
	errChan      chan error
}

// StartRPC creates an RPCServer listening on given port.
func StartRPC(chordServer *ChordServer, kv *KVServer, port int) (*RPCServer, error) {
	return startRPC(chordServer, kv, fmt.Sprintf(":%d", port))
}

// Local start method with more control over address server listens on. Used
// for testing.
func startRPC(ch *ChordServer, kv *KVServer, addr string) (*RPCServer, error) {
	rpcs := &RPCServer{}
	rpcs.chordServer = ch
	rpcs.kvServer = kv
	rpcs.errChan = make(chan error, 1)
	rpcs.baseServ = rpc.NewServer()

	if err := rpcs.baseServ.Register(rpcs); err != nil {
		return nil, fmt.Errorf("rpc registration failed: %s", err)
	}

	var err error
	rpcs.servListener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Start thread for server.
	go func() {
		rpcs.mu.Lock()
		rpcs.running = true
		rpcs.mu.Unlock()

		// Blocks until completion
		rpcs.baseServ.Accept(rpcs.servListener)

		rpcs.mu.Lock()
		rpcs.errChan <- nil
		rpcs.running = false
		rpcs.mu.Unlock()
	}()

	return rpcs, nil
}

// Return true if server is up. Return false otherwise.
func (rpcs *RPCServer) isRunning() bool {
	rpcs.mu.Lock()
	defer rpcs.mu.Unlock()
	return rpcs.running
}

// Wait blocks until server stops.
func (rpcs *RPCServer) Wait() error {
	if !rpcs.isRunning() {
		return nil
	}

	return <-rpcs.errChan
}

// End the server if it is running. Returns nil on success.
func (rpcs *RPCServer) End() error {
	if !rpcs.isRunning() {
		return nil
	}

	if err := rpcs.servListener.Close(); err != nil {
		return fmt.Errorf("error closing listener: %s", err)
	}

	return nil
}

// GetAddr returns the address the servers is listening on.
func (rpcs *RPCServer) GetAddr() net.Addr {
	return rpcs.servListener.Addr()
}

// checkInit used at the start of every RPC to ensure that the server does
// not respond to requests until the server is fully initialized.
func (rpcs *RPCServer) checkInit() {
	rpcs.initBarrier.RLock()
	rpcs.initBarrier.RUnlock()
}

// ChordLookupArgs holds arguments to ChordLookup RPC.
type ChordLookupArgs struct {
	id []byte
}

// ChordLookupReply holds reply to ChordLookup RPC.
type ChordLookupReply struct {
	ip string
}

// ChordLookup returns result of lookup performed by the Chord instance running
// on this RPCServer.
func (rpcs *RPCServer) ChordLookup(args *ChordLookupArgs, reply *ChordLookupReply) error {
	rpcs.checkInit()
	DPrintf("ch [%016x]: ChordLookup for %016x", rpcs.chordServer.node.ipAddr, args.id)
	result := rpcs.chordServer.LookUp(args.id)
	reply.ip = result
	return nil
}
