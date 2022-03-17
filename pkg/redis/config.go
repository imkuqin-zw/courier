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
	"crypto/tls"
	"fmt"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/imkuqin-zw/courier/pkg/config"
)

const (
	defaultSlowThreshold = time.Millisecond * 250
)

type Config struct {
	// Either a single address or a seed list of host:port addresses
	// of cluster/sentinel nodes.
	Addrs []string

	// Database to be selected after connecting to the server.
	// Only single-node and failover clients.
	DB int

	// Common options.

	Password           string
	MaxRetries         int
	MinRetryBackoff    time.Duration
	MaxRetryBackoff    time.Duration
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	MaxConnAge         time.Duration
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	TLSConfig          *tls.Config

	// Only cluster clients.

	MaxRedirects   int
	ReadOnly       bool
	RouteByLatency bool
	RouteRandomly  bool

	// The sentinel master name.
	// Only failover clients.
	MasterName string

	// 慢日志阈值
	SlowThreshold time.Duration
	EnableTrace   bool

	processes         []Processes
	pipelineProcesses []PipelineProcesses
}

func StdConfig(name string) *Config {
	return RawConfig(fmt.Sprintf("%s.%s", "redis", name))
}

func RawConfig(key string) *Config {
	c := new(Config)
	if err := config.Get(key).Scan(c); err != nil {
		logger.Fatalf("fault to scan config, error: %s", err.Error())
	}
	return c
}

//WithProcesses
func (config *Config) WithProcesses(ps ...Processes) *Config {
	if config.processes == nil {
		config.processes = make([]Processes, 0, len(ps))
	}
	config.processes = append(config.processes, ps...)
	return config
}

//WithPipelineProcesses
func (config *Config) WithPipelineProcesses(ps ...PipelineProcesses) *Config {
	if config.pipelineProcesses == nil {
		config.pipelineProcesses = make([]PipelineProcesses, 0, len(ps))
	}
	config.pipelineProcesses = append(config.pipelineProcesses, ps...)
	return config
}

//Build
func (config *Config) Build() *Redis {
	if config.SlowThreshold == 0 {
		config.SlowThreshold = defaultSlowThreshold
	}
	config.WithProcesses(slowLoggerProcess(config.SlowThreshold))
	config.WithPipelineProcesses(slowLoggerPipelineProcess(config.SlowThreshold))
	return newRedis(config)
}
