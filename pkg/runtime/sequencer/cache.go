package sequencer

import (
	"context"
	"errors"
	"mosn.io/layotto/components/sequencer"
	"sync"
)

const defaultSize = 1000
const defaultLimit = 500

type DoubleBuffer struct {
	Key          string
	Size         int
	InUseBuffer  *Buffer
	BackUpBuffer *Buffer
	lock         sync.Mutex
	Store        sequencer.Store
}

type Buffer struct {
	from int64
	to   int64
}

func NewDoubleBuffer(key string, store sequencer.Store) *DoubleBuffer {

	d := &DoubleBuffer{
		Key:   key,
		Size:  defaultSize,
		Store: store,
	}

	return d
}

//init double buffer
func (d *DoubleBuffer) init() error {

	buffer, err := d.getNewBuffer()
	if err != nil {
		return err
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	d.InUseBuffer = buffer
	return nil
}

//get next id
func (d *DoubleBuffer) get() (int64, error) {

	if d.InUseBuffer == nil {
		return 0, errors.New("fail")
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	next := d.InUseBuffer.from
	d.InUseBuffer.from++

	//check swap
	if d.InUseBuffer.from > d.InUseBuffer.to {
		d.swap()
	}

	//equal make sure only one thread enter
	if d.InUseBuffer.to-d.InUseBuffer.from == defaultLimit {
		go func() {
			//how to handle this err?
			buffer, _ := d.getNewBuffer()
			d.BackUpBuffer = buffer

		}()
	}

	return next, nil
}

//swap InUseBuffer and BackUpBuffer, must be locked
func (d *DoubleBuffer) swap() {
	d.InUseBuffer = d.BackUpBuffer
	d.BackUpBuffer = nil
}

//getNewBuffer return a new segment
func (d *DoubleBuffer) getNewBuffer() (*Buffer, error) {
	support, result, err := d.Store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  d.Key,
		Size: d.Size,
	})
	if err != nil {
		return nil, err
	}
	if !support {
		return nil, nil
	}
	return &Buffer{
		from: result.From,
		to:   result.To,
	}, nil
}

var Catch sync.Map

func GetNextIdFromCache(ctx context.Context, store sequencer.Store, req *sequencer.GetNextIdRequest) (bool, int64, error) {

	//check support
	support, _, _ := store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  req.Key,
		Size: 0,
	})

	if !support {
		return false, 0, nil
	}

	actual, ok := Catch.LoadOrStore(req.Key, NewDoubleBuffer(req.Key, store))
	d := actual.(DoubleBuffer)
	//load fail, should init buffer
	if !ok {
		err := d.init()
		if err != nil {
			return true, 0, err
		}
	}

	get, err := d.get()

	if err != nil {
		return true, 0, err
	}
	return true, get, nil
}
