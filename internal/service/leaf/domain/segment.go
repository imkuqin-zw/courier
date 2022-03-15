package domain

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/atomic"
)

type SegmentBuffer struct {
	RWMutex       *sync.RWMutex
	Seq           *atomic.Int64
	NextReady     bool
	InitOk        bool
	ThreadRunning *atomic.Bool
	Step          int
	MinStep       int
	UpdatedAt     int64
}

type SegmentRepo interface {
	FetchNextSegment(ctx context.Context, ID string, maxID uint64) (*Segment, error)
	SaveSegment(ctx context.Context, segment *Segment) error
}

type Segment struct {
	ID   string
	Seq  *atomic.Uint64
	Max  uint64
	Repo SegmentRepo
}

func NewSegment(ID string, seq uint64, max uint64, repo SegmentRepo) *Segment {
	return &Segment{ID: ID, Seq: atomic.NewUint64(seq), Max: max, Repo: repo}
}

func (s *Segment) FetchNextSeq(ctx context.Context) (uint64, error) {
	seq := s.Seq.Inc()
	if seq >= s.Max {
		s, err := s.Repo.FetchNextSegment(ctx, s.ID, s.Max)
		if err != nil {
			return 0, err
		}
		select {
		case <-ctx.Done():
			return 0, errors.New("request timeout")
		default:
			return s.FetchNextSeq(ctx)
		}
	}
	_ = s.Repo.SaveSegment(ctx, s)
	return seq, nil
}
