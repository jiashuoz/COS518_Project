package chord

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAdd(t *testing.T) {
	number1 := big.NewInt(10).Bytes()
	number2 := int64(20)
	result := big.NewInt(0).SetBytes(add(number1, number2))
	fmt.Println(result)
}

func TestIntToBytes(*testing.T) {
	fmt.Println(intToBytes(10))
	// doesn't work with negative number, but that's ok
	fmt.Println(intToBytes(-10))
}
