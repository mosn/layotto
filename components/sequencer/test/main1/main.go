package main

import (
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/components/sequencer/snowflake"
	"mosn.io/pkg/log"
)

const key = "resource_xxx"

var size int = 8000000

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

	// cfg.Properties["boostPower"] = "3"
	// cfg.Properties["paddingFactor"] = "50"
	cfg.Properties["timeBits"] = "28"
	cfg.Properties["workerBits"] = "22"
	cfg.Properties["seqBits"] = "13"
	cfg.Properties["startTime"] = "2022-01-01"

	t1 := time.Now().UnixNano() / 1e6
	if err := s.Init(cfg); err != nil {
		fmt.Println(err)
	}

	fmt.Println("start getting")
	mset := mapset.NewSet()

	var uid int64
	var ccc int64
	var temp int64
	var i int
	var count int
	for i = 0; i < size; i++ {
		resp, err := s.GetNextId(&sequencer.GetNextIdRequest{
			Key: key,
		})
		// mset.Add(resp.NextId)
		// if err != nil{
		// 	count++
		// }
		if err == nil {
			if temp != 0 && resp.NextId <= temp {
				ccc++
			}
			mset.Add(resp.NextId)
			temp = resp.NextId
			count++
		}
		// fmt.Printf("%b\n", resp.NextId)
	}
	// // assert.Equal(t, size, mset.Cardinality())
	fmt.Println(count)
	fmt.Println(uid)
	fmt.Println(time.Now().UnixNano()/1e6 - t1)
	fmt.Println(mset.Cardinality())
	fmt.Println(ccc)
}
