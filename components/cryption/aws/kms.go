// Code generated by github.com/seeflood/protoc-gen-p6 .

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

package aws

import (
	"context"
	"fmt"
	"sync"

	"mosn.io/layotto/components/pkg/actuators"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"

	log "mosn.io/layotto/kit/logger"

	"mosn.io/layotto/components/cryption"
)

const (
	componentName = "kms-aws"
)

var (
	once               sync.Once
	readinessIndicator *actuators.HealthIndicator
	livenessIndicator  *actuators.HealthIndicator
)

func init() {
	readinessIndicator = actuators.NewHealthIndicator()
	livenessIndicator = actuators.NewHealthIndicator()
}

type cy struct {
	client *kms.KMS
	keyID  string
	log    log.Logger
}

func NewCryption() cryption.CryptionService {
	once.Do(func() {
		indicators := &actuators.ComponentsIndicator{ReadinessIndicator: readinessIndicator, LivenessIndicator: livenessIndicator}
		actuators.SetComponentsIndicator(componentName, indicators)
	})
	c := &cy{
		log: log.NewLayottoLogger("cryption/aws"),
	}
	log.RegisterComponentLoggerListener("cryption/aws", c)
	return c
}

func (k *cy) OnLogLevelChanged(outputLevel log.LogLevel) {
	k.log.SetLogLevel(outputLevel)
}

func (k *cy) Init(ctx context.Context, conf *cryption.Config) error {
	accessKey := conf.Metadata[cryption.ClientKey]
	secret := conf.Metadata[cryption.ClientSecret]
	region := conf.Metadata[cryption.Region]
	keyID := conf.Metadata[cryption.KeyID]
	staticCredentials := credentials.NewStaticCredentials(accessKey, secret, "")

	awsConf := &aws.Config{
		Region:      aws.String(region),
		Credentials: staticCredentials,
	}
	client := kms.New(session.New(), awsConf)
	if client == nil {
		readinessIndicator.ReportError("fail to create aws kms client")
		livenessIndicator.ReportError("fail to create aws kms client")
	}
	readinessIndicator.SetStarted()
	livenessIndicator.SetStarted()
	k.client = client
	k.keyID = keyID
	return nil
}

func (k *cy) Decrypt(ctx context.Context, request *cryption.DecryptRequest) (*cryption.DecryptResponse, error) {
	decryptRequest := &kms.DecryptInput{
		CiphertextBlob: request.CipherText,
	}
	decryptResp, err := k.client.Decrypt(decryptRequest)
	if err != nil {
		k.log.Errorf("fail decrypt data, err: %+v", err)
		return nil, fmt.Errorf("fail decrypt data with error: %+v", err)
	}
	resp := &cryption.DecryptResponse{KeyId: *decryptResp.KeyId, PlainText: decryptResp.Plaintext}
	return resp, nil
}

func (k *cy) Encrypt(ctx context.Context, request *cryption.EncryptRequest) (*cryption.EncryptResponse, error) {
	// if keyId specified, use request KeyId
	keyId := k.keyID
	if request.KeyId != "" {
		keyId = request.KeyId
	}
	encryptRequest := &kms.EncryptInput{
		KeyId:     aws.String(keyId),
		Plaintext: request.PlainText,
	}

	encryptResp, err := k.client.Encrypt(encryptRequest)
	if err != nil {
		k.log.Errorf("fail encrypt data, err: %+v", err)
		return nil, fmt.Errorf("fail encrypt data with error: %+v", err)
	}
	resp := &cryption.EncryptResponse{KeyId: *encryptResp.KeyId, CipherText: encryptResp.CiphertextBlob}
	return resp, nil
}
