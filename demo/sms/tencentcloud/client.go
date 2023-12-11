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

package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	smsv1 "mosn.io/layotto/spec/proto/extension/v1/sms"
)

var (
	componentName = "sms_demo"
)

func main() {
	// Dial to the gRPC server
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("conn build failed,err: %+v", err)
		panic(err)
	}
	defer conn.Close()

	// Create a new client
	client := smsv1.NewSmsServiceClient(conn)

	// Make a request to send sms
	request := &smsv1.SendSmsWithTemplateRequest{
		ComponentName: componentName,
		PhoneNumbers:  []string{"+8610001000100"},
		Template: &smsv1.Template{
			TemplateId:     "10000",
			TemplateParams: map[string]string{},
		},
		SignName: "sign_name",
		Metadata: map[string]string{"SdkAppId": "app_id"},
	}
	response, err := client.SendSmsWithTemplate(context.Background(), request)
	if err != nil {
		fmt.Printf("send sms failed: %+v", err)
		panic(err)
	}

	// Print the result of the SendSmsWithTemplate response
	fmt.Println(response)
}
