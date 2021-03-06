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
	"context"

	"github.com/go-redis/redis"
)

type Redis struct {
	config *Config
	redis.UniversalClient
}

func newRedis(config *Config) *Redis {
	opts := &redis.UniversalOptions{
		Addrs:              config.Addrs,
		DB:                 config.DB,
		Password:           config.Password,
		MaxRetries:         config.MaxRetries,
		MinRetryBackoff:    config.MinRetryBackoff,
		MaxRetryBackoff:    config.MaxRetryBackoff,
		DialTimeout:        config.DialTimeout,
		ReadTimeout:        config.ReadTimeout,
		WriteTimeout:       config.WriteTimeout,
		PoolSize:           config.PoolSize,
		MinIdleConns:       config.MinIdleConns,
		MaxConnAge:         config.MaxConnAge,
		PoolTimeout:        config.PoolTimeout,
		IdleTimeout:        config.IdleTimeout,
		IdleCheckFrequency: config.IdleCheckFrequency,
		TLSConfig:          config.TLSConfig,
		MaxRedirects:       config.MaxRedirects,
		ReadOnly:           config.ReadOnly,
		RouteByLatency:     config.RouteByLatency,
		RouteRandomly:      config.RouteRandomly,
		MasterName:         config.MasterName,
	}
	client := redis.NewUniversalClient(opts)
	for _, p := range config.processes {

		client.WrapProcess(p)
	}
	for _, p := range config.pipelineProcesses {
		client.WrapProcessPipeline(p)
	}
	return &Redis{config, client}
}

func (r *Redis) ClusterClient() *redis.ClusterClient {
	if c, ok := r.UniversalClient.(*redis.ClusterClient); ok {
		return c
	}
	return nil
}

func (r *Redis) FailoverOrSimpleClient() *redis.Client {
	if c, ok := r.UniversalClient.(*redis.Client); ok {
		return c
	}
	return nil
}

func (r *Redis) WithContext(ctx context.Context) *Redis {
	r2 := *r
	if r.config.EnableTrace {
		if c := r.ClusterClient(); c != nil {
			r2.UniversalClient = c.WithContext(ctx)
		}
		if c := r.FailoverOrSimpleClient(); c != nil {
			r2.UniversalClient = c.WithContext(ctx)
		}
	}
	r2.UniversalClient.WrapProcess(traceProcess(ctx))
	r2.UniversalClient.WrapProcessPipeline(tracePipelineProcess(ctx))
	return &r2
}
