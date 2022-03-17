// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
