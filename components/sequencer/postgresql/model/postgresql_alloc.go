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
	Key        string     // 也就是biz_tag用来区分业务
	Step       int32      // 记录步长
	CurrentPos int32      // 当前使用的 segment buffer光标; 总共两个buffer缓存区，循环使用
	Buffer     []*Segment // 双buffer 一个作为预缓存作用
	UpdateTime time.Time  // 记录更新时间 方便长时间不用进行清理，防止占用内存
	mutex      sync.Mutex
	IsPreload  bool                   // 是否正在预加载
	Waiting    map[string][]chan byte // 挂起等待, buffer加载时的等待
}

// Segment 号段
type Segment struct {
	Cursor uint64 // 当前发放位置
	Max    uint64 // 最大值
	Min    uint64 // 开始值即最小值
	InitOk bool   // 是否初始化成功
}

// PostgresqlModel postgresql 模型
type PostgresqlModel struct {
	ID          uint64 `json:"id" form:"id"`                   // 主键id
	BizTag      string `json:"biz_tag" form:"biz_tag"`         // 区分业务
	MaxID       uint64 `json:"max_id" form:"max_id"`           // 该biz_tag目前所被分配的ID号段的最大值
	Step        int32  `json:"step" form:"step"`               // 每次分配ID号段长度，默认为1
	Description string `json:"description" form:"description"` // 描述
	UpdateTime  uint64 `json:"update_time" form:"update_time"` // 更新时间
}

func NewPostgresqlAlloc(pModel *PostgresqlModel) *PostgresqlAlloc {
	return &PostgresqlAlloc{
		Key:        pModel.BizTag,
		Step:       pModel.Step,
		CurrentPos: 0, // 第一次只使用第一块cache
		Buffer:     make([]*Segment, 0),
		UpdateTime: time.Now(),
		Waiting:    make(map[string][]chan byte),
		IsPreload:  false,
	}
}

// NewPostgresqlSegment 以初始化DB中的maxId为1举例子，也就是号段从1开始，步长为1000,最大值就是1000，范围是1～1000。下一段范围就是1001～2000
func NewPostgresqlSegment(pModel *PostgresqlModel) *Segment {
	return &Segment{
		Cursor: pModel.MaxID - uint64(pModel.Step+1), // 最小值的前一个值
		Max:    pModel.MaxID - 1,                     // DB默认存的是1 所以这里要减1
		Min:    pModel.MaxID - uint64(pModel.Step),   // 开始的最小值
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
