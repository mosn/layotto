package sequencer

import (
	"context"
	"mosn.io/layotto/components/sequencer"
)

func GetNextIdFromCache(ctx context.Context, store sequencer.Store, req *sequencer.GetNextIdRequest) (bool, int64, error) {
	// TODO
	//	1. check cache
	//	2. load cache with lock
	//	2.1. return if not support
	//	2.2. add segment into cache
	//	2.3. return result
	return false, 0, nil
}
