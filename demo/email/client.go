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
	"fmt"

	"mosn.io/layotto/spec/proto/extension/v1/email"

	"google.golang.org/grpc"
)

const (
	storeName = "email_demo"
)

func TestSendEmail() {
	// Establish a connection.
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to establish connection: %+v", err)
		return
	}

	// Create a client for the email service.
	c := email.NewEmailServiceClient(conn)

	// Create a request to send a common email to specified addresses.
	req := &email.SendEmailRequest{
		ComponentName: storeName,
		Subject:       "email demo",
		Content:       &email.Content{Text: "hi, this is email demo message"},
		Address:       &email.EmailAddress{From: "email_send_from", To: []string{"email_send_to"}},
	}

	// Get response using the client and the request.
	resp, err := c.SendEmail(context.Background(), req)
	if err != nil {
		fmt.Printf("send email request failed: %+v", err)
	}
	fmt.Printf("send email request success, request id: %s \n", resp.RequestId)
}

func TestSendEmailWithTemplate() {
	// Establish a connection.
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("failed to establish connection: %+v", err)
		return
	}

	// Create a client for the email service.
	c := email.NewEmailServiceClient(conn)

	// Create a request to send an email with template to specified receivers_name.
	req := &email.SendEmailWithTemplateRequest{
		ComponentName: storeName,
		Template:      &email.EmailTemplate{TemplateId: "a_template"},
		Address:       &email.EmailAddress{From: "email_send_from", To: []string{"receivers_name"}},
	}

	// Get response using the client and the request.
	resp, err := c.SendEmailWithTemplate(context.Background(), req)
	if err != nil {
		fmt.Printf("send email with template request failed: %+v", err)
	}
	fmt.Printf("send email with template request success, request id: %s \n", resp.RequestId)
}

func main() {
	TestSendEmail()
	TestSendEmailWithTemplate()
}
