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
package utils

import (
	"database/sql"
	"errors"
	"fmt"

	// Package pq is a pure Go Postgres driver for the database/sql package
	_ "github.com/lib/pq"
	"mosn.io/pkg/log"
)

const (
	postgresqlHost      = "host"
	postgresqlPort      = "port"
	postgresqlUsername  = "username"
	postgresqlPassword  = "password"
	postgresqlDBName    = "db"
	postgresqlTableName = "tableName"
	postgresqlBizTag    = "bizTag"
)

type PostgresqlMetaData struct {
	Host      string
	Port      string
	Username  string
	Password  string
	DBName    string
	TableName string
	BizTag    string
}

func ParsePostgresqlMetaData(properties map[string]string) (PostgresqlMetaData, error) {
	meta := PostgresqlMetaData{}
	if val, ok := properties[postgresqlHost]; ok && val != "" {
		meta.Host = val
	} else {
		return meta, errors.New("postgresql store error: missing host address or address error")
	}

	if val, ok := properties[postgresqlPort]; ok && val != "" {
		meta.Port = val
	} else {
		return meta, errors.New("postgresql store error: missing port address or port error")
	}

	if val, ok := properties[postgresqlUsername]; ok && val != "" {
		meta.Username = val
	} else {
		return meta, errors.New("postgresql store error: missing username address or username error")
	}

	if val, ok := properties[postgresqlPassword]; ok && val != "" {
		meta.Password = val
	} else {
		return meta, errors.New("postgresql store error: missing password address or password error")
	}

	if val, ok := properties[postgresqlDBName]; ok && val != "" {
		meta.DBName = val
	} else {
		return meta, errors.New("postgresql store error: missing dbName address or dbName error")
	}

	if val, ok := properties[postgresqlTableName]; ok && val != "" {
		meta.TableName = val
	} else {
		return meta, errors.New("postgresql store error: missing dbName address or dbName error")
	}

	if val, ok := properties[postgresqlBizTag]; ok && val != "" {
		meta.BizTag = val
	} else {
		return meta, errors.New("postgresql store error: missing dbName address or dbName error")
	}

	return meta, nil
}

func NewPostgresqlCli(meta PostgresqlMetaData) *sql.DB {
	info := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", meta.Host, meta.Port, meta.Username, meta.Password, meta.DBName)
	postDB, err := sql.Open("postgres", info)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = postDB.Ping()
	if err != nil {
		log.DefaultLogger.Errorf("faild connected")
		return nil
	}
	log.DefaultLogger.Infof("success connected")
	var ID int
	selectTableExist := fmt.Sprintf(`select count(*) from information_schema.tables where table_schema='public' and table_type='BASE TABLE' and table_name='%s';`, meta.TableName)
	err = postDB.QueryRow(selectTableExist).Scan(&ID)
	if err != nil {
		log.DefaultLogger.Errorf(" Failed to query whether the table exist error, err: %v", err)
	}
	if ID == 0 {
		createSql := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS public.%s
												(
    												id bigint NOT NULL,
													value_id bigint NOT NULL,
													biz_tag character(255) COLLATE pg_catalog."default" NOT NULL,
    												create_time bigint,
   				 									update_time bigint,
    												CONSTRAINT layotto_incr_pkey PRIMARY KEY (id)
												)
											WITH (
    											OIDS = FALSE
											)
												TABLESPACE pg_default;`, meta.TableName)
		_, err = postDB.Exec(createSql)
		if err != nil {
			log.DefaultLogger.Errorf("create table failed， err: %v", err)
			return nil
		}
	}
	return postDB
}
