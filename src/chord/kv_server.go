package chord

import (
	"sync"
)

// KVServer represents node in distributed key value storage system.
type KVServer struct {
	mu    sync.Mutex
	chord    *ChordServer
	state map[string]string
}