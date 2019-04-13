package chord

var Servers = make(map[string]*Server)

func ChangeServer(ipAddr string) *Server {
	return Servers[ipAddr]
}
