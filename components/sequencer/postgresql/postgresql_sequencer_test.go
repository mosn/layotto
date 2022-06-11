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

package postgresql

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"mosn.io/layotto/components/pkg/utils"
	"testing"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
)

func initMap() map[string]string {
	vals := make(map[string]string)
	vals["host"] = "127.0.0.1"
	vals["port"] = "5432"
	vals["username"] = "postgres"
	vals["password"] = "213213"
	vals["db"] = "test_db"
	vals["tableName"] = "layotto_incr"
	vals["bizTag"] = "test11"
	return vals
}

func TestPostgresqlSequencer_Init(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := NewPostgresqlSequencer(log.DefaultLogger)
	p.db = db

	cfg := sequencer.Configuration{
		Properties: initMap(),
		BiggerThan: make(map[string]int64),
	}

	rows := sqlmock.NewRows([]string{"exist"}).AddRow(0)
	mock.ExpectQuery("select exists").WillReturnRows(rows)

	err = p.Init(cfg)
	if err != nil {
		t.Errorf("init postgresql error: %v", err)
	}
	assert.Nil(t, err)
}

func TestPostgresqlSequencer_GetNextId(t *testing.T) {
	p := NewPostgresqlSequencer(log.DefaultLogger)
	meta, err := utils.ParsePostgresqlMetaData(initMap())
	if err != nil {
		t.Errorf("init metadata error: %v", err)
	}


	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()



	mock.ExpectExec("CREATE TABLE IF NOT EXISTS layotto_incr\n(\n    id bigint NOT NULL,\n    value_id bigint NOT NULL,\n    biz_tag character(255) COLLATE pg_catalog.\"default\" NOT NULL,\n    create_time bigint,\n    update_time bigint,\n    CONSTRAINT layotto_incr_pkey PRIMARY KEY (id)\n)\nWITH (\n    OIDS = FALSE\n)\nTABLESPACE pg_default;")
	rows := sqlmock.NewRows([]string{"id", "value_id", "biz_tag", "create_time", "update_time"}).AddRow([]driver.Value{1, 10, meta.BizTag, 111111, 111111}...)
	mock.ExpectBegin()
	mock.ExpectQuery("select").WillReturnRows(rows)
	mock.ExpectExec("update").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()


	p.db = db
	p.metadata = meta

	_, err = p.GetNextId(&sequencer.GetNextIdRequest{Key: meta.BizTag, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: initMap()})

	//if err != nil {
	//	t.Errorf("get id error: %v", err)
	//}

	//assert.Nil(t, err)
}

func TestPostgresqlSequencer_GetSegment(t *testing.T) {
	p := NewPostgresqlSequencer(log.DefaultLogger)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("mock db error: %v", err)
	}
	defer db.Close()

	//rows := sqlmock.NewRows([]string{"value_id", "biz_tag", "create_time", "update_time"}).AddRow(1, "test", 11, 11)
	mock.ExpectBegin()
	//mock.ExpectQuery("select value_id, biz_tag, create_time, update_time from layotto_incr where biz_tag = $1").WithArgs( "test").WillReturnRows(rows)
	//mock.ExpectExec("update layotto_incr set value_id = 3, update_time = ? where biz_tag = ?").WithArgs(1, "test").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	req := &sequencer.GetSegmentRequest{Size: 10, Key: p.metadata.BizTag, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: initMap()}
	p.db = db

	_, _, err = p.GetSegment(req)
	//assert.NoError(t, err)
	//p.Close()
}
//
//func TestLocalNextId(t *testing.T) {
//
//	//db := utils.NewPostgresqlCli(params)
//	p := NewPostgresqlSequencer(log.DefaultLogger)
//	cfg := sequencer.Configuration{
//		Properties: initMap(),
//		BiggerThan: make(map[string]int64),
//	}
//	err := p.Init(cfg)
//	if err != nil {
//		panic(err)
//	}
//	for i := 0; i < 10; i++ {
//		resp, err := p.GetNextId(&sequencer.GetNextIdRequest{Key: p.metadata.BizTag, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: nil})
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println("next id : ", resp.NextId)
//	}
//
//	//ctx, _ := context.WithCancel(context.Background())
//	////
//	//updateParams := fmt.Sprintf(`update %v set value_id = value_id + 1, update_time = $1 where biz_tag = $2`, "layotto_incr")
//	//_, err = db.ExecContext(ctx, updateParams, time.Now().Unix(), "test")
//	//if err != nil {
//	//	fmt.Println("err: ", err)
//	//}
//}

func TestInsert(t *testing.T) {
	//
	//params, err := utils.ParsePostgresqlMetaData(initMap())
	//if err != nil {
	//
	//}
	//cli := utils.NewPostgresqlCli(params)
	//ctx, _ := context.WithCancel(context.Background())

	//insertParams := fmt.Sprintf(`INSERT INTO %s (value_id, biz_tag, create_time, update_time) VALUES (?,?,?,?)`, "layotto_incr")
	//fmt.Println(insertParams)
	//now := time.Now().Unix()
	//_, err = cli.ExecContext(ctx, insertParams, 1, "azh-test", uint64(now), uint64(now))

	//var model PostgresqlModel
	//queryParams := fmt.Sprintf(`select id, value_id, biz_tag, create_time, update_time from %s where biz_tag = $1`, "layotto_incr")
	//err = cli.QueryRow(queryParams, "azh").Scan(&model.ID, &model.ValueId, &model.BizTag, &model.CreateTime, &model.UpdateTime)
	//if err != nil {
	//	fmt.Printf("err: %v\n", err)
	//	//panic(err)
	//}
	//fmt.Println(model.ValueId)

	//sql := fmt.Sprintf(`update %v set value_id = %d, update_time = $1 where biz_tag = $2`, "layotto_incr", 2)
	//_, err = cli.ExecContext(ctx, sql, time.Now().Unix(), "test")
	//if err != nil {
	//	fmt.Printf("err: %v\n", err)
	//	//panic(err)
	//}
}
