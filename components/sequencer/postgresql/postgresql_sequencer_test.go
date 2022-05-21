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
	"testing"

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
	vals["password"] = "213213"
	vals["db"] = "test_db"
	return vals
}

//
//func Test_getId(t *testing.T) {
//
//	ctx := context.Background()
//	//
//	vals := initMap()
//	//
//	db, err := InitPostgresql_mock(vals)
//	//db, _, err := sqlmock.New()
//	if err != nil {
//		fmt.Println("mock error")
//	}
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//	for i := 1; i < 10; i++ {
//		id, err := db.GetId(ctx, "test")
//
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println("id: ", id)
//	}
//
//}
//
//func InitPostgresql_mock(proterties map[string]string) (*service.PostgresqlService, error) {
//	s := &utils.PostgresqlServer{}
//	err := s.Init(proterties)
//	return NewPostgresqlServer_mock(), err
//}
//
//func NewPostgresqlServer_mock() *service.PostgresqlService {
//	db, _, _ := sqlmock.New()
//	postgresqlDB := dao2.NewPostgresqlDB(db)
//	postgresqlDAO := dao2.NewPostgresqlDAO(postgresqlDB, db)
//	seq := model.NewPostgresqlSeq()
//	postgresqlService := service.NewPostgresqlService(postgresqlDAO, seq)
//	return postgresqlService
//}
//
//func Test_map(t *testing.T) {
//	ctx := context.Background()
//	db, err := utils.InitPostgresql(initMap())
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	kv := make(map[string]int64)
//	kv["test"] = 1
//	kv["azh"] = 5000
//	for k, v := range kv {
//		err := db.InitMaxId(ctx, k, v, 1)
//		if err != nil {
//			panic(err)
//		}
//	}
//}
//
//func Test_Init(t *testing.T) {
//	// init
//	p := &PostgresqlSequencer{}
//	config := &sequencer.Configuration{}
//	config.Properties = initMap()
//	p.Init(*config)
//
//	// get id
//	req := &sequencer.GetNextIdRequest{}
//	req.Key = "test"
//	id, err := p.GetNextId(req)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	fmt.Println("next id : ", id.NextId)
//}

func Test_GetID_mock(t *testing.T) {
	// 因为采用双buffer+号段模式，所以我采用分层架构，如果用mock来mock开始事物，需要在dao层修改，进行mock.ExpectBegin()
	// 但那样会破坏原有的代码结构, 所以如果想验证是否有没有问题，可以docker启动一个postgresql，然后执行上面注释掉的test就可以了
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
