package srp

import (
	"crypto/rand"
	"math/big"
)

//
// Get n random bytes. Returns a byte slice.
//
func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

//
// Returns true if a byte slice is equal to 0
//
func isZero(x []byte) bool {
	// Convert x1 from []byte -> *Int
	xBigInt := big.NewInt(0).SetBytes(x)

	// Define a 0 big int to compare to xBigInt
	zeroBigInt := big.NewInt(0)

	isZero := xBigInt.Cmp(zeroBigInt) == 0

	return isZero
}

func padTo(bytes []byte, length int) []byte {
	paddingLength := length - len(bytes)
	padding := make([]byte, paddingLength)
	return append(padding, bytes...)
}
