package in_memory

import (
	"go.uber.org/atomic"
	"mosn.io/layotto/components/sequencer"
	"sync"
)

type InMemorySequencer struct {
	data *sync.Map
}

func NewInMemorySequencer() *InMemorySequencer {
	return &InMemorySequencer{
		data: &sync.Map{},
	}
}

func (s *InMemorySequencer) Init(_ sequencer.Configuration) error {
	return nil
}

func (s *InMemorySequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	seed, ok := s.data.Load(req.Key)
	if !ok {
		seed, _ = s.data.LoadOrStore(req.Key, &atomic.Int64{})
	}

	nextId := seed.(*atomic.Int64).Inc()
	return &sequencer.GetNextIdResponse{NextId: nextId}, nil

}

func (s *InMemorySequencer) GetSegment(req *sequencer.GetSegmentRequest) (bool, *sequencer.GetSegmentResponse, error) {
	seed, ok := s.data.Load(req.Key)
	if !ok {
		seed, _ = s.data.LoadOrStore(req.Key, &atomic.Int64{})
	}

	res := seed.(*atomic.Int64).Add(int64(req.Size))
	return true, &sequencer.GetSegmentResponse{
		From: res - int64(req.Size) + 1,
		To:   res,
	}, nil
}
