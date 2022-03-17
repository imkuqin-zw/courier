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

// +build wireinject

package leaf

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"github.com/google/wire"
	"github.com/imkuqin-zw/courier/internal/leaf/repository"
	"github.com/imkuqin-zw/courier/internal/leaf/usecase"
	"github.com/imkuqin-zw/courier/pkg/gorm"
)

// ProviderSet is segment providers.
var providerSegmentSet = wire.NewSet(
	usecase.NewSegmentUC,
	repository.NewSegmentRepo,
	gorm.New,
)

func newSegmentProvider(
	segmentUC *usecase.SegmentUC,
) []common.RPCService {
	return []common.RPCService{
		segmentUC,
	}
}

func NewSegmentProviderServices() []common.RPCService {
	panic(wire.Build(providerSegmentSet, newSegmentProvider))
}

// ProviderSet is snowflake providers.
var providerSnowflakeSet = wire.NewSet(
	usecase.NewSnowflakeUC,
)

func newSnowflakeProvider(
	snowFlakeUC *usecase.SnowflakeUC,
) []common.RPCService {
	return []common.RPCService{
		snowFlakeUC,
	}
}

func NewSnowflakeProviderServices() []common.RPCService {
	panic(wire.Build(providerSnowflakeSet, newSnowflakeProvider))
}

// ProviderSet is all providers.
var providerAllSet = wire.NewSet(
	providerSnowflakeSet,
	providerSegmentSet,
)

func newAllProvider(
	snowFlakeUC *usecase.SnowflakeUC,
	segmentUC *usecase.SegmentUC,
) []common.RPCService {
	return []common.RPCService{
		snowFlakeUC,
		segmentUC,
	}
}

func NewAllProviderServices() []common.RPCService {
	panic(wire.Build(providerAllSet, newAllProvider))
}
