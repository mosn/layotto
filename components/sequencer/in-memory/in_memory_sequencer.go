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

package in_memory

import (
	"sync"

	"go.uber.org/atomic"

	"mosn.io/layotto/components/pkg/actuators"
	"mosn.io/layotto/components/sequencer"
)

var (
	once               sync.Once
	readinessIndicator *actuators.HealthIndicator
	livenessIndicator  *actuators.HealthIndicator
)

const componentName = "in-memory"

func init() {
	readinessIndicator = actuators.NewHealthIndicator()
	livenessIndicator = actuators.NewHealthIndicator()
}

type InMemorySequencer struct {
	data *sync.Map
}

func registerActuator() {
	once.Do(func() {
		indicators := &actuators.ComponentsIndicator{ReadinessIndicator: readinessIndicator, LivenessIndicator: livenessIndicator}
		actuators.SetComponentsIndicator(componentName, indicators)
	})
}

func NewInMemorySequencer() *InMemorySequencer {
	registerActuator()
	return &InMemorySequencer{
		data: &sync.Map{},
	}
}

func (s *InMemorySequencer) Init(_ sequencer.Configuration) error {
	readinessIndicator.SetStarted()
	livenessIndicator.SetStarted()
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
