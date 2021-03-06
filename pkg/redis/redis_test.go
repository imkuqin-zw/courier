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

package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	conf := &Config{
		Addrs: []string{"localhost:6379"},
	}
	redisCli := conf.Build()
	err := redisCli.Ping()
	require.NoError(t, err, "ping redislock")
	err = redisCli.Set("_test_redis", 1, time.Second)
	require.NoError(t, err, "set")
	val, err := redisCli.Get("_test_redis").Int()
	require.NoError(t, err, "get")
	require.Equal(t, 1, val, "get")
}
