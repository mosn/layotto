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

package tencentcloud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"mosn.io/layotto/components/sms"
	"mosn.io/layotto/components/sms/tencentcloud"
)

type MockSmsClient struct{}

// MockSmsClient.SendSmsWithContext returns error when the TemplateId equals `-1`,
// else returns empty response.
func (*MockSmsClient) SendSmsWithContext(ctx context.Context, request *tcsms.SendSmsRequest) (*tcsms.SendSmsResponse, error) {
	if *request.TemplateId == "-1" {
		return nil, errors.New("need error")
	}
	resp := tcsms.NewSendSmsResponse()
	resp.Response = &tcsms.SendSmsResponseParams{
		SendStatusSet: []*tcsms.SendStatus{
			{
				Code:        common.StringPtr("code"),
				Message:     common.StringPtr("message"),
				PhoneNumber: common.StringPtr("phoneNumber"),
			},
		},
		RequestId: common.StringPtr("requestId"),
	}
	return resp, nil
}

var (
	ctx     = context.Background()
	mockSms = tencentcloud.NewSmsWithClient(&MockSmsClient{})
)

var (
	confSuccess = &sms.Config{
		Metadata: map[string]string{
			sms.ClientKey:    "ck",
			sms.ClientSecret: "cs",
			sms.Region:       "region",
		},
	}
	confNoClientKey = &sms.Config{
		Metadata: map[string]string{
			sms.ClientSecret: "cs",
			sms.Region:       "region",
		},
	}
	confNoClientSecret = &sms.Config{
		Metadata: map[string]string{
			sms.ClientKey: "ck",
			sms.Region:    "region",
		},
	}
	confNoRegion = &sms.Config{
		Metadata: map[string]string{
			sms.ClientKey:    "ck",
			sms.ClientSecret: "cs",
		},
	}
)

var (
	requestSuccess = &sms.SendSmsWithTemplateRequest{
		PhoneNumbers: []string{"+0010001000100"},
		Template:     &sms.Template{TemplateParams: map[string]string{"0": "param"}},
	}
	requestNoPhoneNumbers = &sms.SendSmsWithTemplateRequest{
		Template: &sms.Template{TemplateParams: map[string]string{"0": "param"}},
	}
	requestNoTemplate = &sms.SendSmsWithTemplateRequest{
		PhoneNumbers: []string{"+0010001000100"},
	}
	requestWithWrongTemplate1 = &sms.SendSmsWithTemplateRequest{
		PhoneNumbers: []string{"+0010001000100"},
		Template:     &sms.Template{TemplateParams: map[string]string{"idx": "param"}},
	}
	requestWithWrongTemplate2 = &sms.SendSmsWithTemplateRequest{
		PhoneNumbers: []string{"+0010001000100"},
		Template:     &sms.Template{TemplateParams: map[string]string{"-1": "param"}},
	}
	requestFailed = &sms.SendSmsWithTemplateRequest{
		PhoneNumbers: []string{"+0010001000100"},
		Template: &sms.Template{
			TemplateId:     "-1",
			TemplateParams: map[string]string{"0": "param"},
		},
	}
)

func testInit(t *testing.T) {
	svc := tencentcloud.NewSms()
	err := svc.Init(ctx, confSuccess)
	assert.NoError(t, err)
}

func testInitNoClientKey(t *testing.T) {
	svc := tencentcloud.NewSms()
	err := svc.Init(ctx, confNoClientKey)
	assert.Error(t, err)
}

func testInitNoClientSecret(t *testing.T) {
	svc := tencentcloud.NewSms()
	err := svc.Init(ctx, confNoClientSecret)
	assert.Error(t, err)
}

func testInitNoRegion(t *testing.T) {
	svc := tencentcloud.NewSms()
	err := svc.Init(ctx, confNoRegion)
	assert.Error(t, err)
}

func TestSms_Init(t *testing.T) {
	t.Run("TestInit", testInit)
	t.Run("TestInitNoClientKey", testInitNoClientKey)
	t.Run("TestInitNoClientSecret", testInitNoClientSecret)
	t.Run("TestInitNoRegion", testInitNoRegion)
}

func testSendSms(t *testing.T) {
	_, err := mockSms.SendSmsWithTemplate(ctx, requestSuccess)
	assert.NoError(t, err)
}

func testSendSmsNotInit(t *testing.T) {
	svc := tencentcloud.NewSms()
	_, err := svc.SendSmsWithTemplate(ctx, requestSuccess)
	assert.Error(t, err)
}

func testSendSmsNoPhoneNumbers(t *testing.T) {
	_, err := mockSms.SendSmsWithTemplate(ctx, requestNoPhoneNumbers)
	assert.Error(t, err)
}

func testSendSmsNoTemplate(t *testing.T) {
	_, err := mockSms.SendSmsWithTemplate(ctx, requestNoTemplate)
	assert.Error(t, err)
}

func testSendSmsWithWrongTemplate(t *testing.T) {
	_, err := mockSms.SendSmsWithTemplate(ctx, requestWithWrongTemplate1)
	assert.Error(t, err)
	_, err = mockSms.SendSmsWithTemplate(ctx, requestWithWrongTemplate2)
	assert.Error(t, err)
}

func testSendSmsFailed(t *testing.T) {
	_, err := mockSms.SendSmsWithTemplate(ctx, requestFailed)
	assert.Error(t, err)
}

func TestSms_SendSmsWithTemplate(t *testing.T) {
	t.Run("TestSendSms", testSendSms)
	t.Run("TestSendSmsNotInit", testSendSmsNotInit)
	t.Run("TestSendSmsNoPhoneNumbers", testSendSmsNoPhoneNumbers)
	t.Run("TestSendSmsNoTemplate", testSendSmsNoTemplate)
	t.Run("TestSendSmsWithWrongTemplate", testSendSmsWithWrongTemplate)
	t.Run("TestSendSmsFailed", testSendSmsFailed)
}
