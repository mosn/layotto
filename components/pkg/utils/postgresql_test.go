package utils

/**
* @Author: azh
* @Date: 2022/5/13 22:24
* @Context:
 */

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	model2 "mosn.io/layotto/components/sequencer/postgresql/model"
	"testing"
	"time"
)

func initMap() map[string]string {
	vals := make(map[string]string)
	vals["host"] = "127.0.0.1"
	vals["port"] = "5432"
	vals["username"] = "postgres"
	vals["password"] = "213213"
	vals["db"] =  "test_db"
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
