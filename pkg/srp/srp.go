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

import "math/big"

// Compute the SRP-6a version of the multiplier parameter "k"
func getK(grp *Group, S []byte) []byte {
	return grp.Hash(S)
}

// Compute "u"
func getu(grp *Group, A, B []byte) *big.Int {
	return new(big.Int).SetBytes(grp.Hash(A, B))
}

// Compute the SRP-6a version of the multiplier parameter "k"
func getk(grp *Group) *big.Int {
	k := grp.Hash(grp.Pad(grp.n.Bytes()), grp.Pad(grp.g.Bytes()))
	return new(big.Int).SetBytes(k)
}

func ComputeM1(grp *Group, A, B, K []byte) []byte {
	return grp.Hash(A, B, K)
}

func ComputeM2(grp *Group, A, M1, K []byte) []byte {
	return grp.Hash(A, M1, K)
}
