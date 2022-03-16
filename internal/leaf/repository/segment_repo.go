package repository

import (
	"context"
	"runtime"
	"sync"
	"time"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/imkuqin-zw/courier/internal/leaf/domain"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type SegmentRepo interface {
	domain.SegmentRepo
	GetSegmentByID(ctx context.Context, ID string) (*domain.Segment, error)
}

type Segment struct {
	ID        uint   `gorm:"primarykey"`
	Tag       string `gorm:"uniqueIndex"`
	Step      uint32
	MaxSeq    uint64
	Describe  string
	CreatedAt int64
	UpdatedAt int64
	DeletedAt soft_delete.DeletedAt
}

type segmentStep struct {
	Step          uint32
	LastFetchedAt time.Time
}

type stepManager struct {
	minStableTime time.Duration
	maxStableTime time.Duration
	mu            sync.RWMutex
	steps         map[string]*segmentStep
}

func (sm *stepManager) Add(tag string, step uint32) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.steps[tag] = &segmentStep{
		Step:          step,
		LastFetchedAt: time.Now(),
	}
}

func (sm *stepManager) GetStep(tag string) uint32 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	step, ok := sm.steps[tag]
	if ok {
		return step.Step
	}
	return 0
}

func (sm *stepManager) Update(tag string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	step := sm.steps[tag]
	now := time.Now()
	duration := now.Sub(step.LastFetchedAt)
	switch {
	case duration < sm.minStableTime:
		step.Step <<= 1
	case duration > sm.maxStableTime && step.Step > 2:
		step.Step >>= 1
	default:
	}
	step.LastFetchedAt = now
}

type segmentRepo struct {
	segments    [2]sync.Map
	stepManager *stepManager
	initialized sync.Map
	preLoading  sync.Map
	replacing   sync.Map
	db          *gorm.DB
}

func NewSegmentRepo(db *gorm.DB) SegmentRepo {
	return &segmentRepo{
		stepManager: &stepManager{
			minStableTime: time.Minute * 15,
			maxStableTime: time.Minute * 30,
			steps:         make(map[string]*segmentStep),
		},
		db: db,
	}
}

func (r *segmentRepo) GetSegmentByID(ctx context.Context, ID string) (*domain.Segment, error) {
	v, ok := r.segments[0].Load(ID)
	if ok {
		return v.(*domain.Segment), nil
	}

	_, b := r.initialized.LoadOrStore(ID, struct{}{})
	if b {
		if b {
			return r.loopLoad(ID), nil
		}
	}
	segment, err := r.fetchSegment(ctx, ID)
	if err != nil {
		return nil, err
	}
	r.segments[0].Store(ID, segment)
	return segment, nil
}

func (r *segmentRepo) FetchNextSegment(ctx context.Context, ID string, maxSeq uint64) (*domain.Segment, error) {
	_, b := r.replacing.LoadOrStore(ID, struct{}{})
	if b {
		return r.loopLoad(ID), nil
	}
	defer r.replacing.Delete(ID)

	v, ok := r.segments[0].Load(ID)
	if ok && v.(domain.Segment).Max > maxSeq {
		return v.(*domain.Segment), nil
	}

	if segment, ok := r.replaceSegment(ID); ok {
		return segment, nil
	}

	if err := r.syncPreloadSegment(ctx, ID, maxSeq); err != nil {
		return nil, err
	}

	return r.loopReplace(ID), nil
}

func (r *segmentRepo) SaveSegment(ctx context.Context, ID string, seq, max uint64, step uint32) {
	if float32(max-seq)/float32(step) <= 0.9 {
		r.asyncPreloadSegment(ctx, ID, max)
	}
}

func (r *segmentRepo) preloadSegment(ctx context.Context, ID string, maxSeq uint64) error {
	if _, ok := r.segments[1].Load(ID); ok {
		return nil
	}
	if v, ok := r.segments[0].Load(ID); ok && maxSeq < v.(domain.Segment).Max {
		return nil
	}
	segment, err := r.fetchSegmentWithStep(ctx, ID)
	if err != nil {
		return err
	}
	r.segments[1].Store(ID, segment)
	return nil
}

func (r *segmentRepo) syncPreloadSegment(ctx context.Context, ID string, maxSeq uint64) error {
	_, b := r.preLoading.LoadOrStore(ID, struct{}{})
	if b {
		return nil
	}
	defer r.preLoading.Delete(ID)
	return r.preloadSegment(ctx, ID, maxSeq)
}

func (r *segmentRepo) asyncPreloadSegment(ctx context.Context, ID string, maxSeq uint64) {
	_, b := r.preLoading.LoadOrStore(ID, struct{}{})
	if b {
		return
	}
	defer r.preLoading.Delete(ID)
	go func() {
		if err := r.preloadSegment(ctx, ID, maxSeq); err != nil {
			logger.Errorf("fault to preload segment, err: %s", err)
		}
	}()
	return
}

func (r *segmentRepo) loopLoad(ID string) *domain.Segment {
	for {
		if v, ok := r.segments[0].Load(ID); ok {
			return v.(*domain.Segment)
		}
		runtime.Gosched()
	}
}

func (r *segmentRepo) loopReplace(ID string) *domain.Segment {
	for {
		if segment, ok := r.replaceSegment(ID); ok {
			return segment
		}
		runtime.Gosched()
	}
}

func (r *segmentRepo) replaceSegment(ID string) (*domain.Segment, bool) {
	if v, ok := r.segments[1].Load(ID); ok {
		r.segments[0].Store(ID, v)
		r.segments[1].Delete(ID)
		return v.(*domain.Segment), true
	}
	return nil, false
}

func (r *segmentRepo) fetchSegment(ctx context.Context, ID string) (*domain.Segment, error) {
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
	r.stepManager.Add(ID, m.Step)
	segment := domain.NewSegment(m.Tag, m.MaxSeq, m.Step, m.MaxSeq-uint64(m.Step), r)
	return segment, nil
}

func (r *segmentRepo) fetchSegmentWithStep(ctx context.Context, ID string) (*domain.Segment, error) {
	var m Segment
	step := r.stepManager.GetStep(ID)
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("tag = ?", ID).UpdateColumn("max_seq", gorm.Expr("max_seq + ?", step)).Error
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
	r.stepManager.Update(ID)
	segment := domain.NewSegment(m.Tag, m.MaxSeq, step, m.MaxSeq-uint64(m.Step), r)
	return segment, nil
}
