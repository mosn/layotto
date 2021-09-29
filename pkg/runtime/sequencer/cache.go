package sequencer

import (
	"context"
	"errors"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/utils"
	"sync"
)

const defaultSize = 1000
const defaultLimit = 700

// DoubleBuffer is double segment id buffer.
// There are two buffers in DoubleBuffer: InUseBuffer is in use, BackUpBuffer is a backup buffer.
// Their default capacity is 1000. When the InUseBuffer usage exceeds 30%, the BackUpBuffer will be initialized.
// When InUseBuffer is used up, swap them.
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

	d.InUseBuffer = buffer

	return nil
}

//getId next id
func (d *DoubleBuffer) getId() (int64, error) {

	if d.InUseBuffer == nil {
		return 0, errors.New("[DoubleBuffer] Get error: InUseBuffer nil ")
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	//check swap
	if d.InUseBuffer.from > d.InUseBuffer.to {
		err := d.swap()
		if err != nil {
			return 0, err
		}
	}
	next := d.InUseBuffer.from
	d.InUseBuffer.from++

	//when InUseBuffer id more than half used, initialize BackUpBuffer.
	//equal make sure only one thread enter
	if d.InUseBuffer.to-d.InUseBuffer.from == defaultLimit {
		utils.GoWithRecover(func() {

			//lock to avoid data race
			d.lock.Lock()
			defer d.lock.Unlock()
			//one case: here and swap are performed simultaneously,
			//swap fast,no check will cover old BackUpBuffer
			if d.BackUpBuffer != nil {
				return
			}
			buffer, err := d.getNewBuffer()
			if err != nil {
				log.DefaultLogger.Errorf("[DoubleBuffer] [getNewBuffer] error: %v", err)
				return
			}

			d.BackUpBuffer = buffer

		}, nil)
	}

	return next, nil
}

//swap InUseBuffer and BackUpBuffer, must be locked
func (d *DoubleBuffer) swap() error {

	//check again
	if d.BackUpBuffer == nil {
		buffer, err := d.getNewBuffer()
		if err != nil {
			return err
		}
		d.BackUpBuffer = buffer
	}

	d.InUseBuffer = d.BackUpBuffer
	d.BackUpBuffer = nil
	return nil
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
		return nil, errors.New("[DoubleBuffer] unSupport Segment id")
	}
	return &Buffer{
		from: result.From,
		to:   result.To,
	}, nil
}

var BufferCatch = map[string]*DoubleBuffer{}

//common lock is enough ？
var rwLock sync.RWMutex

func GetNextIdFromCache(ctx context.Context, store sequencer.Store, req *sequencer.GetNextIdRequest) (bool, int64, error) {

	//check support
	support, _, _ := store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  req.Key,
		Size: 0,
	})

	if !support {
		return false, 0, nil
	}

	var d *DoubleBuffer
	var err error
	if _, ok := BufferCatch[req.Key]; !ok {
		d, err = getDoubleBufferInWL(req.Key, store)
	} else {
		d, err = getDoubleBufferInRL(req.Key)
	}

	if err != nil {
		return true, 0, err
	}

	id, err := d.getId()

	if err != nil {
		return true, 0, err
	}

	return true, id, nil
}

//DoubleBuffer for this key not exist
func getDoubleBufferInWL(key string, store sequencer.Store) (*DoubleBuffer, error) {
	d := NewDoubleBuffer(key, store)
	rwLock.Lock()
	defer rwLock.Unlock()
	//double check
	if _, ok := BufferCatch[key]; ok {
		return BufferCatch[key], nil
	}
	err := d.init()
	if err != nil {
		return nil, err
	}
	BufferCatch[key] = d
	return d, nil
}

// DoubleBuffer for this key  exist
func getDoubleBufferInRL(key string) (*DoubleBuffer, error) {
	/*	rwLock.RLock()
		defer rwLock.RUnlock()*/
	return BufferCatch[key], nil
}
