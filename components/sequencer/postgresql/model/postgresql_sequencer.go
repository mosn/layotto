package model

/**
* @Author: azh
* @Date: 2022/5/13 22:16
* @Context:
 */

import (
	"fmt"
	"sync"
	"time"
)

const expiredTime = time.Minute * 10

// 全局分配器
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

// 清理超过10min没用过的内存
func (p *PostgresqlSeq) clear() {
	for {
		now := time.Now()
		// 15分钟后
		mm, _ := time.ParseDuration("10m")
		next := now.Add(mm)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		fmt.Println("start clear goroutine")
		p.cache.Range(func(key, value interface{}) bool {
			alloc := value.(*PostgresqlAlloc)
			if next.Sub(alloc.UpdateTime) > expiredTime {
				fmt.Printf("clear biz_tag: %s cache", key)
				p.cache.Delete(key)
			}
			return true
		})
	}
}
