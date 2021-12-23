//
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
package default_api

import (
	"mosn.io/layotto/components/lock"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TryLockRequest2ComponentRequest(req *runtimev1pb.TryLockRequest) *lock.TryLockRequest {
	result := &lock.TryLockRequest{}
	if req == nil {
		return result
	}
	result.ResourceId = req.ResourceId
	result.LockOwner = req.LockOwner
	result.Expire = req.Expire
	return result
}

func TryLockResponse2GrpcResponse(compResponse *lock.TryLockResponse) *runtimev1pb.TryLockResponse {
	result := &runtimev1pb.TryLockResponse{}
	if compResponse == nil {
		return result
	}
	result.Success = compResponse.Success
	return result
}

func UnlockGrpc2ComponentRequest(req *runtimev1pb.UnlockRequest) *lock.UnlockRequest {
	result := &lock.UnlockRequest{}
	if req == nil {
		return result
	}
	result.ResourceId = req.ResourceId
	result.LockOwner = req.LockOwner
	return result
}

func UnlockComp2GrpcResponse(compResp *lock.UnlockResponse) *runtimev1pb.UnlockResponse {
	result := &runtimev1pb.UnlockResponse{}
	if compResp == nil {
		return result
	}
	result.Status = runtimev1pb.UnlockResponse_Status(compResp.Status)
	return result
}
