package converter

import (
	"github.com/dapr/components-contrib/state"
	"github.com/stretchr/testify/assert"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"testing"
)

func TestGetResponse2GetStateResponse(t *testing.T) {
	resp := GetResponse2GetStateResponse(&state.GetResponse{
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: make(map[string]string),
	})
	assert.Equal(t, resp.Data, []byte("v"))
	assert.Equal(t, resp.Etag, "")
	assert.True(t, len(resp.Metadata) == 0)
}

func TestGetResponse2BulkStateItem(t *testing.T) {
	itm := GetResponse2BulkStateItem(&state.GetResponse{
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: make(map[string]string),
	}, "key")
	assert.Equal(t, itm.Key, "key")
	assert.Equal(t, itm.Data, []byte("v"))
	assert.Equal(t, itm.Etag, "")
	assert.Equal(t, itm.Error, "")
	assert.True(t, len(itm.Metadata) == 0)
}

func TestBulkGetResponse2BulkStateItem(t *testing.T) {
	itm := BulkGetResponse2BulkStateItem(&state.BulkGetResponse{
		Key:      "key",
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: nil,
		Error:    "",
	})
	assert.Equal(t, itm.Key, "key")
	assert.Equal(t, itm.Data, []byte("v"))
	assert.Equal(t, itm.Etag, "")
	assert.Equal(t, itm.Error, "")
	assert.True(t, len(itm.Metadata) == 0)
}

func TestStateItem2SetRequest(t *testing.T) {
	req := StateItem2SetRequest(&runtimev1pb.StateItem{
		Key:      "",
		Value:    []byte("v"),
		Etag:     nil,
		Metadata: nil,
		Options: &runtimev1pb.StateOptions{
			Concurrency: runtimev1pb.StateOptions_CONCURRENCY_UNSPECIFIED,
			Consistency: runtimev1pb.StateOptions_CONSISTENCY_UNSPECIFIED,
		},
	}, "appid||key")
	assert.Equal(t, req.Key, "appid||key")
	assert.Equal(t, req.Value, []byte("v"))
	assert.Nil(t, req.ETag)
	assert.Equal(t, req.Options.Consistency, "")
	assert.Equal(t, req.Options.Concurrency, "")
}

func TestDeleteStateRequest2DeleteRequest(t *testing.T) {
	req := DeleteStateRequest2DeleteRequest(&runtimev1pb.DeleteStateRequest{
		StoreName: "redis",
		Key:       "",
		Etag:      nil,
		Options: &runtimev1pb.StateOptions{
			Concurrency: runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE,
			Consistency: runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL,
		},
		Metadata: nil,
	}, "appid||key")
	assert.Equal(t, req.Key, "appid||key")
	assert.Nil(t, req.ETag)
	assert.Equal(t, req.Options.Consistency, "eventual")
	assert.Equal(t, req.Options.Concurrency, "last-write")
}

func TestStateItem2DeleteRequest(t *testing.T) {
	req := StateItem2DeleteRequest(&runtimev1pb.StateItem{
		Key:      "",
		Value:    []byte("v"),
		Etag:     nil,
		Metadata: nil,
		Options: &runtimev1pb.StateOptions{
			Concurrency: runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE,
			Consistency: runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL,
		},
	}, "appid||key")
	assert.Equal(t, req.Key, "appid||key")
	assert.Nil(t, req.ETag)
	assert.Nil(t, req.ETag)
	assert.Equal(t, req.Options.Consistency, "eventual")
	assert.Equal(t, req.Options.Concurrency, "last-write")
}
