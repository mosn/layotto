# File API

## API定义

```
Put(*PutFileStu) error
Get(*GetFileStu) (io.ReadCloser, error)
List(*ListRequest) (*ListResp, error)
Del(*DelRequest) error
```

## 调研

请参照：

```
https://github.com/mosn/layotto/issues/98
```

## 接口解释

---

### Put接口

#### 入参类型
put接口用于上传文件。其入参类型如下：

```
type PutFileStu struct {
	DataStream io.Reader  //用于读取grpc stream传输的文件数据
	FileName   string     //文件名字
	Metadata   map[string]string //拓展字段
}

```

#### 返回类型

返回error类型

----

### Get接口

#### 入参类型

get接口用于下载文件。其入参类型如下：

```
    type GetFileStu struct {
    FileName string  //FileName
    Metadata   map[string]string //扩展字段，比如bucketName，endpoint等
    }
```
#### 返回值类型

返回类型为 io.ReadCloser, error， io.ReadCloser封装了read和write接口，用户可以自己实现，只要支持流式传输即可，比如net.Pipe()类型。


---

### List接口

#### 入参类型

List接口用于查询某个目录(bucket)下的文件。其入参类型如下：

```
    type ListRequest struct {
        DirectoryName string //目录名字
        Metadata      map[string]string //扩展字段
    }
```
#### 返回值类型

```
    type ListResp struct {
    FilesName []string //目录下所有文件列表
    }
```
---

### Del

#### 入参类型

Del接口用于删除某个文件。其入参类型如下：

```
    type DelRequest struct {
        FileName string //删除的文件名
        Metadata map[string]string //扩展字段
    }
```

#### 返回值类型

返回error类型

---