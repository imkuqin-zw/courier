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
