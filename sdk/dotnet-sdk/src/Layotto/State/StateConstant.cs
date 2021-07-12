namespace Layotto.State
{
    public static class StateConstant
    {
        public const string UndefinedType = "undefined";

        public static Protocol.StateOptions StateOptionDefault = new Protocol.StateOptions
        {
            Concurrency = Protocol.StateOptions.Types.StateConcurrency.ConcurrencyLastWrite,
            Consistency = Protocol.StateOptions.Types.StateConsistency.ConsistencyStrong
        };
    }
}