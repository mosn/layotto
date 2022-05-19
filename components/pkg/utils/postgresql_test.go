package utils

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
* @Date: 2022/5/13 22:24
* @Context:
 */

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	model2 "mosn.io/layotto/components/sequencer/postgresql/model"
)

func initMap() map[string]string {
	vals := make(map[string]string)
	vals["host"] = "127.0.0.1"
	vals["port"] = "5432"
	vals["username"] = "postgres"
	vals["password"] = "213213"
	vals["db"] = "test_db"
	return vals
}
func Test_readDB(t *testing.T) {
	s := &PostgresqlServer{}
	s.init(initMap())
	db := NewPostgresqlClient(s.conf)
	fmt.Println(db)
	fmt.Println(time.Now().Unix())
}

func Test_yaml(t *testing.T) {
	s := &PostgresqlServer{}
	s.init(initMap())
	fmt.Println(s.conf.Postgresql.Host)
}

func Test_redisHost(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {

	}
	fmt.Println(s.Addr())
}

func Test_query(t *testing.T) {
	s := &PostgresqlServer{}
	s.init(initMap())
	postDB := NewPostgresqlClient(s.conf)
	sql, err := postDB.Query("select * from layotto_alloc;")
	if err != nil {
		fmt.Println("query error")
	}
	defer sql.Close()
	for sql.Next() {
		model := model2.PostgresqlModel{}
		err := sql.Scan(&model.ID, &model.BizTag, &model.MaxID, &model.Step, &model.Description, &model.UpdateTime)
		if err != nil {
			fmt.Println("query fialed")
			continue
		}
		fmt.Println(model.MaxID)
	}
}
