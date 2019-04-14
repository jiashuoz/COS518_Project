package chord

import (
	"math/big"
)

// intToBytes takes a int64 and convert it to []byte, big endian
func intToBytes(x int64) []byte {
	if x == 0 {
		return []byte{0}
	}
	return big.NewInt(x).Bytes()
}

// add takes one number in bytes and second number in int64, return the result in bytes
// does not work with negative numbers
func addBytesInt64(numberInBytes []byte, addend int64) []byte {
	addend1 := big.NewInt(0).SetBytes(numberInBytes)
	addend2 := big.NewInt(addend)
	return addend1.Add(addend1, addend2).Bytes()
}

func addBytesBigint(numberInBytes []byte, addend *big.Int) []byte {
	addend1 := big.NewInt(0).SetBytes(numberInBytes)
	return addend.Add(addend, addend1).Bytes()
}
