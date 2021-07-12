using System.ComponentModel;

namespace Layotto.State
{
    public enum StateConcurrency
    {
        [Description(StateConstant.UndefinedType)]
        Undefined = 0,
        [Description("first-write")] FirstWrite = 1,
        [Description("last-write")] LastWrite = 2
    }

    public static class StateConcurrencyExtensions
    {
        public static Protocol.StateOptions.Types.StateConcurrency GetPBConcurrency(this StateConcurrency s)
        {
            return (Protocol.StateOptions.Types.StateConcurrency) (int) s;
        }
    }
}