/*
 * Copyright 2021 Layotto Authors
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/anypb"
	"io/ioutil"
	"log"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"net/http"
)

func main() {

	http.HandleFunc("/http", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://server.default.svc.cluster.local:9090/")
		if err != nil {
			log.Printf("http get error: %v", err)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("read body error: %v", err)
			return
		}
		w.Write(body)
	})

	http.HandleFunc("/grpc", func(w http.ResponseWriter, r *http.Request) {
		conn, err := grpc.Dial("localhost:34904", grpc.WithInsecure())
		if err != nil {
			log.Printf("connect layotto error: %v", err)
			return
		}
		cli := runtimev1pb.NewRuntimeClient(conn)
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		body, _ := ioutil.ReadAll(r.Body)

		name := r.Header.Get("name")
		if name != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "name", name)
		}

		resp, err := cli.InvokeService(
			ctx,
			&runtimev1pb.InvokeServiceRequest{
				Id: "server.default.svc.cluster.local:9090",
				Message: &runtimev1pb.CommonInvokeRequest{
					Method:      "/hello",
					ContentType: "",
					Data:        &anypb.Any{Value: body}}},
		)
		if err != nil {
			log.Printf("invoke service error: %v", err)
		}else{
			w.Write(resp.Data.Value)
			w.Write([]byte("\n"))
		}
	})

	http.ListenAndServe(":9080", nil)

}
