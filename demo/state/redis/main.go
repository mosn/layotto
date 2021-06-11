package main

import (
	"context"
	"fmt"
	client "github.com/layotto/layotto/sdk/go-sdk/client"
	"time"
)

func main() {
	ctx := context.Background()
	json := `{ "data": "hello" }`
	data := []byte(json)
	store := "redis"

	// create the cli
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	// save state with the key key1
	fmt.Printf("saving data: %s\n", string(data))
	if err := cli.SaveState(ctx, store, "key1", data); err != nil {
		panic(err)
	}
	fmt.Println("data saved")

	// get state for key key1
	item, err := cli.GetState(ctx, store, "key1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("data retrieved [key:%s etag:%s]: %s\n", item.Key, item.Etag, string(item.Value))

	// save state with options
	item2 := &client.SetStateItem{
		Etag: &client.ETag{
			Value: "2",
		},
		Key: item.Key,
		Metadata: map[string]string{
			"created-on": time.Now().UTC().String(),
		},
		Value: item.Value,
		Options: &client.StateOptions{
			Concurrency: client.StateConcurrencyLastWrite,
			Consistency: client.StateConsistencyStrong,
		},
	}
	if err := cli.SaveBulkState(ctx, store, item2); err != nil {
		panic(err)
	}
	fmt.Println("data item saved")

	// delete state for key key1
	if err := cli.DeleteState(ctx, store, "key1"); err != nil {
		panic(err)
	}
	fmt.Println("data deleted")

}
