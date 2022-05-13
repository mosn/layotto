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

func main() {

}

func Test_readDB(t *testing.T) {
	s := &PostgresqlServer{}
	s.init()
	db := NewPostgresqlClient(s.conf)
	fmt.Println(db)
	fmt.Println(time.Now().Unix())
}

func Test_yaml(t *testing.T) {
	s := &PostgresqlServer{}
	s.init()
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
	s.init()
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
