package model

//
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

/**
* @Author: azh
* @Date: 2022/5/13 22:16
* @Context:
 */

import (
	"mosn.io/pkg/log"
	"sync"
	"time"
)

const expiredTime = time.Minute * 20

// Global allocator
// key: biz_tag value: SegmentBuffer
type PostgresqlSeq struct {
	cache sync.Map
}

func (p *PostgresqlSeq) GetId(bizTag string) *PostgresqlAlloc {
	if seq, ok := p.cache.Load(bizTag); ok {
		return seq.(*PostgresqlAlloc)
	}
	return nil
}

func (seq *PostgresqlSeq) Add(alloc *PostgresqlAlloc) string {
	seq.cache.Store(alloc.Key, alloc)
	return alloc.Key
}

func (seq *PostgresqlSeq) Update(key string, bean *PostgresqlAlloc) {
	if element, ok := seq.cache.Load(key); ok {
		alloc := element.(*PostgresqlAlloc)
		alloc.Buffer = bean.Buffer
		alloc.UpdateTime = bean.UpdateTime
	}
}

func NewPostgresqlSeq() *PostgresqlSeq {
	seq := &PostgresqlSeq{}
	go seq.clear()
	return seq
}

// Clean up unused memory for more than 20min
func (p *PostgresqlSeq) clear() {
	for {
		now := time.Now()
		// after 20min
		mm, _ := time.ParseDuration("20m")
		next := now.Add(mm)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		log.DefaultLogger.Infof("start clear goroutine")
		p.cache.Range(func(key, value interface{}) bool {
			alloc := value.(*PostgresqlAlloc)
			if next.Sub(alloc.UpdateTime) > expiredTime {
				log.DefaultLogger.Infof("clear biz_tag: %s cache", key)
				p.cache.Delete(key)
			}
			return true
		})
	}
}
