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
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"

	"mosn.io/layotto/components/email"
)

type AliyunEmail struct {
	client *dm20151123.Client
}

func NewAliyunEmail() email.EmailService {
	return &AliyunEmail{}
}

var _ email.EmailService = (*AliyunEmail)(nil)

func (a *AliyunEmail) Init(ctx context.Context, conf *email.Config) error {
	accessKeyID := conf.Metadata[email.ClientKey]
	accessKeySecret := conf.Metadata[email.ClientSecret]
	endpoint := conf.Metadata[email.Endpoint]
	config := &openapi.Config{
		// accessKey ID
		AccessKeyId: tea.String(accessKeyID),
		// accessKey Secret
		AccessKeySecret: tea.String(accessKeySecret),
		// endpoint, ref https://api.aliyun.com/product/Dm
		Endpoint: tea.String(endpoint),
	}

	client, err := dm20151123.NewClient(config)
	if err != nil {
		return err
	}
	a.client = client
	return nil
}

// SendEmail .
func (a *AliyunEmail) SendEmail(ctx context.Context, req *email.SendEmailRequest) (*email.SendEmailResponse, error) {
	if !a.checkSendRequest(req) {
		return nil, email.ErrInvalid
	}

	// AccountName is the email send from
	accountName := req.Address.From
	// ToAddress is target addresses the email send to
	toAddress := strings.Join(req.Address.To, ",")

	sendMailRequest := &dm20151123.SingleSendMailRequest{}
	sendMailRequest.
		SetAccountName(accountName).
		// AddressType = 1: use the email send from
		// ref https://help.aliyun.com/document_detail/29444.html
		SetAddressType(1).
		// ReplyToAddress = false: the email no need to reply
		// ref https://help.aliyun.com/document_detail/29444.html
		SetReplyToAddress(false).
		SetSubject(req.Subject).
		SetToAddress(toAddress).
		SetTextBody(req.Content.Text)

	resp, err := a.client.SingleSendMail(sendMailRequest)
	if err != nil {
		return nil, err
	}
	return &email.SendEmailResponse{
		RequestId: *resp.Body.RequestId,
	}, nil
}

func (a *AliyunEmail) checkSendRequest(r *email.SendEmailRequest) bool {
	// make sure content not empty
	if r.Content == nil || r.Content.Text == "" {
		return false
	}
	if info := r.Address; info == nil || info.From == "" || len(info.To) == 0 {
		return false
	}
	return true
}

// SendEmailWithTemplate .
// template must have been applied in aliyun console, and there need the template name
// receivers must have been filled in aliyun console, and there need the receivers list name
func (a *AliyunEmail) SendEmailWithTemplate(ctx context.Context, req *email.SendEmailWithTemplateRequest) (*email.SendEmailWithTemplateResponse, error) {
	if !a.checkSendWithTemplateRequest(req) {
		return nil, email.ErrInvalid
	}

	// AccountName is the email send from
	accountName := req.Address.From
	// ReceiversName is the name of the recipient list that is created in advance and uploaded with recipients.
	// Only take the element with index zero
	receiversName := req.Address.To[0]

	sendMailWithTemplateRequest := &dm20151123.BatchSendMailRequest{}
	sendMailWithTemplateRequest.
		SetAccountName(accountName).
		// AddressType = 1: use the email send from
		// ref https://help.aliyun.com/document_detail/29444.html
		SetAddressType(1).
		SetReceiversName(receiversName).
		SetTemplateName(req.Template.TemplateId)

	resp, err := a.client.BatchSendMail(sendMailWithTemplateRequest)
	if err != nil {
		return nil, err
	}
	return &email.SendEmailWithTemplateResponse{
		RequestId: *resp.Body.RequestId,
	}, nil
}

func (a *AliyunEmail) checkSendWithTemplateRequest(r *email.SendEmailWithTemplateRequest) bool {
	// make sure template exist
	if r.Template == nil || r.Template.TemplateId == "" {
		return false
	}
	if info := r.Address; info == nil || info.From == "" || len(info.To) == 0 {
		return false
	}
	return true
}
