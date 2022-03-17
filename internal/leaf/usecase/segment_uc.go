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
