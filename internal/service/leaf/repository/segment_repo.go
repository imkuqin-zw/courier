package repository

import (
	"context"
	"runtime"
	"sync"

	"github.com/imkuqin-zw/courier/internal/service/leaf/domain"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Segment struct {
	ID        uint   `gorm:"primarykey"`
	Tag       string `gorm:"uniqueIndex"`
	Step      uint64
	MaxSeq    uint64
	Describe  string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt soft_delete.DeletedAt
}

type SegmentRepo struct {
	segments    [2]sync.Map
	initialized sync.Map
	preLoading  sync.Map
	replacing   sync.Map
	db          *gorm.DB
}

func (r *SegmentRepo) GetSegmentByID(ctx context.Context, ID string) (*domain.Segment, error) {
	v, ok := r.segments[0].Load(ID)
	if ok {
		return v.(*domain.Segment), nil
	}

	_, b := r.initialized.LoadOrStore(ID, struct{}{})
	if b {
		for {
			if v, ok := r.segments[0].Load(ID); ok {
				return v.(*domain.Segment), nil
			}
			runtime.Gosched()
		}
	}
	segment, err := r.fetchSegment(ctx, ID)
	if err != nil {
		return nil, err
	}
	r.segments[0].Store(ID, segment)
	return segment, nil
}

func (r *SegmentRepo) FetchNextSegment(ctx context.Context, ID string, maxID uint64) (*domain.Segment, error) {
	_, b := r.replacing.LoadOrStore(ID, struct{}{})
	if b {
		for {
			if v, ok := r.segments[0].Load(ID); ok {
				return v.(*domain.Segment), nil
			}
			runtime.Gosched()
		}
	}

	v, ok := r.segments[0].Load(ID)
	if ok && v.(domain.Segment).Max > maxID {
		return v.(*domain.Segment), nil
	}

	if v, ok := r.segments[1].Load(ID); ok {
		r.segments[0].Store(ID, v)
		r.segments[1].Delete(ID)
		return v.(*domain.Segment), nil
	}

	segment, err := r.fetchSegment(ctx, ID)
	if err != nil {
		return nil, err
	}

	r.replacing.Delete(ID)
	return nil, nil
}

func (r *SegmentRepo) SaveSegment(ctx context.Context, segment *domain.Segment) error {
	return nil
}

func (r *SegmentRepo) fetchSegment(ctx context.Context, ID string) (*domain.Segment, error) {
	var m Segment
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("tag = ?", ID).UpdateColumn("max_seq", gorm.Expr("max_seq + step")).Error
		if err != nil {
			return err
		}

		if err := tx.Where("tag = ?", ID).Take(&m).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	segment := domain.NewSegment(m.Tag, m.MaxSeq, m.MaxSeq-m.Step, r)
	return segment, nil
}
