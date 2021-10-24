import {
  SequencerOptions as SequencerOptionsPB,
} from '../../proto/runtime_pb';

export const SequencerOptionsAutoIncrement = SequencerOptionsPB.AutoIncrement;

export type SequencerOptions = {
  increment: SequencerOptionsPB.AutoIncrement;
};
