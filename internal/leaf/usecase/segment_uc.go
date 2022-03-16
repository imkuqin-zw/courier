package usecase

import (
	"context"

	pbLeaf "github.com/imkuqin-zw/courier/api/leaf"
	"github.com/imkuqin-zw/courier/internal/leaf/repository"
)

type SegmentUC struct {
	repo repository.SegmentRepo
	pbLeaf.UnimplementedSegmentServer
}

func (s *SegmentUC) FetchNext(ctx context.Context, req *pbLeaf.SegmentFetchNextReq) (*pbLeaf.SegmentFetchNextRes, error) {
	segment, err := s.repo.GetSegmentByID(ctx, req.Tag)
	if err != nil {
		return nil, err
	}
	seq, err := segment.FetchNextSeq(ctx)
	if err != nil {
		return nil, err
	}
	return &pbLeaf.SegmentFetchNextRes{Seq: seq}, nil
}

func NewSegmentUC(repo repository.SegmentRepo) *SegmentUC {
	return &SegmentUC{repo: repo}
}
