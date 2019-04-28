package chord

import (
	"net"
	"net/rpc"
	"time"
)

func (n *Node) openConn() (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", n.IP(), 5*time.Second)
	if err != nil {
		return nil, err
	}
	return rpc.NewClient(conn), nil
}

// RemoteLookup performs Lookup RPC on remote node.
func (n *Node) RemoteLookup(id []byte) (ip string, e error) {
	DPrintf("Remote lookup %v -> (%v)", n.ID(), id)
	client, err := n.openConn()
	if err != nil {
		return "", err
	}
	defer client.Close()

	args := &LookupArgs{n.ID()}
	var reply LookupReply
	err = client.Call("Server.Lookup", args, &reply)
	if err != nil {
		return "", err
	}

	return reply.IP, nil
}

// RemoteGetPred returns the predecessor of the specified node.
func (n *Node) RemoteGetPred() (*Node, error) {
	client, err := n.openConn()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	args := GetPredArgs{}
	var reply GetPredReply
	err = client.Call("Server.GetPred", args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply.N, nil
}

// RemoteGetSucc returns the successor of the specified node.
// THE SERVER LOCKS THE CHORD MUTEX TO ACCESS THE SUCCESSOR
func (n *Node) RemoteGetSucc() (*Node, error) {
	client, err := n.openConn()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	args := new(GetSuccArgs)
	var reply GetSuccReply
	err = client.Call("Server.GetSucc", args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply.N, nil
}

// // RemoteFindClosestNode find the closest node from n to hash identifier h.
// // Returns the closest known Node and the Chord instance for that node if the
// // Node is the actual Node that is responsible for h. If the returned Node is
// // not responsible for h then the Chord return is nil.
// func (n *Node) RemoteFindClosestNode(h UHash) (Node, Chord, error) {
// 	DPrintf("ch [%s]: RemoteFindClosestNode (%016x): opening connection", n.String(), h)
// 	client, err := n.openConn()
// 	if err != nil {
// 		return Node{}, Chord{}, err
// 	}
// 	defer client.Close()

// 	args := &FindClosestArgs{h}
// 	var reply FindClosestReply
// 	DPrintf("ch [%s]: RemoteFindClosestNode (%016x): making RPC", n.String(), h)
// 	err = client.Call("RPCServer.FindClosestNode", args, &reply)
// 	if err != nil {
// 		DPrintf("ch [%s]: RemoteFindClosest (%016x): RPC failed: %s", n.String(), h, err)
// 		return Node{}, Chord{}, err
// 	}

// 	DPrintf("ch [%s]: RemoteFindClosest (%016x): RPC succeeded: "+
// 		"{Done:%v, N:%s, Ch:%s}", n.String(), h, reply.Done,
// 		reply.N.String(), reply.ChFields.N.String())

// 	if reply.Done {
// 		return reply.N, *deserializeChord(&reply.ChFields), nil
// 	}
// 	return reply.N, Chord{}, nil
// }
