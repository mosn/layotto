package in_memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/sequencer"
)

func TestNew(t *testing.T) {
	s := NewInMemorySequencer()

	assert.NotNil(t, s)
}

func TestInit(t *testing.T) {
	s := NewInMemorySequencer()
	assert.NotNil(t, s)

	err := s.Init(sequencer.Configuration{})
	assert.NoError(t, err)
}

func TestGetNextId(t *testing.T) {
	s := NewInMemorySequencer()
	assert.NotNil(t, s)

	err := s.Init(sequencer.Configuration{})
	assert.NoError(t, err)

	var resp *sequencer.GetNextIdResponse
	resp, err = s.GetNextId(&sequencer.GetNextIdRequest{Key: "666"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.NextId)

	resp, err = s.GetNextId(&sequencer.GetNextIdRequest{Key: "666"})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.NextId)

	resp, err = s.GetNextId(&sequencer.GetNextIdRequest{Key: "777"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.NextId)
}

func TestGetSegment(t *testing.T) {
	s := NewInMemorySequencer()
	assert.NotNil(t, s)

	err := s.Init(sequencer.Configuration{})
	assert.NoError(t, err)

	var resp *sequencer.GetSegmentResponse
	var res bool
	res, resp, err = s.GetSegment(&sequencer.GetSegmentRequest{Key: "666", Size: 5})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.From)
	assert.Equal(t, int64(5), resp.To)
	assert.True(t, res)

	res, resp, err = s.GetSegment(&sequencer.GetSegmentRequest{Key: "666", Size: 5})
	assert.NoError(t, err)
	assert.Equal(t, int64(6), resp.From)
	assert.Equal(t, int64(10), resp.To)
	assert.True(t, res)

	res, resp, err = s.GetSegment(&sequencer.GetSegmentRequest{Key: "777", Size: 5})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.From)
	assert.Equal(t, int64(5), resp.To)
	assert.True(t, res)

}
