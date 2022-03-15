package srp

import (
	"errors"
	"math/big"
)

func Create(grp *Group, U, p []byte) (s []byte, v []byte, err error) {
	// Generate s random salt. Default to 32 bytes.
	//   <salt> = random()
	s, err = randomBytes(32)
	if err != nil {
		return
	}

	x := getx(grp, s, U, p)
	v = cmputev(grp, x)
	return
}

func InitiateHandshake(grp *Group) (a []byte, A []byte, err error) {
	// Create a random secret "a" value
	//   a = random()
	a, err = randomBytes(32)
	if err != nil {
		return nil, nil, err
	}

	A = getA(grp, a)
	return
}

func CompleteHandshake(grp *Group, A, a, s, U, p, B []byte) (S []byte, K []byte, err error) {
	// "B" cannot be zero
	if isZero(B) {
		return nil, nil, errors.New("\"B\" value is zero. Aborting handshake")
	}

	u := getu(grp, A, B)

	// "u" cannot be zero
	if u.Cmp(big.NewInt(0)) == 0 {
		return nil, nil, errors.New("\"u\" value is zero. Aborting handshake")
	}

	x := getx(grp, s, U, p)
	k := getk(grp)
	S = getClientS(grp, k, x, new(big.Int).SetBytes(B), new(big.Int).SetBytes(a), u)
	K = getK(grp, S)
	return
}

// Compute the secret "x" value
func getx(grp *Group, s, U, p []byte) *big.Int {
	// x = SHA(<salt> | SHA(<username> | ":" | <raw password>))
	x := new(big.Int).SetBytes(grp.Hash(s, grp.Hash(U, []byte{':'}, p)))
	return x
}

// Compute the verifier.
func cmputev(grp *Group, x *big.Int) []byte {
	// <verifier> = v = g^x % N
	return new(big.Int).Mod(new(big.Int).Exp(grp.g, x, grp.n), grp.n).Bytes()
}

// Compute "A" based on "a"
func getA(grp *Group, a []byte) (A []byte) {
	// A = g^a % N
	//    let I = g^a
	I := new(big.Int).Exp(grp.g, new(big.Int).SetBytes(a), grp.n)
	A = grp.Pad(new(big.Int).Mod(I, grp.n).Bytes())
	return
}

// Compute the pseudo-session key, "S"
func getClientS(grp *Group, k, x, B, a, u *big.Int) (S []byte) {
	// S = (B - kg^x) ^ (a + ux) % N
	//
	//  let j = kg^x
	//  	l = (B - j),
	//      r = (a + ux)
	//
	//  ... so that S = l ^ r
	j := new(big.Int).Mul(k, new(big.Int).Exp(grp.g, x, grp.n))
	l := new(big.Int).Sub(B, j)
	r := new(big.Int).Add(a, new(big.Int).Mul(u, x))
	S = grp.Pad(new(big.Int).Mod(new(big.Int).Exp(l, r, grp.n), grp.n).Bytes())
	return
}
