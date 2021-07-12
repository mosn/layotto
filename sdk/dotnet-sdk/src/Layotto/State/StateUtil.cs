using Google.Protobuf;
using Layotto.Protocol;

namespace Layotto.State
{
    public class StateUtil
    {
        public static Protocol.StateItem ToProtoSaveStateItem(SetStateItem stateItem)
        {
            var res = new Protocol.StateItem
            {
                Key = stateItem.Key,
                Metadata = {stateItem.Metadata},
                Value = ByteString.CopyFrom(stateItem.Value.Span),
                Options = ToProtoStateOptions(stateItem.Options)
            };

            if (stateItem.ETag != null) res.Etag = new Etag {Value = stateItem.ETag.Value};

            return res;
        }

        public static Protocol.StateOptions ToProtoStateOptions(StateOptions options)
        {
            if (options == null) return CopyStateOptionDefaultPB();

            return new Protocol.StateOptions
            {
                Concurrency = options.Concurrency.GetPBConcurrency(),
                Consistency = options.Consistency.GetPBConsistency()
            };
        }

        public static Protocol.StateOptions CopyStateOptionDefaultPB()
        {
            return new Protocol.StateOptions
            {
                Concurrency = StateConstant.StateOptionDefault.Concurrency,
                Consistency = StateConstant.StateOptionDefault.Consistency
            };
        }

        public StateOptions CopyStateOptionDefault()
        {
            return new StateOptions
            {
                Concurrency = (StateConcurrency) (int) StateConstant.StateOptionDefault.Concurrency,
                Consistency = (StateConsistency) (int) StateConstant.StateOptionDefault.Consistency
            };
        }
    }
}