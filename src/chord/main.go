package chord

var Servers = make(map[string]*ChordServer)

// rpc simulator
func ChangeServer(ipAddr string) *ChordServer {
	return Servers[ipAddr]
}

