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

func Test_connected(t *testing.T) {

}

func Test_getId(t *testing.T) {

	ctx := context.Background()
	db := utils.InitPostgresql()

	for i := 1; i < 10; i++ {
		id, err := db.GetId(ctx, "azh")

		if err != nil {
			panic(err)
		}
		fmt.Println("id: ", id)
	}

}

func Test_map(t *testing.T) {
	ctx := context.Background()
	db := utils.InitPostgresql()
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

func Test_all(t *testing.T) {
	req := &sequencer.GetNextIdRequest{}
	req.Key = "azh"
	p := &PostgresqlSequencer{}
	config := &sequencer.Configuration{}
	p.Init(*config)
	res, err := p.GetNextId(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

}
