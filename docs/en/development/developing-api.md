# Development specification when adding API
Thank you for your support in Layotto!

This document describes how to design and implement a new Layotto API. Layotto is written in Go , if you are not familiar with Go, you can check it out [Go Tutorial](https://tour.golang.org/welcome/1) 。

When developing a new API, you can refer to the code and documentation of other existing APIs, which will make the development much easier.

**Q: Why should the specification be formulated?**

A: Currently there is a lack of documentation, which makes it difficult for users to use;

Contributors can't understand code due to the lack of comments, such as : https://github.com/mosn/layotto/issues/112

It's time-consuming to make up comments for previous documents. We hope future development will contain these.

**Q: Following the specification is too much trouble, will it discourage developers who want to contribute code?**

A: **This specification only restricts "what is needed when adding a new Layotto API" (such as a new API to generate distributed self-incrementing id)**, other prs, such as adding a new component or a new SDK, do not need to follow this specification and are free enough.

## Too long; Don't read
Proposal before development.Make sure the proposal is detailed enough.

Write four documents for users
- Quick start
- Usage document
- API common configuration
- Component configuration

No need to write design documents, but you need to write detailed comments for proto API and component API and take the comments as doc.

The pr of the new API requires two people to code review.


## 1、Publish API proposal to the community and have a full discussion
### 1.1. Publish a detailed proposal
#### 1.1.1. Why the proposal should be detailed
If the proposal is not detailed enough, there is nothing worthy of review by others, and no problem can be found;

The purpose of the review is to brainstorm ideas, and everyone will help analyze the shortcomings of the current design, expose the problems as soon as possible, and avoid rework later.

#### 1.1.2. The content of the proposal
The proposal needs to include the following:

- demand analysis
  - Why create this API
  - Define the boundary of the requirements, which features support and which do not support  
- Product research on the market
- grpc/http API design
- Component API design
- Explain your design

An example of an excellent proposal：https://github.com/dapr/dapr/issues/2988

### 1.2. Proposal review
After submitting the proposal, everyone can discuss it in text;

For important or complex API, we can hold community meetings for design review.


## 2、Development
### 2.1. Code specification
TODO

### 2.2. Test Specification
- Unit test
- Client demo, which can be used for demonstration and integration testing

### 2.3. Document specification
Principle: A user guide is required. However, a design doc for developers is optional, since it may be expired or inconsistent after a period of time. Posting links to the proposal issue and having comments in the code will be sufficient.

#### 2.3.1. Quick start
need to have:

- What can this API do
- What does this quickstart do? What effect you want to achieve? It’s better to provide a picture for illustration
- Steps of operation

Correct example：[Dapr pub-sub quickstart](https://github.com/dapr/quickstarts/tree/v1.0.0/pub-sub) 
Before the operation, explain what effect this demo want to achieve with illustration

![img.png](../../img/development/api/img.png)

Counter-example: The document only writes operation steps 1234, and users do not understand what they want to do.

#### 2.3.2. Using document
The documentation path is under "Reference - API reference", for example, see [State API](https://mosn.io/layotto/#/en/api_reference/state/reference)

>The study found that Dapr has a lot of documentation, such as the State API:
>
> https://docs.dapr.io/developing-applications/building-blocks/state-management/  
> https://docs.dapr.io/reference/api/state_api/
>
> https://docs.dapr.io/operations/components/setup-state-store/
>
> https://docs.dapr.io/reference/components-reference/supported-state-stores/
>
> Since we are in the early stages of the project, Layotto documentation can be slightly less

Need to have ：
##### What is this API? What is the problem to solve
##### What scenarios are appropriate for using this API
##### How to use this API
- List of interfaces.For example：

![img_4.png](../../img/development/api/img_4.png)

List out which interfaces are there. On the one hand, the users of the province go to the proto and don’t know which APIs are related. On the other hand, it can avoid the disgust of users due to the lack of interface documentation
- About the interface`s input and output parameters: use proto comments as interface documentation
  Considering that the interface document needs to be written in Both Chinese and English and may be inconsistent with the code after a long time, it is recommended not to write the interface document but to write the proTO comment in sufficient detail as the interface document. Such as:

```protobuf
// GetStateRequest is the message to get key-value states from specific state store.
message GetStateRequest {
  // Required. The name of state store.
  string store_name = 1;

  // Required. The key of the desired state
  string key = 2;

  // (optional) read consistency mode
  StateOptions.StateConsistency consistency = 3;

  // (optional) The metadata which will be sent to state store components.
  map<string, string> metadata = 4;
}
   
// StateOptions configures concurrency and consistency for state operations
message StateOptions {
  // Enum describing the supported concurrency for state.
  // The API server uses Optimized Concurrency Control (OCC) with ETags.
  // When an ETag is associated with an save or delete request, the store shall allow the update only if the attached ETag matches with the latest ETag in the database.
  // But when ETag is missing in the write requests, the state store shall handle the requests in the specified strategy(e.g. a last-write-wins fashion).
  enum StateConcurrency {
    CONCURRENCY_UNSPECIFIED = 0;
    // First write wins
    CONCURRENCY_FIRST_WRITE = 1;
    // Last write wins
    CONCURRENCY_LAST_WRITE = 2;
  }
  // Enum describing the supported consistency for state.
  enum StateConsistency {
    CONSISTENCY_UNSPECIFIED = 0;
    //  The API server assumes data stores are eventually consistent by default.A state store should:
    //
    // - For read requests, the state store can return data from any of the replicas
    // - For write request, the state store  should asynchronously replicate updates to configured quorum after acknowledging the update request.
    CONSISTENCY_EVENTUAL = 1;

    // When a strong consistency hint is attached, a state store should:
    //
    // - For read requests, the state store should return the most up-to-date data consistently across replicas.
    // - For write/delete requests, the state store should synchronisely replicate updated data to configured quorum before completing the write request.
    CONSISTENCY_STRONG = 2;
  }

  StateConcurrency concurrency = 1;
  StateConsistency consistency = 2;
}
```      

Clear comment should to be written ：
- Mandatory parameter or optional parameter；
- Explain what this field means; It is not enough to just explain the literal meaning. It is necessary to explain the usage mechanism behind it. For example, the consistency and concurrency above should explain what kind of guarantee the server can provide after the user passes a certain option.

  (The comments on consistency and concurrency are actually pasted after simplifying the description on the Dapr document, saving the need to write bilingual documents.)
- If the comment is not clear, explain it on the document

##### Why is it designed this way?
Post a document link if there is a design document or a proposal issue link if there is no document

#### 2.3.3. This document describes the general CONFIGURATION of the API
For example: https://mosn.io/layotto/#/zh/component_specs/state/common

- Configuration file structure
- Explain the general configuration of the API, such as keyPrefix

#### 2.3.4. Documents that describe component configuration
For example: https://mosn.io/layotto/#/zh/component_specs/state/redis

- Configuration item description for this component
- If you want to start this component and run the demo, how to start it


### 2.4. Annotation specifications
#### proto comment as doc

See above

#### Component API comment as doc

If you don't write a bilingual design document, Then the annotation of the component API should carry the role of the design document (explain to other developers). You can post the link of proposal Issue

The criterion for judging good writing is "After sending it out, if community enthusiasts want to contribute components, can they not ask questions face-to-face and see the project for themselves to get started developing components"

If you feel that the explanation is not clear, write a design document, or add a proposal issue and write in more detail

#### Other Matters needing attention

Make sure there are no Chinese comments；

Don't write meaningless comments (repeat method names), such as：

```protobuf
//StopSubscribe stop subs
StopSubscribe()
```

## 3、Submit a pull request
### 3.1. Pr that does not conform to the development specification may not be incorporated into the trunk

### 3.2. The number of cr
The code review of the new API requires two people to review, and the photo is automatically checked by a robot and changed to one person for review.

Other people who pull the request are random and not bound.
