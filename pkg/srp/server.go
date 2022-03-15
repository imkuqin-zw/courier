package srp

import (
	"errors"
	"math/big"
)

func HandShake(grp *Group, A, v []byte) (B, K []byte, err error) {
	// "A" cannot be zero
	if isZero(A) {
		err = errors.New("Server found \"A\" to be zero. Aborting handshake")
		return
	}

	// Create a random secret "b"
	b, err := randomBytes(32)
	if err != nil {
		return
	}

	bigIntv := new(big.Int).SetBytes(v)
	bigIntb := new(big.Int).SetBytes(b)
	bigIntA := new(big.Int).SetBytes(A)

	k := getk(grp)

	B = getB(grp, k, bigIntv, bigIntb)

	u := getu(grp, A, B)

	S := getServerS(grp, bigIntv, bigIntA, bigIntb, u)

	K = getK(grp, S)
	return
}

// Compute a value "B" based on "b"
func getB(grp *Group, k, v, b *big.Int) (B []byte) {
	//   B = (kv + g^b) % N
	gModPowB := new(big.Int).Exp(grp.g, b, grp.n)

	kMulV := new(big.Int).Mul(k, v)

	l := new(big.Int).Add(kMulV, gModPowB)

	B = grp.Pad(new(big.Int).Mod(l, grp.n).Bytes())

	return
}

// Compute the pseudo-session key, "S"
func getServerS(grp *Group, v, A, b, u *big.Int) []byte {
	// S = (A * v^u) ^ b % N
	// 	   let I = A * v^u
	I := new(big.Int).Mul(A, new(big.Int).Exp(v, u, grp.n))
	S := grp.Pad(new(big.Int).Mod(new(big.Int).Exp(I, b, grp.n), grp.n).Bytes())
	return S
}
