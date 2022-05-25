package postgresql

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
* @Date: 2022/5/13 22:20
* @Context:
 */

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/sequencer"
	model2 "mosn.io/layotto/components/sequencer/postgresql/model"

	"github.com/DATA-DOG/go-sqlmock"

	_ "github.com/lib/pq"
)

//type name struct {
//	id uint64
//}

//func Test_connected(t *testing.T) {
//
//	//_ := context.Background()
//
//	vals := initMap()
//
//	_, err := utils.InitPostgresql(vals)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
func initMap() map[string]string {
	vals := make(map[string]string)
	vals["host"] = "127.0.0.1"
	vals["port"] = "5432"
	vals["username"] = "postgres"
	vals["password"] = "21321311"
	vals["db"] = "test_db"
	return vals
}

func Test_Init(t *testing.T) {
	config := &sequencer.Configuration{Properties: initMap()}
	p := &PostgresqlSequencer{}
	err := p.Init(*config)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_GetNextId(t *testing.T) {
	p := &PostgresqlSequencer{}
	config := &sequencer.Configuration{Properties: initMap()}
	err := p.Init(*config)
	if err != nil {
		fmt.Println(err)
	}
	if p.client == nil {
		fmt.Println("postgresql client is nil")
		return
	}

	req := &sequencer.GetNextIdRequest{Key: "test"}

	id, err := p.GetNextId(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(id)
}

func Test_Create(t *testing.T) {
	model := &model2.PostgresqlModel{}
	p := &PostgresqlSequencer{}
	err := p.Create(model)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_GetSegment(t *testing.T) {
	req := &sequencer.GetSegmentRequest{Key: "test", Size: 10}
	p := &PostgresqlSequencer{}
	config := &sequencer.Configuration{Properties: initMap()}
	err := p.Init(*config)
	if err != nil {
		fmt.Println(err)
	}
	if p.client == nil {
		fmt.Println("postgresql client is nil")
		return
	}

	_, id, err := p.GetSegment(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("get next id: ", id)
	assert.NoError(t, err)
}

func Test_GetID_mock(t *testing.T) {
	// 因为采用双buffer+号段模式，所以我采用分层架构，如果用mock来mock开始事物，需要在dao层修改，进行mock.ExpectBegin()
	// 但那样会破坏原有的代码结构, 所以如果想验证是否有没有问题，可以docker启动一个postgresql，然后执行上面test就可以了
	// 然后连接后执行postgresql.sql脚本，然后配置文件配置用户密码等信息，就可以验证了
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectExec("update layotto_alloc").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}
