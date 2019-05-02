package chord

// chord_node_test.go should test functions in chord_node.go

import (
	"bytes"
	"testing"
)

func TestReverseHash(t *testing.T) {
	ipAddress := "127.0.0.1"
	startPortNumer := 5000

	testIps := reverseHash(numBits, ipAddress, startPortNumer)

	for i := 0; i < len(testIps); i++ {
		got := hash(testIps[i])
		if !bytes.Equal(got, intToByteArray(i)) {
			t.Errorf("ReverseHash(%s) got %d, want %d", testIps[i], got, i)
		}
	}

	DPrintf("The corresponding IDs/IPs are:\n")
	for i := 0; i < len(testIps); i++ {
		DPrintf("[%d]: %s", i, testIps[i])
	}
}
