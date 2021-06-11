// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package grpc

import runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"

func stateConsistencyToString(c runtimev1pb.StateOptions_StateConsistency) string {
	switch c {
	case runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case runtimev1pb.StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}

	return ""
}

func stateConcurrencyToString(c runtimev1pb.StateOptions_StateConcurrency) string {
	switch c {
	case runtimev1pb.StateOptions_CONCURRENCY_FIRST_WRITE:
		return "first-write"
	case runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE:
		return "last-write"
	}

	return ""
}
