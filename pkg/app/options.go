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

package app

import (
	"fmt"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common"
)

const (
	baseAppName     = "com.github.imkuqin_zw.courier"
	shutdownTimeout = time.Second * 60
)

type Options struct {
	appName              string
	disableAppCfgFile    bool
	disableAppCfgDynamic bool
	providerFactory      func() []common.RPCService
	consumerFactory      func() []common.RPCService
}

type Option func(*Options)

func WithAppName(appName string) Option {
	return func(options *Options) {
		if appName == "" {
			return
		}
		options.appName = fmt.Sprintf("%s.%s", baseAppName, appName)
	}
}

func DisableAppCfgFile() Option {
	return func(options *Options) {
		options.disableAppCfgFile = true
	}
}

func DisableAppCfgDynamic() Option {
	return func(options *Options) {
		options.disableAppCfgDynamic = true
	}
}

func WithProvidersFactory(f func() []common.RPCService) Option {
	return func(options *Options) {
		options.providerFactory = f
	}
}

func WithConsumersFactory(f func() []common.RPCService) Option {
	return func(options *Options) {
		options.consumerFactory = f
	}
}

func applyOpts(o *Options, ops ...Option) {
	for _, f := range ops {
		f(o)
	}
}
