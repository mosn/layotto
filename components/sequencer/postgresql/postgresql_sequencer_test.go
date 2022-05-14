package postgresql

/**
* @Author: azh
* @Date: 2022/5/13 22:20
* @Context:
 */

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
	"testing"
)

type name struct {
	id uint64
}
func Test_connected(t *testing.T) {

	//_ := context.Background()

	vals := initMap()

	_, err := utils.InitPostgresql(vals)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initMap() map[string]string {
	vals := make(map[string]string)
	vals["host"] = "127.0.0.1"
	vals["port"] = "5432"
	vals["username"] = "postgres"
	vals["password"] = "213213"
	vals["db"] =  "test_db"
	return vals
}

func Test_getId(t *testing.T) {

	ctx := context.Background()

	vals := initMap()

	db, err := utils.InitPostgresql(vals)
	if err != nil {
		fmt.Println(err.Error())
	}

	for i := 1; i < 10; i++ {
		id, err := db.GetId(ctx, "test")

		if err != nil {
			panic(err)
		}
		fmt.Println("id: ", id)
	}

}

func Test_map(t *testing.T) {
	ctx := context.Background()
	db, err := utils.InitPostgresql(initMap())
	if err != nil {
		fmt.Println(err.Error())
	}
	kv := make(map[string]int64)
	kv["test"] = 1
	kv["azh"] = 5000
	for k, v := range kv {
		err := db.InitMaxId(ctx, k, v, 1)
		if err != nil {
			panic(err)
		}
	}
}

func Test_Init(t *testing.T)  {
	// init
	p := &PostgresqlSequencer{}
	config := &sequencer.Configuration{}
	config.Properties = initMap()
	p.Init(*config)

	// get id
	req := &sequencer.GetNextIdRequest{}
	req.Key = "test"
	id, err := p.GetNextId(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("next id : ", id.NextId)
}

