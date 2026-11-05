/*
 * Copyright 2026 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestHostingModes(t *testing.T) {
	for _, mode := range []string{"api", "native", "both"} {
		t.Run(mode, func(t *testing.T) {
			r, api, native, err := assemble(mode, "127.0.0.1:0", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := r.Start(ctx); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				c, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				if err := r.Shutdown(c); err != nil {
					t.Error(err)
				}
			})
			if !r.Status(ctx).Ready {
				t.Fatal("runtime not ready")
			}
			if mode == "native" && api != nil {
				t.Fatal("native mode created API")
			}
			if mode == "api" && native != nil {
				t.Fatal("API mode created native service")
			}
			if api != nil {
				conn, err := grpc.DialContext(ctx, api.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()
				resp, err := pb.NewRuntimeClient(conn).SayHello(ctx, &pb.SayHelloRequest{ServiceName: "hello", Name: "runtime"})
				if err != nil || resp.GetHello() != "hello, runtime" {
					t.Fatalf("response=%v err=%v", resp, err)
				}
			}
			if native != nil {
				conn, err := net.DialTimeout("tcp", native.Address(), time.Second)
				if err != nil {
					t.Fatal(err)
				}
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(time.Second))
				if _, err := conn.Write([]byte("hello")); err != nil {
					t.Fatal(err)
				}
				data := make([]byte, 5)
				if _, err := io.ReadFull(conn, data); err != nil || string(data) != "hello" {
					t.Fatalf("data=%q err=%v", data, err)
				}
				_ = conn.Close()
			}
			if err := r.Shutdown(ctx); err != nil {
				t.Fatal(err)
			}
		})
	}
}
