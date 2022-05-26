package main

import (
	"context"
	"flag"
	"fmt"
	"sync"

	"github.com/google/uuid"

	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	resourceId = "resource_a"
)

var storeName string

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
}

func main() {
	flag.Parse()
	if storeName == "" {
		panic("storeName is empty.")
	}
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	// 1. Client trylock
	owner1 := uuid.New().String()
	fmt.Println("client1 prepare to tryLock...")
	resp, err := cli.TryLock(ctx, &runtimev1pb.TryLockRequest{
		StoreName:  storeName,
		ResourceId: resourceId,
		LockOwner:  owner1,
		Expire:     1000,
	})
	if err != nil {
		panic(err)
	}
	if !resp.Success {
		panic("TryLock failed")
	}
	fmt.Printf("client1 got lock!ResourceId is %s\n", resourceId)
	var wg sync.WaitGroup
	wg.Add(1)
	//	2. Client2 tryLock fail
	go func() {
		fmt.Println("client2 prepare to tryLock...")
		owner2 := uuid.New().String()
		resp, err := cli.TryLock(ctx, &runtimev1pb.TryLockRequest{
			StoreName:  storeName,
			ResourceId: resourceId,
			LockOwner:  owner2,
			Expire:     10,
		})
		if err != nil {
			panic(err)
		}
		if resp.Success {
			panic("client2 got lock?!")
		}
		fmt.Printf("client2 failed to get lock.ResourceId is %s\n", resourceId)
		wg.Done()
	}()
	wg.Wait()
	// 3. client 1 unlock
	fmt.Println("client1 prepare to unlock...")
	unlockResp, err := cli.Unlock(ctx, &runtimev1pb.UnlockRequest{
		StoreName:  storeName,
		ResourceId: resourceId,
		LockOwner:  owner1,
	})
	if err != nil {
		panic(err)
	}
	if unlockResp.Status != 0 {
		panic("client1 failed to unlock!")
	}
	fmt.Println("client1 succeeded in unlocking")
	// 4. client 2 get lock
	wg.Add(1)
	go func() {
		fmt.Println("client2 prepare to tryLock...")
		owner2 := uuid.New().String()
		resp, err := cli.TryLock(ctx, &runtimev1pb.TryLockRequest{
			StoreName:  storeName,
			ResourceId: resourceId,
			LockOwner:  owner2,
			Expire:     10,
		})
		if err != nil {
			panic(err)
		}
		if !resp.Success {
			panic("client2 failed to get lock?!")
		}
		fmt.Printf("client2 got lock.ResourceId is %s\n", resourceId)
		// 5. client2 unlock
		unlockResp, err := cli.Unlock(ctx, &runtimev1pb.UnlockRequest{
			StoreName:  storeName,
			ResourceId: resourceId,
			LockOwner:  owner2,
		})
		if err != nil {
			panic(err)
		}
		if unlockResp.Status != 0 {
			panic("client2 failed to unlock!")
		}
		fmt.Println("client2 succeeded in unlocking")
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("Demo success!")
}
