using System.ComponentModel;

namespace Layotto.State
{
    public enum StateOperationType
    {
        [Description(StateConstant.UndefinedType)]
        Undefined = 0,
        [Description("upsert")] Upsert = 1,
        [Description("delete")] Delete
    }

    public static class StateOperationTypeExtensions
    {
        private static readonly string[] Names = {StateConstant.UndefinedType, "upsert", "delete"};

        public static string GetString(this StateOperationType s)
        {
            return Names[(int) s];
        }
    }
}