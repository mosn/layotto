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

package aliyun

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/email"
)

func TestInit(t *testing.T) {
	a := &AliyunEmail{}
	conf := &email.Config{
		Metadata: map[string]string{
			email.ClientKey:    "aki",
			email.ClientSecret: "aks",
			email.Endpoint:     "endpoint",
		},
	}
	err := a.Init(context.TODO(), conf)
	assert.Nil(t, err)
}

func TestSendEmail(t *testing.T) {
	a := &AliyunEmail{}
	conf := &email.Config{
		Metadata: map[string]string{
			email.ClientKey:    "aki",
			email.ClientSecret: "aks",
			email.Endpoint:     "endpoint",
		},
	}
	err := a.Init(context.TODO(), conf)
	assert.Nil(t, err)

	request := &email.SendEmailRequest{
		Subject: "a_subject",
		Address: &email.EmailAddress{
			From: "email_send_from",
			To:   []string{"email_send_to"},
		},
		Content: &email.Content{Text: "some words"},
	}
	_, err = a.SendEmail(context.TODO(), request)
	assert.Error(t, err)
}

func TestSendEmailWithTemplate(t *testing.T) {
	a := &AliyunEmail{}
	conf := &email.Config{
		Metadata: map[string]string{
			email.ClientKey:    "aki",
			email.ClientSecret: "aks",
			email.Endpoint:     "endpoint",
		},
	}
	err := a.Init(context.TODO(), conf)
	assert.NoError(t, err)

	request := &email.SendEmailWithTemplateRequest{
		Template: &email.EmailTemplate{
			TemplateId: "a_template",
		},
		Address: &email.EmailAddress{
			From: "email_send_from",
			To:   []string{"receivers_name"},
		},
	}
	_, err = a.SendEmailWithTemplate(context.TODO(), request)
	assert.Error(t, err)
}
