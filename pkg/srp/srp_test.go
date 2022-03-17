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
	"bytes"
	"fmt"
	"testing"
)

var (
	grp = Group2048()
	U   = []byte("1141137429")
	p   = []byte("zhangwei")
)

func TestCreate(t *testing.T) {
	s, v, err := Create(grp, U, p) //客户端
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("s: %x", s))
	fmt.Println(fmt.Sprintf("v: %x", v))
	a, A, err := InitiateHandshake(grp) //客户端
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("a: %x", a))
	fmt.Println(fmt.Sprintf("A: %x", A))
	B, serverK, err := HandShake(grp, A, v) //服务端，通过u找到v
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("B: %x", B))
	//fmt.Println(fmt.Sprintf("serverS: %x", S))
	fmt.Println(fmt.Sprintf("serverK: %x", serverK))
	_, clientK, err := CompleteHandshake(grp, A, a, s, U, p, B)
	if err != nil {
		t.Fatal("Error in CompleteHandshake()")
	}
	//-------------------------
	//fmt.Println(fmt.Sprintf("clientS: %x", clientS))
	fmt.Println(fmt.Sprintf("clientK: %x", clientK))
	clientM1 := ComputeM1(grp, A, B, clientK)

	serverM1 := ComputeM1(grp, A, B, serverK)
	fmt.Println(fmt.Sprintf("clientM1: %x", clientM1))
	if !bytes.Equal(clientM1, serverM1) {
		t.Fatal("M1 not success")
	}

	serverM2 := ComputeM2(grp, A, serverM1, serverK)

	clientM2 := ComputeM2(grp, A, clientM1, clientK)
	fmt.Println(fmt.Sprintf("clientM2: %x", clientM2))
	if !bytes.Equal(serverM2, clientM2) {
		t.Fatal("M2 not success")
	}
}

/**
client					server
			U, A
		 ---------->

			s, B
		 <----------
S K						S K


*/
