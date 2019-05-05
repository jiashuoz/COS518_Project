package chord

var Servers = make(map[string]*ChordServer)

const numBits = 3

// rpc simulator
func ChangeServer(ipAddr string) *ChordServer {
	return Servers[ipAddr]
}
