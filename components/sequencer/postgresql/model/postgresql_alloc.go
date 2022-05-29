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
	"sync"
	"time"
)

type PostgresqlAlloc struct {
	Key        string     // Which is biz_ Tag is used to distinguish business
	Step       int32      // Recording step
	CurrentPos int32      // Current segment buffer cursor; There are two buffer buffers in total, which are used circularly
	Buffer     []*Segment // Two buffers, one as a pre cache
	UpdateTime time.Time  // It is convenient to record the update time without cleaning for a long time to prevent memory occupation
	mutex      sync.Mutex
	IsPreload  bool                   // Is preloading
	Waiting    map[string][]chan byte // suspend wait, wait when the buffer is loaded
}

// Segment 号段
type Segment struct {
	Cursor uint64 // current distribution location
	Max    uint64 // the maximum
	Min    uint64 // start value is the minimum value
	InitOk bool   // whether initialization succeeded
}

// PostgresqlModel postgresql model
type PostgresqlModel struct {
	ID          uint64 `json:"id" form:"id"`                   // primary key ID
	BizTag      string `json:"biz_tag" form:"biz_tag"`         // differentiated business
	MaxID       uint64 `json:"max_id" form:"max_id"`           // this biz_ Tag the maximum value of the currently assigned ID number segment
	Step        int32  `json:"step" form:"step"`               // each time the ID number segment length is allocated, the default value is 1
	Description string `json:"description" form:"description"` // describe
	UpdateTime  uint64 `json:"update_time" form:"update_time"` // update time
}

func NewPostgresqlAlloc(pModel *PostgresqlModel) *PostgresqlAlloc {
	return &PostgresqlAlloc{
		Key:        pModel.BizTag,
		Step:       pModel.Step,
		CurrentPos: 0, // use only the first cache for the first time
		Buffer:     make([]*Segment, 0),
		UpdateTime: time.Now(),
		Waiting:    make(map[string][]chan byte),
		IsPreload:  false,
	}
}

// NewPostgresqlSegment For example, the maxid in the initialization DB is 1,
// that is, the number segment starts from 1, the step size is 1000,
// the maximum value is 1000, and the range is 1 ~ 1000. The next range is 1001 ~ 2000
func NewPostgresqlSegment(pModel *PostgresqlModel) *Segment {
	return &Segment{
		Cursor: pModel.MaxID - uint64(pModel.Step+1), // previous value of minimum value
		Max:    pModel.MaxID - 1,                     // DB stores 1 by default, so 1 should be subtracted here
		Min:    pModel.MaxID - uint64(pModel.Step),   // start min
		InitOk: true,
	}
}

func (p *PostgresqlAlloc) Lock() {
	p.mutex.Lock()
}

func (p *PostgresqlAlloc) Unlock() {
	p.mutex.Unlock()
}

func (p *PostgresqlAlloc) HasSeq() bool {
	return p.Buffer[p.CurrentPos].InitOk && p.Buffer[p.CurrentPos].Cursor < p.Buffer[p.CurrentPos].Max
}

func (p *PostgresqlAlloc) WakeUp() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, waitChan := range p.Waiting[p.Key] {
		close(waitChan)
	}
	p.Waiting[p.Key] = p.Waiting[p.Key][:0]
}

func (p *PostgresqlAlloc) HasID(id uint64) bool {
	return id != 0
}
