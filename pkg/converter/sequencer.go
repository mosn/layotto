package converter

import (
	"errors"
	"mosn.io/layotto/components/sequencer"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func GetNextIdRequest2ComponentRequest(req *runtimev1pb.GetNextIdRequest) (*sequencer.GetNextIdRequest, error) {
	result := &sequencer.GetNextIdRequest{}
	if req == nil {
		return result, nil
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
