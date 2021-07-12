using System.ComponentModel;

namespace Layotto.State
{
    public enum StateConsistency
    {
        [Description(StateConstant.UndefinedType)]
        Undefined = 0,
        [Description("eventual")] Eventual = 1,
        [Description("strong")] Strong = 2
    }

    public static class StateConsistencyExtensions
    {
        public static Protocol.StateOptions.Types.StateConsistency GetPBConsistency(this StateConsistency s)
        {
            return (Protocol.StateOptions.Types.StateConsistency) (int) s;
        }
    }
}