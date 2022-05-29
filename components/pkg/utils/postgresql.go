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
* @Date: 2022/5/13 22:20
* @Context:
 */

import (
	"database/sql"
	"errors"
	"fmt"
	"mosn.io/pkg/log"
	"sync"

	"mosn.io/layotto/components/sequencer/postgresql/config"
	dao2 "mosn.io/layotto/components/sequencer/postgresql/dao"
	"mosn.io/layotto/components/sequencer/postgresql/model"
	"mosn.io/layotto/components/sequencer/postgresql/service"
)

type PostgrepsqlMetaData struct {
	Host     string
	Username string
	Password string
	Port     int
}

const (
	postgresqlHost     = "host"
	postgresqlPort     = "port"
	postgresqlUsername = "username"
	postgresqlPassword = "password"
	postgresqlDBName   = "db"
)

var postDB *sql.DB
var err error

type product struct {
	ProductNo string
	Name      string
	Price     float64
}

type PostgresqlServer struct {
	conf *config.Server
	once sync.Once
}

func (s *PostgresqlServer) NewPostgresqlServer() *service.PostgresqlService {
	db := NewPostgresqlClient(s.conf)
	postgresqlDB := dao2.NewPostgresqlDB(db)
	postgresqlDAO := dao2.NewPostgresqlDAO(postgresqlDB, db)
	seq := model.NewPostgresqlSeq()
	postgresqlService := service.NewPostgresqlService(postgresqlDAO, seq)
	return postgresqlService
}

func InitPostgresql(proterties map[string]string) (*service.PostgresqlService, error) {
	s := &PostgresqlServer{}
	err := s.Init(proterties)
	return s.NewPostgresqlServer(), err
}

func NewPostgresqlClient(conf *config.Server) *sql.DB {
	info := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgresql.Host, conf.Postgresql.Port, conf.Postgresql.Username, conf.Postgresql.Password, conf.Postgresql.Db)
	postDB, err = sql.Open("postgres", info)
	if err != nil {
		return nil
	}
	err = postDB.Ping()
	if err != nil {
		log.DefaultLogger.Infof("faild connected")
		return nil
	}
	log.DefaultLogger.Infof("success connected")
	return postDB
}

func (s *PostgresqlServer) Init(proterties map[string]string) error {
	conf := &config.Server{}

	if val, ok := proterties[postgresqlHost]; ok && val != "" {
		conf.Postgresql.Host = val
	} else {
		return errors.New("posgresql store error: missing host address")
	}

	if val, ok := proterties[postgresqlPort]; ok && val != "0" {
		conf.Postgresql.Port = val
	} else {
		return err
	}

	if val, ok := proterties[postgresqlUsername]; ok {
		conf.Postgresql.Username = val
	} else {
		return errors.New("posgresql store error: username error")
	}

	if val, ok := proterties[postgresqlPassword]; ok && val != "" {
		conf.Postgresql.Password = val
	} else {
		return errors.New("posgresql store error: password error")
	}

	if val, ok := proterties[postgresqlDBName]; ok && val != "" {
		conf.Postgresql.Db = val
	} else {
		return errors.New("posgresql store error: db not found error")
	}
	s.conf = conf
	return nil
}
