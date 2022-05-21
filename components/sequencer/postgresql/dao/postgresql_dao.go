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
* @Date: 2022/5/13 22:15
* @Context:
 */

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"mosn.io/layotto/components/sequencer/postgresql/model"
)

type PostgresqlDao struct {
	sql *sql.DB
	db  *DB
}

func (p *PostgresqlDao) NextSegment(ctx context.Context, bizTag string) (*model.PostgresqlModel, error) {
	// start transition
	tx, err := p.sql.Begin()
	defer func() {
		if err != nil {
			p.rollback(tx)
		}
	}()
	if err = p.checkError(err); err != nil {
		return nil, err
	}

	err = p.db.UpdateMaxID(ctx, bizTag, tx)
	if err = p.checkError(err); err != nil {
		return nil, err
	}

	postgresqlModel, err := p.db.Get(ctx, bizTag, tx)
	if err = p.checkError(err); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err = p.checkError(err); err != nil {
		return nil, err
	}

	return postgresqlModel, nil
}

// rollback
func (p *PostgresqlDao) rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		fmt.Println("rollback error")
	}
}

func (p *PostgresqlDao) checkError(err error) error {
	if err == nil {
		return nil
	}
	//if message, ok := err.(*postgresql.PostgresqlError); ok {
	// fmt.Printf("it's sql error; str:%v", message.Message)
	//}
	return err
}

func NewPostgresqlDAO(db *DB, sql *sql.DB) *PostgresqlDao {
	return &PostgresqlDao{
		db:  db,
		sql: sql,
	}
}

func (p *PostgresqlDao) InitMaxId(ctx context.Context, bizTag string, maxId int64, step int64) error {
	err := p.db.InitMaxId(ctx, bizTag, maxId, step)
	if err = p.checkError(err); err != nil {
		return err
	}
	return nil
}

func (p *PostgresqlDao) Add(ctx context.Context, model *model.PostgresqlModel) error {
	return p.db.Create(ctx, model)
}
