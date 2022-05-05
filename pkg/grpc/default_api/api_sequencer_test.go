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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/sequencer"
	mock_sequencer "mosn.io/layotto/pkg/mock/components/sequencer"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestGetNextId(t *testing.T) {
	t.Run("sequencers not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "abc",
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = Sequencer store is not configured", err.Error())
	})

	t.Run("seq key empty", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "abc",
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Key is empty in sequencer store abc", err.Error())
	})

	t.Run("sequencer store not found", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "abc",
			Key:       "next key",
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Sequencer store abc not found", err.Error())
	})

	t.Run("auto increment is strong", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		mockSequencerStore.EXPECT().GetNextId(gomock.Any()).
			DoAndReturn(func(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
				assert.Equal(t, "sequencer|||next key", req.Key)
				assert.Equal(t, sequencer.STRONG, req.Options.AutoIncrement)
				return &sequencer.GetNextIdResponse{
					NextId: 10,
				}, nil
			})
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "mock",
			Key:       "next key",
			Options: &runtimev1pb.SequencerOptions{
				Increment: runtimev1pb.SequencerOptions_STRONG,
			},
		}
		rsp, err := api.GetNextId(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, int64(10), rsp.NextId)
	})

	t.Run("net error", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		mockSequencerStore.EXPECT().GetNextId(gomock.Any()).Return(nil, fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "mock",
			Key:       "next key",
			Options: &runtimev1pb.SequencerOptions{
				Increment: runtimev1pb.SequencerOptions_STRONG,
			},
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "net error", err.Error())
	})
}
