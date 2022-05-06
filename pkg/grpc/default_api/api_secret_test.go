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

package default_api

import (
	"context"
	"testing"

	"github.com/dapr/components-contrib/secretstores"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	moke_secret "mosn.io/layotto/pkg/mock/components/secret"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestGetSecret(t *testing.T) {
	fakeStore := moke_secret.FakeSecretStore{}
	fakeStores := map[string]secretstores.SecretStore{
		"store1": fakeStore,
		"store2": fakeStore,
	}
	testCases := []struct {
		testName         string
		storeName        string
		key              string
		errorExcepted    bool
		expectedResponse string
		expectedError    codes.Code
	}{
		{
			testName:         "Good Key from store",
			storeName:        "store1",
			key:              "good-key",
			errorExcepted:    false,
			expectedResponse: "life is good",
			expectedError:    codes.OK,
		},
		{
			testName:         "error occur with error-key",
			storeName:        "store2",
			key:              "error-key",
			errorExcepted:    true,
			expectedResponse: "null",
			expectedError:    codes.Internal,
		},
	}
	// Setup API server
	fakeAPI := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil,
		nil, fakeStores)

	// Run test server
	port, _ := freeport.GetFreePort()
	server := startTestRuntimeAPIServer(port, fakeAPI)
	defer server.Stop()

	// Create gRPC test client
	clientConn := createTestClient(port)
	defer clientConn.Close()
	// act
	client := runtimev1pb.NewRuntimeClient(clientConn)
	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			req := &runtimev1pb.GetSecretRequest{
				StoreName: tt.storeName,
				Key:       tt.key,
			}
			resp, err := client.GetSecret(context.Background(), req)

			if !tt.errorExcepted {
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, resp.Data[tt.key], tt.expectedResponse, "Expected responses to be same")
			} else {
				assert.Error(t, err, "Expected error")
				assert.Equal(t, tt.expectedError, status.Code(err))
			}
		})
	}
}

func TestGetBulkSecret(t *testing.T) {

	fakeStore := moke_secret.FakeSecretStore{}
	fakeStores := map[string]secretstores.SecretStore{
		"store1": fakeStore,
	}
	expectedResponse := "life is good"

	testCases := []struct {
		testName         string
		storeName        string
		key              string
		errorExcepted    bool
		expectedResponse string
		expectedError    codes.Code
	}{
		{
			testName:         "Good Key from unrestricted store",
			storeName:        "store1",
			key:              "good-key",
			errorExcepted:    false,
			expectedResponse: expectedResponse,
		},
	}
	// Setup API server
	fakeAPI := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil,
		nil, fakeStores)

	// Run test server
	port, _ := freeport.GetFreePort()
	server := startTestRuntimeAPIServer(port, fakeAPI)
	defer server.Stop()

	// Create gRPC test client
	clientConn := createTestClient(port)
	defer clientConn.Close()

	// act
	client := runtimev1pb.NewRuntimeClient(clientConn)

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			req := &runtimev1pb.GetBulkSecretRequest{
				StoreName: tt.storeName,
			}
			resp, err := client.GetBulkSecret(context.Background(), req)

			if !tt.errorExcepted {
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, resp.Data[tt.key].Secrets[tt.key], tt.expectedResponse, "Expected responses to be same")
			} else {
				assert.Error(t, err, "Expected error")
				assert.Equal(t, tt.expectedError, status.Code(err))
			}
		})
	}

}
