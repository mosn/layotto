# Sequencer API design document

This document discusses the API "Generate Distributed Unique (auto-ad) id"

## Requirements

### Generating global unique id

Q: When does a global unique id be generated?

A: db does not help you automatically generate it.Like：

- db has a split table that does not help you create id automatically, you need a global unique business id
- No db, like request to generate a traceId

### There is an incremental need for this id.More specifically there are various：

- No increment required.This situation can be resolved by UUID, although the disadvantage is lengthy.**This API does not consider this for the time**
- “Growing trends”.Without a certain increase, most of the cases are on the increase.

Q: What scenarios require trend incremental?

1. For b+ tree db (e.g. MYSQL), the incremental main key can make better use of the cache (cache friendly).

2. Sort the most recent data.For example, needs are up to 100 news. Developers do not want to add timestamp fields or indexes. If the id itself is incremental, the latest 100 messages are sorted by id and sorted directly to：

```
Select * from message order by message-id limit 100
```

This is common when using nosql, because nosql has difficulty indexing the timestamp field

- Increase in lone transfer in sharding.e.g.[Tidb的自增id](https://docs.pingcap.com/zh/tidb/stable/auto-increment) ensures the ID increment generated on a single server and does not guarantee a global (multi-servers) one-tone.

- Global monochrome

The desired id is necessarily incremental, and there is no regression.

### There may be a need for custom id schema

For example, the required id is in the format "Uid first 8 and auto-id"

### Possible information security related requirements

If the ID is continuous, the pick-ups of malicious users will be very easy and the specified URL will be downloaded directly in order;

If an order number is more dangerous, competing for a direct knowledge of the user's volume of one day.So in some application scenarios, IDs are unruly and unruly.

## Product research

| **System**                        | **Securing the generated id unique**                                                                                                                                                             | **Trends** | **Strict increment**                                          | **Availability**              | **Information Security** |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------- | ------------------------------------------------------------- | ----------------------------- | ------------------------ |
| Solo Redis                        | yes.[Special configuration needs to be configured to redis, open and write disks for each of the two disk policies](https://redis.io/topics/persistence) to avoid data loss      | yes        | yes.premised on delay restarting without data | has single point failure risk |                          |
| redis main copy from +sentinel    | no.Copy is asynchronous, that is, waiting for sync copy using Wait command or may drop data after fo, see[文档](https://redis.io/topics/reapplication)                             | yes        | Depending on which data will not be lost                      |                               |                          |
| redis cluster                     | Id.                                                                                                                                                                              | yes        | Id.                                           |                               |                          |
| snowflake                         | no.(Clock callbacks, etc, can cause idrepeat; need to rely on external storage; or declare that NTP must be closed or reliable NTP used to prevent callbacks) | yes        | No                                                            |                               | Good                     |
| Leaf snowflake                    | yes                                                                                                                                                                                              | yes        | No                                                            |                               | Good                     |
| Leaf segment                      | yes                                                                                                                                                                                              | yes        | No                                                            |                               |                          |
| Leaf segment only one Leaf server | yes                                                                                                                                                                                              | yes        | yes                                                           | has single point failure risk |                          |
| zookeeper                         | yes                                                                                                                                                                                              | yes        | yes                                                           |                               |                          |
| etcd                              | yes                                                                                                                                                                                              | yes        | yes                                                           |                               |                          |
| mysql Single Library Table        | yes                                                                                                                                                                                              | yes        | yes                                                           | has single point failure risk |                          |

## 3. grpc API design

### Proto definition

```protobuf
// Sequencer API
rpc GetNextIdd (GetNextIdReques) returnns (GetNextIdResponse) {}


message GetNextIdest Lum
  string store_name = 1;
  /key is the identification of a consequence.
  string key = 2;
  
  Sequencer options = 3;
  // The metadata which will be sent to the component.
  map<string, string> metadata = 4;
}

/// Sequencer Options requirements for increased and unproductive guarantees
message Sequencer Options Fact  
  enum AutoAddition
    // WEAK meaning a "best forest" increasing service. t there is no strict guarantee   
    WAK = 0;
    // STRONG means a strict guarantee of global monotonically increasing
    STRONG = 1;
  }
  
/// enum Uniqueness$6
// / WEAK means a "best ffort" unqueness guarantee.
/// But it might duplicate in some corner cases.
// WEAK = 0;
// STRONG means a strict guarantee of global unity
// STRONG = 1;
// }

  AutoIncreasing revenue = 1;
// Uniqueness union=2;
}

message GetNextIdResponse
  int64 next_id=1;
}
```

Explanation：

In fact, a GetNextId interface requires a key to be passed as a namespace (e.g. an order form name, "order_table"), and a sequencer ensures that the id generated is unique and incremental in that naming space.

SequencerOptions. AutoCreation is used to state user demand for incremental increments, whether trend increase (WEAK) or strict global increment (STRONG)

**Q: Do you want to spell the user as needed?**

A: The API and Layotto run regardless of this, are handled by the sdk or the user themselves, or a special component can do the feature.

**Q: Return type string or int64**

If you return string, if a user uses a return int64 implementation to convert the returned string into int64 in the user's code, then migrate to another component, this conversion may be misreported

If you return int64, the component does not help the user to make some customized spelling.

Select int64 for portability.Do things spelling in sdk

**Q: What to do with the int64 spill?**

Do not consider this for now

### Disputes concerning singularity

The API was originally defined as the user sender's SequencerOptions. Uniqueness enumeration, in which the WEAK representative is "trying to ensure global uniqueness, but the very low probability is likely to be repeat", requires business code to be prepared for repeating and retrying if you get an id to write the library; and STRONG's representatives are strictly global and the user code does not take into account the need for repetition, retry.

- Defines the reason for this enumeration value (benefits)

If strict uniqueness is to be assured, the component will be more productive.For example, a one-stop cron policy will not work and may cause data to be lost and a repeat id to be generated when the machine is restarted; for example, writing a snowflake algorithm directly in the sidecar will not work because there may be a clock-back problem that leads to duplicate calls (NTP clock-synchronization, leap seconds, etc.).[Leaf的snowflake实现](https://tech.meituan.com/2017/04/21/mt-leaf.html) relies on zookeper's callback to the clock;

- Define the disadvantage of this enumeration

More understanding costs for users

- Conclusion

There is a controversy. No value is added for this period.The result of the default return must be a global unique (STRONG).

## Component API

```go
package sequencer

type Store interface {
	// Init this component.
	//
	// The number 'BiggerThan' means that the id generated by this component must be bigger than this number.
	//
	// If the component find that currently the storage can't guarantee this,
	// it can do some initialization like inserting a new id equal to or bigger than this 'BiggerThan' into the storage,
	// or just return an error
	Init(config Configuration) error

	GetNextId(*GetNextIdRequest) (*GetNextIdResponse, error)

	// GetSegment returns a range of id.
	// 'support' indicates whether this method is supported by the component.
	// Layotto runtime will cache the result if this method is supported.
	GetSegment(*GetSegmentRequest) (support bool, result *GetSegmentResponse, err error)
}

type GetNextIdRequest struct {
	Key      string
	Options  SequencerOptions
	Metadata map[string]string
}

type SequencerOptions struct {
	AutoIncrement AutoIncrementOption
}

type AutoIncrementOption string

const (
	WEAK   AutoIncrementOption = "weak"
	STRONG AutoIncrementOption = "strong"
)

type GetNextIdResponse struct {
	NextId int64
}

type GetSegmentRequest struct {
	Size     int
	Key      string
	Options  SequencerOptions
	Metadata map[string]string
}

type GetSegmentResponse struct {
	Segment []int64
}

type Configuration struct {
	BiggerThan map[string]int64
	Properties map[string]string
}
```

**Q: What is the BiggerThan's field?**

All IDs that require components to generate are larger than "biggerThan".

This configuration item is designed to facilitate transplantation by users.For example, the system originally used mysql for hairdressing, id has been generated to 1000 and later migrated to PostgreSQL, and BiggerThan needs to be configured for 1000 so that the PostgreSQL components will be set up at initialization, force id above 1000 or report errors when the id is not found to meet the requirements and direct start.

**Q: What BiggerThan is a map?**

Because each key may have its own biggerThan.

For example, the original app1 made a split table, used some kind of distribution service to generate orders and commodity table ids, which reached 1000 and reached 2000.

The app1 then wants to replace the store as a sequence, then he would like to state that the order form id is over 1,000 and the product table id is above 2000.

Another example is[Leaf的设计](https://tech.meituan.com/2017/04/21/mt-leaf.htm) and one max_id per biz_tag (Leaf's max_id is our biggerThan)

![leaf_max_id.png](/img/sequencer/design/leaf_max_id.png)

**Q: Do not cache at runtime layer?**

Component implementation method： if runtime is cached

```go
GetSegment (*GetSegmentRequest) (support bool, result *GetSegmentResponse, err error)
```

You can define interfaces first, components are not implemented first, and there is performance need to implement them later

## References

[设计分布式唯一id生成](https://www.jianshu.com/p/fb9478687e55)

[Architectural Chat ID generation](https://www.w3cschool.cn/architectroad/architectroad-distributed-id.html)

[Leaf - USG point ID generation system] (https://tech.meituan.com/2017/04/21/mt-leaf.html)
