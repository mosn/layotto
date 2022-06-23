# 使用工具自动测试 Quickstart 文档
Quickstart 是项目的门面, 如果新用户进入仓库后，发现 Quickstart 文档跑不起来，可能会失望的走掉。

所以我们要经常性的测试 Quickstart, 保证能正常运行。

但是……定期手动测试 Quickstart、修复文档中的异常，这个过程实在太花时间了：

<img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*fTI5RbfAK3gAAAAAAAAAAAAAARQnAQ" width="30%" height="30%">

烦死了！

我们用工具自动测试文档吧！

## 原理
用工具按顺序执行 markdown 文档里的所有 shell 脚本, 即, 所有用

~~~markdown
```shell
```
~~~

包裹起来的脚本。

注意，不会执行用 

~~~markdown
```bash
```
~~~

包裹起来的脚本哦。

## step 1. 安装 `mdx`
见 https://github.com/seeflood/mdx#installation

## step 2. 关闭本地可能导致冲突的软件
关闭本地的 Layotto, 避免运行文档时出现端口冲突。

同样的，如果文档中会用 Docker 启动 redis 之类的容器，您需要先关闭、删除可能导致端口冲突、容器名冲突的容器。

## step 3. 运行文档
举个例子，运行 state API 的 Quickstart 文档:

```shell
mdx docs/en/start/state/start.md 
```

## step 4. 报错了？测试驱动开发，优化你的文档吧！

可以把每个文档看成一个 UT,它应该有准备、执行、验证、释放资源 4个阶段。

如果文档运行报错了，说明这个 case 需要优化一下。

这也是"测试驱动开发"的思想，优化文档，让文档具有"可测试性"吧。

比如，我运行 state API 的 Quickstart 文档，发现报错:

```bash
SaveState succeeded.key:key1 , value: hello world 
GetState succeeded.[key:key1 etag:1]: hello world
SaveBulkState succeeded.[key:key1 etag:2]: hello world
SaveBulkState succeeded.[key:key2 etag:2]: hello world
GetBulkState succeeded.key:key2 ,value:hello world ,etag:1 ,metadata:map[] 
GetBulkState succeeded.key:key1 ,value:hello world ,etag:2 ,metadata:map[] 
GetBulkState succeeded.key:key3 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key4 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key5 ,value: ,etag: ,metadata:map[] 
panic: error deleting state: rpc error: code = Aborted desc = failed deleting state with key key1: possible etag mismatch. error from state store: ERR Error running script (call to f_9b5da7354cb61e2ca9faff50f6c43b81c73c0b94): @user_script:1: user_script:1: failed to delete key1

goroutine 1 [running]:
main.testDelete(0x16bc760, 0xc0000ac000, 0x16c56a0, 0xc0000b90e0, 0x15f30e1, 0x5, 0x15f2539, 0x4)
        /Users/qunli/projects/layotto/demo/state/redis/client.go:73 +0x13d
main.main()
        /Users/qunli/projects/layotto/demo/state/redis/client.go:57 +0x2f4
exit status 2
```

经过一顿排查，发现是 demo client 在删除指定 key 时没有传 `etag` 字段，导致 demo 运行出现异常。

看，通过自动测试文档，我们发现了一例 Quickstart bug :)

### Q: 如何编写具有"可测试性"的文档
注: 您可以参考能跑通测试的示例文档： docs/en/start/state/start.md

以下解释一些常见细节：
#### demo 代码应该在不符合预期时 panic
比如，`demo/state/redis/client.go` 这个 demo 里，如果调用 Layotto 时出现 error, 应该直接 panic:

```go
if err := cli.SaveBulkState(ctx, store, item, &item2); err != nil {
	panic(err)
}
```

除了判断 error 外，demo 还应该校验测试结果，如果不符合预期则直接 panic。这就相当于 UT 里，调用了某个方法后，需要对调用结果做校验。

这样的好处是：一但 Quickstart 不符合预期，demo 就会异常退出，让自动化工具能够发现"测试失败了！快找人来修！"

#### 最好在文档结尾删除容器、释放资源
写 UT 时，我们会在最后阶段做释放资源、恢复 Mock 之类的事情；为了让文档具有"可测试性"，也要做类似的事情。

例如在文档最后删除 redis Docker 容器：

```shell
docker rm -f redis-test
```

注: Layotto 的 github workflow 每次执行一个 md 之后，会删除所有容器、关闭 layotto,etcd 等应用。
所以即使文档里不删除容器，也不影响 github workflow 跑测试。
#### 不想让某条命令被执行，怎么办?
`mdx` 默认情况下只会执行 shell 代码块，即这么写的代码块： 

```shell
```shell
```

如果不希望某个代码块被执行，可以把 shell 改成别的，比如:

```bash
```bash
```

#### 某条 shell 命令会 hang 住、影响测试，怎么办?
还是以 docs/en/start/state/start.md 为例。

其中有一段脚本会运行 Layotto, 但是如果运行它就会 hang 住，导致测试工具没法继续运行下一条命令：

```bash
./layotto start -c ../../configs/config_redis.json
```

怎么办呢？

##### 解决方案1:

用 @background 注解，见 https://github.com/seeflood/mdx#background

~~~
```shell @background
./layotto start -c ../../configs/config_standalone.json
```
~~~

##### 解决方案2: 

不运行这段脚本，加一段"以后台方式运行 Layotto" 的"隐藏脚本"，这段隐藏脚本用注释包裹住，所以不会被阅读文档的人看到，但是`mdx` 依然会运行它:

```bash
    ```bash
    ./layotto start -c ../../configs/config_redis.json
    ```
    
    <!-- The command below will be run when testing this file 
    ```shell
    nohup ./layotto start -c ../../configs/config_redis.json &
    ```
    -->
```

#### 切换目录的命令怎么处理?
我们可以假设，当前目录是项目根路径。

那么切换路径可以这样写：

```bash
# change directory to ${your project path}/demo/state/redis/
 cd demo/state/redis/
 go run .
```

如果你运行了这个命令后，想再回到根路径怎么办?

##### 解决方案1.
用 `${project_path}` 变量,代表项目根路径，见 https://github.com/seeflood/mdx#cd-project_path

```shell 
cd ${project_path}/demo/state/redis/
```

##### 解决方案2. 
加一段隐藏脚本，用来切换目录。例如这么写:

    <!-- The command below will be run when testing this file 
    ```shell
    cd ../../
    # if we should wait for layotto start, we can:
    # sleep 1s 
    ```
    -->
    
    ```shell
    # open a new terminal tab
    # change directory to ${your project path}/demo/state/redis/
    cd demo/state/redis/
    go run .
    ```

### 其他 markdown 注解
mdx 工具提供了很多"markdown 注解"，帮助您编写"可以运行的 markdown 文件"。感兴趣可以查看[mdx文档](https://github.com/seeflood/mdx#usage)

### 修复报错，看看效果吧!
经过一顿修复，我再次运行文档:

```shell
mdx docs/en/start/state/start.md
```

文档不报错了，能正常运行并退出:

```bash
admindeMacBook-Pro-2:layotto qunli$ mdx docs/en/start/state/start.md
latest: Pulling from library/redis
Digest: sha256:69a3ab2516b560690e37197b71bc61ba245aafe4525ebdece1d8a0bc5669e3e2
Status: Image is up to date for redis:latest
docker.io/library/redis:latest
REPOSITORY                     TAG         IMAGE ID       CREATED         SIZE
redis                          latest      bba24acba395   3 days ago      113MB
pseudomuto/protoc-gen-doc      latest      35472df9ecbb   6 weeks ago     39.5MB
apache/skywalking-oap-server   8.0.1-es7   887769fd3bf6   21 months ago   191MB
apache/skywalking-ui           8.0.1       42b3b496616e   21 months ago   127MB
5835d4652c057ce7ea109564c3e36351335ec53c3dedb02650f2056d3cc3edd5
appending output to nohup.out
runtime client initializing for: 127.0.0.1:34904
SaveState succeeded.key:key1 , value: hello world 
GetState succeeded.[key:key1 etag:1]: hello world
SaveBulkState succeeded.[key:key1 etag:2]: hello world
SaveBulkState succeeded.[key:key2 etag:2]: hello world
GetBulkState succeeded.key:key1 ,value:hello world ,etag:2 ,metadata:map[] 
GetBulkState succeeded.key:key4 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key3 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key5 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key2 ,value:hello world ,etag:1 ,metadata:map[] 
DeleteState succeeded.key:key1
DeleteState succeeded.key:key2
redis-test
```

## step 5. 修改 CI,自动测试新写的 quickstart 文档
如果您新写了一篇 quickstart 文档, 并且自测能正常运行，下一步可以修改 CI，实现"每次有人提 Pull request 时，工具自动测试这篇 quickstart 文档能跑通"。

修改方法是：

1. 修改脚本 `etc/script/test-quickstart.sh`，把您的文档添加到其中:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*ZPRlRa7a-0QAAAAAAAAAAAAAARQnAQ)

2. 如果需要在文档运行前、运行后自动释放一些资源（比如自动 kill 进程、删除 docker 容器），可以在脚本里添加要释放的资源。举个例子，如果想实现"每次运行完一篇文档后，自动 kill etcd 进程"，可以在脚本中添加:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*0th0Q7yn5MIAAAAAAAAAAAAAARQnAQ)

3. 完成上述改动后，就可以测试新的 CI 了。 
   
在项目根目录下运行 

```shell
make style.quickstart
```

会测试这些文档:

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*I7LRSryXwWYAAAAAAAAAAAAAARQnAQ)

> [!TIP|label: 本地运行需谨慎，该脚本会删除一些 docker 容器]
> 该命令会删除包含图中关键字的 Docker 容器，如果您不希望删除这些容器，还是不要本地运行了：
> ![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*N3CIRb0883kAAAAAAAAAAAAAARQnAQ)


而如果运行:

```shell
make style.quickstart QUICKSTART_VERSION=1.17
```

会测试以下文档(这些文档在 golang 1.17 及以上的版本才能运行成功):

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*X3F9QJSKq3QAAAAAAAAAAAAAARQnAQ)
