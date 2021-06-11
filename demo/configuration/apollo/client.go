package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"strconv"
	"sync"
	"time"
)

const (
	storeName = "apollo"
	appid     = "testApplication_yang"
	group     = "application"
)

func main() {
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	testSet(cli)
	// Since configuration data in apollo cache is eventual-consistent,we need to sleep a while before querying new data
	time.Sleep(time.Second * 1)
	testGet(cli)

	testDelete(cli)
	//sleep a while before querying deleted data
	time.Sleep(time.Second * 1)
	testGet(cli)
	// besides sdk,u can also call our server with grpc
	testSubscribeWithGrpc()
	cli.Close()
}

func testSubscribeWithGrpc() {
	// 1. connect with grpc
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := runtimev1pb.NewRuntimeClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get client for subscribe
	var wg sync.WaitGroup
	wg.Add(1)
	cli, err := c.SubscribeConfiguration(ctx)
	// client receive changes
	go func() {
		for {
			data, err := cli.Recv()
			if err != nil {
				fmt.Println("subscribe failed", err)
				continue
			}
			fmt.Println("receive subscribe resp", data)

		}
	}()
	// client send subscribe request
	go func() {
		cli.Send(&runtimev1pb.SubscribeConfigurationRequest{StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"haha", "key1"}})
		cli.Send(&runtimev1pb.SubscribeConfigurationRequest{StoreName: storeName, Group: "application", Label: "prod", AppId: appid, Keys: []string{"heihei"}})
	}()

	// loop write in another gorountine
	go func() {
		idx := 0
		for {
			time.Sleep(1 * time.Second)
			idx++
			newItmes := make([]*runtimev1pb.ConfigurationItem, 0, 10)
			newItmes = append(newItmes, &runtimev1pb.ConfigurationItem{Group: "application", Label: "prod", Key: "heihei", Content: "heihei" + strconv.Itoa(idx)})
			fmt.Println("write start")
			c.SaveConfiguration(ctx, &runtimev1pb.SaveConfigurationRequest{StoreName: storeName, AppId: appid, Items: newItmes})
		}
	}()
	wg.Wait()
}

func testSet(cli client.Client) {
	item1 := &client.ConfigurationItem{Group: group, Label: "prod", Key: "key1", Content: "value1", Tags: map[string]string{
		"release": "1.0.0",
		"feature": "print",
	}}
	item2 := &client.ConfigurationItem{Group: group, Label: "prod", Key: "haha", Content: "heihei", Tags: map[string]string{
		"release": "1.0.0",
		"feature": "haha",
	}}
	saveRequest := &client.SaveConfigurationRequest{StoreName: storeName, AppId: appid}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)
	if cli.SaveConfiguration(context.Background(), saveRequest) != nil {
		fmt.Println("save key failed")
		return
	}
	fmt.Println("save key success")
}

func testGet(cli client.Client) {
	getRequest := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "prod", Keys: []string{"key1", "haha"}}
	items, err := cli.GetConfiguration(context.Background(), getRequest)
	//validate
	if err != nil {
		fmt.Printf("get configuration failed %+v \n", err)
	}
	for _, item := range items {
		fmt.Printf("get configuration after save, %+v \n", item)
	}
}

func testDelete(cli client.Client) {
	request := &client.ConfigurationRequestItem{StoreName: storeName, AppId: appid, Group: group, Label: "prod", Keys: []string{"key1", "haha"}}
	if cli.DeleteConfiguration(context.Background(), request) != nil {
		fmt.Println("delete key failed")
		return
	}
	fmt.Printf("delete keys success\n")
}
