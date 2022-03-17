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

package main

import (
	"github.com/imkuqin-zw/courier/internal/leaf"
	"github.com/imkuqin-zw/courier/pkg/app"
	"github.com/imkuqin-zw/courier/pkg/config"
)

func main() {
	app.Init(app.WithAppName("leaf"))
	start()
	app.WaitSignals()
}

func start() {
	snowflakeEnable := len(config.Get("dubbo.provider.services.SnowflakeUC").Bytes()) > 0
	segmentEnable := len(config.Get("dubbo.provider.services.SegmentUC").Bytes()) > 0
	if snowflakeEnable && segmentEnable {
		app.Start(app.WithProvidersFactory(leaf.NewAllProviderServices))
	} else if snowflakeEnable {
		app.Start(app.WithProvidersFactory(leaf.NewSnowflakeProviderServices))
	} else if segmentEnable {
		app.Start(app.WithProvidersFactory(leaf.NewSegmentProviderServices))
	}
}
