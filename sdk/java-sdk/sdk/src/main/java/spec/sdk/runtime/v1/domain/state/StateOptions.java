package spec.sdk.runtime.v1.domain.state;

public class StateOptions {
  private final Consistency consistency;
  private final Concurrency concurrency;

  /**
   * Represents options for a state API call.
   * @param consistency The consistency mode.
   * @param concurrency The concurrency mode.
   */
  public StateOptions(Consistency consistency, Concurrency concurrency) {
    this.consistency = consistency;
    this.concurrency = concurrency;
  }

  public Concurrency getConcurrency() {
    return concurrency;
  }

  public Consistency getConsistency() {
    return consistency;
  }

  public enum Consistency {
    EVENTUAL("eventual"),
    STRONG("strong");

    private final String value;

    Consistency(String value) {
      this.value = value;
    }

    public String getValue() {
      return this.value;
    }

    public static Consistency fromValue(String value) {
      return Consistency.valueOf(value);
    }
  }

  public enum Concurrency {
    FIRST_WRITE("first-write"),
    LAST_WRITE("last-write");

    private final String value;

    Concurrency(String value) {
      this.value = value;
    }

    public String getValue() {
      return this.value;
    }

    public static Concurrency fromValue(String value) {
      return Concurrency.valueOf(value);
    }
  }
}
