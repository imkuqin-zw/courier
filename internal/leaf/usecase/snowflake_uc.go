package usecase

import (
	"context"

	pbLeaf "github.com/imkuqin-zw/courier/api/leaf"
)

type SnowflakeUC struct {
	pbLeaf.UnimplementedSnowflakeServer
}

func (s *SnowflakeUC) FetchNext(ctx context.Context, req *pbLeaf.SnowflakeFetchNextReq) (*pbLeaf.SnowflakeFetchNextRes, error) {
	return &pbLeaf.SnowflakeFetchNextRes{Seq: 0}, nil
}

func NewSnowflakeUC() *SnowflakeUC {
	return &SnowflakeUC{}
}
