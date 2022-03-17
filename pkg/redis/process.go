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
	"fmt"
	"strings"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Processes func(func(redis.Cmder) error) func(redis.Cmder) error

type PipelineProcesses func(func([]redis.Cmder) error) func([]redis.Cmder) error

func traceProcess(ctx context.Context) Processes {
	return func(next func(redis.Cmder) error) func(redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			span, _ := opentracing.StartSpanFromContext(ctx, strings.ToUpper(cmd.Name()))
			ext.SpanKind.Set(span, "client")
			ext.DBType.Set(span, "redis")
			ext.DBStatement.Set(span, fmt.Sprintf("%v", cmd.Args()))
			defer span.Finish()

			return next(cmd)
		}
	}
}

func tracePipelineProcess(ctx context.Context) PipelineProcesses {
	return func(next func([]redis.Cmder) error) func([]redis.Cmder) error {
		return func(cmds []redis.Cmder) error {
			pipelineSpan, ctx := opentracing.StartSpanFromContext(ctx, "(pipeline)")
			ext.DBType.Set(pipelineSpan, "redis")
			ext.SpanKind.Set(pipelineSpan, "client")

			for i := len(cmds); i > 0; i-- {
				cmdName := strings.ToUpper(cmds[i-1].Name())
				if cmdName == "" {
					cmdName = "(empty command)"
				}
				span, _ := opentracing.StartSpanFromContext(ctx, cmdName)
				ext.DBStatement.Set(span, fmt.Sprintf("%v", cmds[i-1].Args()))
				span.Finish()
			}

			defer pipelineSpan.Finish()

			return next(cmds)
		}
	}
}

func slowLoggerProcess(slowThreshold time.Duration) Processes {
	return func(next func(redis.Cmder) error) func(redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			begin := time.Now()
			err := next(cmd)
			cost := time.Since(begin)
			if slowThreshold < cost {
				logger.Warnf("redis slow: cost[%s], cmd:[%s]", cost, fmt.Sprintf("%v", cmd.Args()))
			}
			return err
		}
	}
}

func slowLoggerPipelineProcess(slowThreshold time.Duration) PipelineProcesses {
	return func(next func([]redis.Cmder) error) func([]redis.Cmder) error {
		return func(cmds []redis.Cmder) error {
			begin := time.Now()

			if err := next(cmds); err != nil {
				return err
			}

			cost := time.Since(begin)
			if slowThreshold < cost {
				cmdsStr := make([]string, len(cmds))
				for i := len(cmds); i > 0; i-- {
					cmdsStr[i] = fmt.Sprintf("%v", cmds[i].Args())
				}
				logger.Warnf("redis slow: cost[%s], cmd:[%s]", cost, fmt.Sprintf("%v", cmdsStr))
			}
			return nil
		}
	}
}
