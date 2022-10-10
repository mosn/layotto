// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// CODE ATTRIBUTION: https://github.com/dapr/go-sdk
// Modified the import package to use layotto's pb
// We use same sdk code with Dapr's for state API because we want to keep compatible with Dapr state API
package client

import (
	"context"
	"time"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/pkg/errors"
)

const (
	// StateConsistencyUndefined is the undefined value for state consistency.
	StateConsistencyUndefined StateConsistency = 0
	// StateConsistencyEventual represents eventual state consistency value.
	StateConsistencyEventual StateConsistency = 1
	// StateConsistencyStrong represents strong state consistency value.
	StateConsistencyStrong StateConsistency = 2

	// StateConcurrencyUndefined is the undefined value for state concurrency.
	StateConcurrencyUndefined StateConcurrency = 0
	// StateConcurrencyFirstWrite represents first write concurrency value.
	StateConcurrencyFirstWrite StateConcurrency = 1
	// StateConcurrencyLastWrite represents last write concurrency value.
	StateConcurrencyLastWrite StateConcurrency = 2

	// StateOperationTypeUndefined is the undefined value for state operation type.
	StateOperationTypeUndefined OperationType = 0
	// StateOperationTypeUpsert represents upsert operation type value.
	StateOperationTypeUpsert OperationType = 1
	// StateOperationTypeDelete represents delete operation type value.
	StateOperationTypeDelete OperationType = 2
	// UndefinedType represents undefined type value
	UndefinedType = "undefined"
)

type (
	// StateConsistency is the consistency enum type.
	StateConsistency int
	// StateConcurrency is the concurrency enum type.
	StateConcurrency int
	// OperationType is the operation enum type.
	OperationType int
)

// GetPBConsistency get consistency pb value
func (s StateConsistency) GetPBConsistency() runtimev1pb.StateOptions_StateConsistency {
	return runtimev1pb.StateOptions_StateConsistency(s)
}

// GetPBConcurrency get concurrency pb value
func (s StateConcurrency) GetPBConcurrency() runtimev1pb.StateOptions_StateConcurrency {
	return runtimev1pb.StateOptions_StateConcurrency(s)
}

// String returns the string value of the OperationType.
func (o OperationType) String() string {
	names := [...]string{
		UndefinedType,
		"upsert",
		"delete",
	}
	if o < StateOperationTypeUpsert || o > StateOperationTypeDelete {
		return UndefinedType
	}

	return names[o]
}

// String returns the string value of the StateConsistency.
func (s StateConsistency) String() string {
	names := [...]string{
		UndefinedType,
		"strong",
		"eventual",
	}
	if s < StateConsistencyStrong || s > StateConsistencyEventual {
		return UndefinedType
	}

	return names[s]
}

// String returns the string value of the StateConcurrency.
func (s StateConcurrency) String() string {
	names := [...]string{
		UndefinedType,
		"first-write",
		"last-write",
	}
	if s < StateConcurrencyFirstWrite || s > StateConcurrencyLastWrite {
		return UndefinedType
	}

	return names[s]
}

var (
	stateOptionDefault = &runtimev1pb.StateOptions{
		Concurrency: runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE,
		Consistency: runtimev1pb.StateOptions_CONSISTENCY_STRONG,
	}
)

// StateOperation is a collection of StateItems with a store name.
type StateOperation struct {
	Type OperationType
	Item *SetStateItem
}

// StateItem represents a single state item.
type StateItem struct {
	Key      string
	Value    []byte
	Etag     string
	Metadata map[string]string
}

// BulkStateItem represents a single state item.
type BulkStateItem struct {
	Key      string
	Value    []byte
	Etag     string
	Metadata map[string]string
	Error    string
}

// SetStateItem represents a single state to be persisted.
type SetStateItem struct {
	Key      string
	Value    []byte
	Etag     *ETag
	Metadata map[string]string
	Options  *StateOptions
}

// DeleteStateItem represents a single state to be deleted.
type DeleteStateItem SetStateItem

// ETag represents an versioned record information
type ETag struct {
	Value string
}

// StateOptions represents the state store persistence policy.
type StateOptions struct {
	Concurrency StateConcurrency
	Consistency StateConsistency
}

// StateOption StateOptions's function type
type StateOption func(*StateOptions)

// WithConcurrency set StateOptions's Concurrency
func WithConcurrency(concurrency StateConcurrency) StateOption {
	return func(so *StateOptions) {
		so.Concurrency = concurrency
	}
}

// WithConsistency set StateOptions's consistency
func WithConsistency(consistency StateConsistency) StateOption {
	return func(so *StateOptions) {
		so.Consistency = consistency
	}
}

func toProtoSaveStateItem(si *SetStateItem) (item *runtimev1pb.StateItem) {
	s := &runtimev1pb.StateItem{
		Key:      si.Key,
		Metadata: si.Metadata,
		Value:    si.Value,
		Options:  toProtoStateOptions(si.Options),
	}

	if si.Etag != nil {
		s.Etag = &runtimev1pb.Etag{
			Value: si.Etag.Value,
		}
	}

	return s
}

func toProtoStateOptions(so *StateOptions) (opts *runtimev1pb.StateOptions) {
	if so == nil {
		return copyStateOptionDefaultPB()
	}
	return &runtimev1pb.StateOptions{
		Concurrency: runtimev1pb.StateOptions_StateConcurrency(so.Concurrency),
		Consistency: runtimev1pb.StateOptions_StateConsistency(so.Consistency),
	}
}

func copyStateOptionDefaultPB() *runtimev1pb.StateOptions {
	return &runtimev1pb.StateOptions{
		Concurrency: stateOptionDefault.GetConcurrency(),
		Consistency: stateOptionDefault.GetConsistency(),
	}
}

func copyStateOptionDefault() *StateOptions {
	return &StateOptions{
		Concurrency: StateConcurrency(stateOptionDefault.GetConcurrency()),
		Consistency: StateConsistency(stateOptionDefault.GetConsistency()),
	}
}

func toProtoDuration(d time.Duration) *duration.Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &duration.Duration{
		Seconds: secs,
		Nanos:   int32(nanos),
	}
}

// ExecuteStateTransaction provides way to execute multiple operations on a specified store.
func (c *GRPCClient) ExecuteStateTransaction(ctx context.Context, storeName string, meta map[string]string, ops []*StateOperation) error {
	// 1. parameter validation
	if storeName == "" {
		return errors.New("nil storeName")
	}
	if len(ops) == 0 {
		return nil
	}
	// 2. prepare request
	items := make([]*runtimev1pb.TransactionalStateOperation, 0)
	for _, op := range ops {
		item := &runtimev1pb.TransactionalStateOperation{
			OperationType: op.Type.String(),
			Request:       toProtoSaveStateItem(op.Item),
		}
		items = append(items, item)
	}

	req := &runtimev1pb.ExecuteStateTransactionRequest{
		Metadata:   meta,
		StoreName:  storeName,
		Operations: items,
	}
	// 3. send request
	_, err := c.protoClient.ExecuteStateTransaction(ctx, req)
	if err != nil {
		return errors.Wrap(err, "error executing state transaction")
	}
	return nil
}

// SaveState saves the raw data into store, default options: strong, last-write
func (c *GRPCClient) SaveState(ctx context.Context, storeName, key string, data []byte, so ...StateOption) error {
	var stateOptions = new(StateOptions)
	for _, o := range so {
		o(stateOptions)
	}
	if len(so) == 0 {
		stateOptions = copyStateOptionDefault()
	}
	item := &SetStateItem{Key: key, Value: data, Options: stateOptions}
	return c.SaveBulkState(ctx, storeName, item)
}

// SaveBulkState saves the multiple state item to store.
func (c *GRPCClient) SaveBulkState(ctx context.Context, storeName string, items ...*SetStateItem) error {
	if storeName == "" {
		return errors.New("nil store")
	}
	if items == nil {
		return errors.New("nil item")
	}

	req := &runtimev1pb.SaveStateRequest{
		StoreName: storeName,
		States:    make([]*runtimev1pb.StateItem, 0),
	}

	for _, si := range items {
		item := toProtoSaveStateItem(si)
		req.States = append(req.States, item)
	}

	_, err := c.protoClient.SaveState(ctx, req)
	if err != nil {
		return errors.Wrap(err, "error saving state")
	}
	return nil
}

// GetBulkState retrieves state for multiple keys from specific store.
func (c *GRPCClient) GetBulkState(ctx context.Context, storeName string, keys []string, meta map[string]string, parallelism int32) ([]*BulkStateItem, error) {
	if storeName == "" {
		return nil, errors.New("nil store")
	}
	if len(keys) == 0 {
		return nil, errors.New("keys required")
	}
	items := make([]*BulkStateItem, 0)

	req := &runtimev1pb.GetBulkStateRequest{
		StoreName:   storeName,
		Keys:        keys,
		Metadata:    meta,
		Parallelism: parallelism,
	}

	results, err := c.protoClient.GetBulkState(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "error getting state")
	}

	if results == nil || results.Items == nil {
		return items, nil
	}

	for _, r := range results.Items {
		item := &BulkStateItem{
			Key:      r.Key,
			Etag:     r.Etag,
			Value:    r.Data,
			Metadata: r.Metadata,
			Error:    r.Error,
		}
		items = append(items, item)
	}

	return items, nil
}

// GetState retrieves state from specific store using default consistency option.
func (c *GRPCClient) GetState(ctx context.Context, storeName, key string) (item *StateItem, err error) {
	return c.GetStateWithConsistency(ctx, storeName, key, nil, StateConsistencyStrong)
}

// GetStateWithConsistency retrieves state from specific store using provided state consistency.
func (c *GRPCClient) GetStateWithConsistency(ctx context.Context, storeName, key string, meta map[string]string, sc StateConsistency) (item *StateItem, err error) {
	if err := hasRequiredStateArgs(storeName, key); err != nil {
		return nil, errors.Wrap(err, "missing required arguments")
	}

	req := &runtimev1pb.GetStateRequest{
		StoreName:   storeName,
		Key:         key,
		Consistency: runtimev1pb.StateOptions_StateConsistency(sc),
		Metadata:    meta,
	}

	result, err := c.protoClient.GetState(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "error getting state")
	}

	return &StateItem{
		Etag:     result.Etag,
		Key:      key,
		Value:    result.Data,
		Metadata: result.Metadata,
	}, nil
}

// DeleteState deletes content from store using default state options.
func (c *GRPCClient) DeleteState(ctx context.Context, storeName, key string) error {
	return c.DeleteStateWithETag(ctx, storeName, key, nil, nil, nil)
}

// DeleteStateWithETag deletes content from store using provided state options and etag.
func (c *GRPCClient) DeleteStateWithETag(ctx context.Context, storeName, key string, etag *ETag, meta map[string]string, opts *StateOptions) error {
	if err := hasRequiredStateArgs(storeName, key); err != nil {
		return errors.Wrap(err, "missing required arguments")
	}

	req := &runtimev1pb.DeleteStateRequest{
		StoreName: storeName,
		Key:       key,
		Options:   toProtoStateOptions(opts),
		Metadata:  meta,
	}

	if etag != nil {
		req.Etag = &runtimev1pb.Etag{
			Value: etag.Value,
		}
	}

	_, err := c.protoClient.DeleteState(ctx, req)
	if err != nil {
		return errors.Wrap(err, "error deleting state")
	}

	return nil
}

// DeleteBulkState deletes content for multiple keys from store.
func (c *GRPCClient) DeleteBulkState(ctx context.Context, storeName string, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	items := make([]*DeleteStateItem, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		item := &DeleteStateItem{
			Key: keys[i],
		}
		items = append(items, item)
	}

	return c.DeleteBulkStateItems(ctx, storeName, items)
}

// DeleteBulkState deletes content for multiple keys from store.
func (c *GRPCClient) DeleteBulkStateItems(ctx context.Context, storeName string, items []*DeleteStateItem) error {
	if len(items) == 0 {
		return nil
	}

	states := make([]*runtimev1pb.StateItem, 0, len(items))
	for i := 0; i < len(items); i++ {
		item := items[i]
		if err := hasRequiredStateArgs(storeName, item.Key); err != nil {
			return errors.Wrap(err, "missing required arguments")
		}

		state := &runtimev1pb.StateItem{
			Key:      item.Key,
			Metadata: item.Metadata,
			Options:  toProtoStateOptions(item.Options),
		}
		if item.Etag != nil {
			state.Etag = &runtimev1pb.Etag{
				Value: item.Etag.Value,
			}
		}
		states = append(states, state)
	}

	req := &runtimev1pb.DeleteBulkStateRequest{
		StoreName: storeName,
		States:    states,
	}
	_, err := c.protoClient.DeleteBulkState(ctx, req)

	return err
}

func hasRequiredStateArgs(storeName, key string) error {
	if storeName == "" {
		return errors.New("store")
	}
	if key == "" {
		return errors.New("key")
	}
	return nil
}
