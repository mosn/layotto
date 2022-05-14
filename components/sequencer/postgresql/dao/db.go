package dao

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
* @Date: 2022/5/13 22:13
* @Context:
 */
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	//_ "github.com/bmizerany/pq"
	"mosn.io/layotto/components/sequencer/postgresql/model"
	"time"
)

const postgresqlTableName = "layotto_alloc"

type DB struct {
	db *sql.DB
}

func (db *DB) UpdateMaxID(ctx context.Context, bizTag string, tx *sql.Tx) error {
	querySql := fmt.Sprintf(`update %v set max_id = max_id + step, update_time = $1 where biz_tag = $2`, postgresqlTableName)
	var err error
	var res sql.Result
	now := uint64(time.Now().Unix())
	if tx != nil {
		res, err = tx.ExecContext(ctx, querySql, now, bizTag)
	} else {
		res, err = db.db.ExecContext(ctx, querySql, now, bizTag)
	}
	if err != nil {
		fmt.Printf("update max_id error, bizTag: %s; err: %v\n", bizTag, err)
		fmt.Println("错误： ", err.Error())
		return err
	}
	rowsId, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsId == 0 {
		return errors.New("no update")
	}
	return nil
}

func (db *DB) InitMaxId(ctx context.Context, bizTag string, maxId int64, step int64) error {
	sql := fmt.Sprintf(`update %v set max_id = %d, step = %d, update_time = $1 where biz_tag = $2`, postgresqlTableName, maxId, step)
	res, err := db.db.Exec(sql, time.Now().Unix(), bizTag)
	if err != nil {
		fmt.Printf("init max_id error, bizTag: %s; err: %v\n", bizTag, err)
		return err
	}
	rowsId, err := res.RowsAffected()
	if rowsId == 0 {
		return errors.New("no init max_id")
	}
	return nil
}

func (db *DB) Get(ctx context.Context, bizTag string, tx *sql.Tx) (*model.PostgresqlModel, error) {
	querySql := fmt.Sprintf(`select id, biz_tag, max_id, step, description, update_time from %s where biz_tag = $1`, postgresqlTableName)
	var postgresqlModel model.PostgresqlModel
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, querySql, bizTag).Scan(&postgresqlModel.ID, &postgresqlModel.BizTag, &postgresqlModel.MaxID, &postgresqlModel.Step, &postgresqlModel.Description, &postgresqlModel.UpdateTime)
	} else {
		err = db.db.QueryRowContext(ctx, querySql, bizTag).Scan(&postgresqlModel.ID, &postgresqlModel.BizTag, &postgresqlModel.MaxID, &postgresqlModel.Step, &postgresqlModel.Description, &postgresqlModel.UpdateTime)
	}
	if err != nil {
		fmt.Printf("get postgresqlModel error, biz_tag: %s, err: %v", bizTag, err)
		return nil, err
	}
	return &postgresqlModel, nil
}

func NewPostgresqlDB(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

func (db *DB) Create(ctx context.Context, model *model.PostgresqlModel) error {
	createSql := fmt.Sprintf(`INSERT INTO %s (biz_tag, max_id, step, description, update_time) VALUES (?, ?, ?, ?, ?)`, postgresqlTableName)
	now := time.Now().Unix()
	res, err := db.db.ExecContext(ctx, createSql, model.BizTag, model.MaxID, model.Step, model.Description, uint64(now))
	if err != nil {
		fmt.Printf("insert error; model: %v, err: %v\n", model, err)
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

