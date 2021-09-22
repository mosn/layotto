File 设计文档

### API
layotto的file接口主要针对于文件系统实现文件的增删改查的能力。在protobuf文件中增加新的接口定义：

```json
    Put(*PutFileStu) error
    Get(*GetFileStu) (io.ReadCloser, error)
    List(*ListRequest) (*ListResp, error)
    Del(*DelRequest) error
    Complete(int64, bool) error
```

### 核心抽象

文件接口的抽象是按照平常文件的操作来定义的，代码中对于文件的操作往往分为几步：

``` go
    handler := open("file_name",option) //打开文件获取文件句柄
    handler.Put(data[]) //写文件
    handler.Close()//关闭文件句bing
```

### 代码实现

Put和Get对于文件操作来说，都是属于流式操作，List和Del都是属于Unary的操作，本身不会太复杂，着重说一下文件的Put和Get的操作：

#### Get接口
``` go
Get(*GetFileStu) (io.ReadCloser, error)
``` 

Get操作在读取文件流的时候，只需要将数据流返回给Api层即可，在api.go里面会分批的读取数据流里面的数据，然后将数据通过stream返回给应用测：


![img.png](../../../img/file/put.png)

这个地方的stream是一个包含了读写的interface,可以自行实现：

``` go
type ReadCloser interface {
	Reader
	Closer
}
``` 

#### Put接口

``` go
    Put(*PutFileStu) error
```
Put接口中定义了一个自增的streamId用来关联stream和file的关系。这里的设计也可以更改为每次收到stream请求的时候直接new一个对应的file变量出来
然后进行增/删/改/查的能力。这样的设计需要在增加新的file component的时候，在api里面也增加对应的new操作。所以这里将这个操作移动到了file component里面去。


