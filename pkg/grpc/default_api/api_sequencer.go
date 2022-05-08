/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package default_api

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/messages"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (a *api) GetNextId(ctx context.Context, req *runtimev1pb.GetNextIdRequest) (*runtimev1pb.GetNextIdResponse, error) {
	// 1. validate
	if len(a.sequencers) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrSequencerStoresNotConfigured)
		log.DefaultLogger.Errorf("[runtime] [grpc.GetNextId] error: %v", err)
		return &runtimev1pb.GetNextIdResponse{}, err
	}
	if req.Key == "" {
		err := status.Errorf(codes.InvalidArgument, messages.ErrSequencerKeyEmpty, req.StoreName)
		return &runtimev1pb.GetNextIdResponse{}, err
	}
	// 2. convert
	compReq, err := GetNextIdRequest2ComponentRequest(req)
	if err != nil {
		return &runtimev1pb.GetNextIdResponse{}, err
	}
	// modify key
	compReq.Key, err = runtime_sequencer.GetModifiedSeqKey(compReq.Key, req.StoreName, a.appId)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetNextId] error: %v", err)
		return &runtimev1pb.GetNextIdResponse{}, err
	}
	// 3. find store component
	store, ok := a.sequencers[req.StoreName]
	if !ok {
		return &runtimev1pb.GetNextIdResponse{}, status.Errorf(codes.InvalidArgument, messages.ErrSequencerStoreNotFound, req.StoreName)
	}
	var next int64
	// 4. invoke component
	if compReq.Options.AutoIncrement == sequencer.WEAK {
		// WEAK
		next, err = a.getNextIdWithWeakAutoIncrement(ctx, store, compReq)
	} else {
		// STRONG
		next, err = a.getNextIdFromComponent(ctx, store, compReq)
	}
	// 5. convert response
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetNextId] error: %v", err)
		return &runtimev1pb.GetNextIdResponse{}, err
	}
	return &runtimev1pb.GetNextIdResponse{
		NextId: next,
	}, nil
}

func (a *api) getNextIdWithWeakAutoIncrement(ctx context.Context, store sequencer.Store, compReq *sequencer.GetNextIdRequest) (int64, error) {
	// 1. try to get from cache
	support, next, err := runtime_sequencer.GetNextIdFromCache(ctx, store, compReq)

	if !support {
		// 2. get from component
		return a.getNextIdFromComponent(ctx, store, compReq)
	}
	return next, err
}

func (a *api) getNextIdFromComponent(ctx context.Context, store sequencer.Store, compReq *sequencer.GetNextIdRequest) (int64, error) {
	var next int64
	resp, err := store.GetNextId(compReq)
	if err == nil {
		next = resp.NextId
	}
	return next, err
}

func GetNextIdRequest2ComponentRequest(req *runtimev1pb.GetNextIdRequest) (*sequencer.GetNextIdRequest, error) {
	result := &sequencer.GetNextIdRequest{}
	if req == nil {
		return nil, errors.New("cannot convert it since request is nil")
	}

	result.Key = req.Key
	var incrOption = sequencer.WEAK
	if req.Options != nil {
		if req.Options.Increment == runtimev1pb.SequencerOptions_WEAK {
			incrOption = sequencer.WEAK
		} else if req.Options.Increment == runtimev1pb.SequencerOptions_STRONG {
			incrOption = sequencer.STRONG
		} else {
			return nil, errors.New("options.Increment is illegal")
		}
	}
	result.Options = sequencer.SequencerOptions{AutoIncrement: incrOption}
	result.Metadata = req.Metadata
	return result, nil
}
