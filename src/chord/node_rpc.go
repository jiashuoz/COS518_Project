package chord

import (
	"net"
	"net/rpc"
	"time"
)

func (node *Node) openConn() (*rpc.Client, error) {
	conn, err := net.DialTimeout("tcp", node.IP, 5*time.Second)
	if err != nil {
		checkError(err)
		return nil, err
	}
	return rpc.NewClient(conn), nil
}

// FindSuccessorRPC sends RPC call to remote node
func (node *Node) FindSuccessorRPC(id []byte) (Node, error) {
	client, err := node.openConn()
	if err != nil {
		checkError(err)
		return Node{}, err
	}
	defer client.Close()

	args := FindSuccessorArgs{}
	var reply FindSuccessorReply
	err = client.Call("RPC.FindSuccessor", args, &reply)
	if err != nil {
		checkError(err)
		return Node{}, err
	}
	return reply.N, nil
}

// GetSuccessorRPC sends RPC call to remote node
func (node *Node) GetSuccessorRPC() (Node, error) {
	client, err := node.openConn()
	if err != nil {
		checkError(err)
		return Node{}, err
	}
	defer client.Close()

	args := GetSuccessorArgs{}
	var reply GetSuccessorReply
	err = client.Call("RPC.GetSuccessor", args, &reply)
	if err != nil {
		checkError(err)
		return Node{}, err
	}
	return reply.N, nil
}

// FindClosestNodeRPC sends RPC call to remote node
func (node *Node) FindClosestNodeRPC(id []byte) (Node, error) {
	client, err := node.openConn()
	if err != nil {
		checkError(err)
		return Node{}, err
	}
	defer client.Close()

	args := FindClosestNodeArgs{}
	var reply FindClosestNodeReply
	args.ID = id
	err = client.Call("RPC.FindClosestNode", args, &reply)
	if err != nil {
		checkError(err)
		return Node{}, err
	}
	return reply.N, nil
}

// RemoteGetPred returns the predecessor of the specified node.
// func (n *Node) RemoteGetPred() (*Node, error) {
// 	client, err := n.openConn()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer client.Close()

// 	args := GetPredArgs{}
// 	var reply GetPredReply
// 	err = client.Call("Server.GetPred", args, &reply)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &reply.N, nil
// }

// RemoteGetSucc returns the successor of the specified node.
// THE SERVER LOCKS THE CHORD MUTEX TO ACCESS THE SUCCESSOR
// func (n *Node) RemoteGetSucc() (*Node, error) {
// 	client, err := n.openConn()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer client.Close()

// 	args := new(GetSuccArgs)
// 	var reply GetSuccReply
// 	err = client.Call("Server.GetSucc", args, &reply)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &reply.N, nil
// }

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
