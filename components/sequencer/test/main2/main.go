package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/components/sequencer/snowflake"
	"mosn.io/pkg/log"
)

const key = "resource_xxx"

var size int = 600000

const (
	defaultMysqlUrl     = "localhost:3306"
	defaultDatabaseName = "layotto_sequencer"
	defaultTableName    = "layotto_sequencer_snowflake"
	defaultUserName     = "root"
	defaultPassword     = "oooooo"
)

func main() {

	s := snowflake.NewSnowFlakeSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	cfg.Properties["mysqlHost"] = defaultMysqlUrl
	cfg.Properties["databaseName"] = defaultDatabaseName
	cfg.Properties["tableName"] = defaultTableName
	cfg.Properties["userName"] = defaultUserName
	cfg.Properties["password"] = defaultPassword

	cfg.Properties["boostPower"] = "3"
	cfg.Properties["paddingFactor"] = "50"
	cfg.Properties["timeBits"] = "28"
	cfg.Properties["workerBits"] = "22"
	cfg.Properties["seqBits"] = "13"
	cfg.Properties["startTime"] = "2022-01-01"

	if err := s.Init(cfg); err != nil {
		fmt.Println(err)
	}

	mset := mapset.NewSet()

	var count int64
	var errCount int
	var temp int64

	var wg sync.WaitGroup
	var mu sync.Mutex

	cores := 2 * runtime.NumCPU()
	wg.Add(cores)
	t1 := time.Now().UnixNano() / 1e6
	for i := 0; i < cores; i++ {
		go func() {
			for j := 0; j < size; j++ {
				mu.Lock()
				resp, err := s.GetNextId(&sequencer.GetNextIdRequest{
					Key: key,
				})
				// if err != nil{
				// 	// fmt.Println(err)
				// 	atomic.AddInt64(&count, 1)
				// }
				// if err == nil {
				// 	mset.Add(resp.NextId)
				// } else {
				// 	fmt.Println(err)
				// 	atomic.AddInt64(&count, 1)
				// }

				if err == nil {
					// fmt.Println(resp.NextId)
					mset.Add(resp.NextId)
					if temp != 0 && resp.NextId <= temp {
						errCount++
					}
					temp = resp.NextId
				} else {
					atomic.AddInt64(&count, 1)
				}
				mu.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	// // assert.Equal(t, size, mset.Cardinality())
	fmt.Println(count)
	fmt.Println(time.Now().UnixNano()/1e6 - t1)
	fmt.Println(mset.Cardinality())
	fmt.Println(errCount)
}
