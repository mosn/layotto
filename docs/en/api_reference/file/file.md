# File API

## grpc API definition

```
Put(*PutFileStu) error
Get(*GetFileStu) (io.ReadCloser, error)
List(*ListRequest) (*ListResp, error)
Del(*DelRequest) error
CompletePut(int64, bool) error
```

## Research

Refer：

```
https://github.com/mosn/layotto/issues/98
```

## Explanation

### Put

#### Entry type
The put interface is used to upload files. The input types are as follows：

```
type PutFileStu struct {
	Data        []byte //data receive
	FileName    string //fileName want put
	Metadata    map[string]string //extended fields, sdk can transmit any field, and components can be implemented in detail
	StreamId    int64 //during the file upload process, grpc's client and server will establish a transmission stream, and the corresponding file handle can be found in the component through StreamId
	ChunkNumber int //chunk number, starting from 1
}

```
#### return type

error

----

### Get

#### Entry type

get interface used download file：

```
    type GetFileStu struct {
    ObjectName string  //FileName
    Metadata   map[string]string //extended fields， eg.bucketName，endpoint
    }
```
#### return type

The return type is io.ReadCloser, error. io.ReadCloser implements the read and write interfaces and can be implemented by yourself, as long as it supports streaming, such as net.Pipe() type

---

### List

#### Entry type

The List interface is used to query files in a certain directory (bucket). The input types are as follows:

```
     type ListRequest struct {
         DirectoryName string //Directory name
         Metadata map[string]string //Extension field
     }
```
#### Return value type

```
     type ListResp struct {
     FilesName []string //List of all files in the directory
     }
```
---

### Del

#### Entry type

The Del interface is used to delete a file. The input types are as follows:

```
     type DelRequest struct {
         FileName string //File name to delete
         Metadata map[string]string //Extension field
     }
```

#### Return value type

Return error type

---

### CompletePut

The concept of completePut is to close the handle of the opened file after the file transfer is complete
#### Input type

```

Parameter 1, corresponds to the StreamId in the input parameters of the Put interface
Parameter 2, represents whether EOF is received, if EOF is received, it is true, if an error occurs, it is false
  
```

#### Return value type

Return error type

---