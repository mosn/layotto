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
	"errors"
	"mosn.io/layotto/components/sequencer"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func GetNextIdRequest2ComponentRequest(req *runtimev1pb.GetNextIdRequest) (*sequencer.GetNextIdRequest, error) {
	result := &sequencer.GetNextIdRequest{}
	if req == nil {
		return nil, errors.New("Cannot convert it since request is nil.")
	}

	result.Key = req.Key
	var incrOption = sequencer.WEAK
	if req.Options != nil {
		if req.Options.Increment == runtimev1pb.SequencerOptions_WEAK {
			incrOption = sequencer.WEAK
		} else if req.Options.Increment == runtimev1pb.SequencerOptions_STRONG {
			incrOption = sequencer.STRONG
		} else {
			return nil, errors.New("Options.Increment is illegal.")
		}
	}
	result.Options = sequencer.SequencerOptions{AutoIncrement: incrOption}
	result.Metadata = req.Metadata
	return result, nil
}
