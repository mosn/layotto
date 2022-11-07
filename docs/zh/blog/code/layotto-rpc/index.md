# Layotto 源码解析 —— 处理 RPC 请求

>本文主要以 Dubbo Json RPC 为例来分析 Layotto RPC 处理流程。
> 
>作者：[王志龙](https://github.com/rayowang) | 2022年4月21日

- [概述](#概述)
- [源码分析](#源码分析)
  * [0x00 Layotto 初始化 RPC](#_0x00-layotto-初始化-rpc)
  * [0x01 Dubbo-go-sample client 发起请求](#_0x01-dubbo-go-sample-client-发起请求)
  * [0x02 Mosn EventLoop 读协程处理请求数据](#_0x02-mosn-eventloop-读协程处理请求数据)
  * [0x03 Grpc Sever 作为 NetworkFilter 处理请求](#_0x03-grpc-sever-作为-networkfilter-处理请求)
  * [0x04 Layotto 发送 RPC 请求并写入 Local 虚拟连接](#_0x04-layotto-发送-rpc-请求并写入-local-虚拟连接)
  * [0x05 Mosn 读取 Remote 并执行 Filter 和代理转发](#_0x05-mosn-读取-remote-并执行-filter-和代理转发)
  * [0x06 Dubbo-go-sample server 收到请求返回响应](#_0x06-dubbo-go-sample-server-收到请求返回响应)
  * [0x07 Mosn 框架处理响应并写回 Remote 虚拟连接](#_0x07-mosn-框架处理响应并写回-remote-虚拟连接)
  * [0x08 Layotto 接收 RPC 响应并读取 Local 虚拟连接](#_0x08-layotto-接收-rpc-响应并读取-local-虚拟连接)
  * [0x09 Grpc Sever 处理数据帧返回给客户端](#_0x09-grpc-sever-处理数据帧返回给客户端)
  * [0x10 Dubbo-go-sample client 接收响应](#_0x10-dubbo-go-sample-client-接收响应)
- [总结](#总结)

## 概述
Layotto 作为区别于网络代理 Service Mesh 的分布式原语集合且使用标准协议的 Runtime，具有明确和丰富的语义 API，而 RPC API 就是众多 API 中的一种。通过 RPC API 应用程序开发者可以通过与同样使用 Sidecar 架构的应用本地 Layotto 实例进行交互，从而间接的调用不同服务的方法，并可以利用内置能力完成分布式追踪和诊断，流量调控，错误处理，安全链路等操作。并且 Layotto 的 RPC API 基于 Mosn 的 Grpc handler 设计，除了 Http/Grpc，与其它服务通信时还可以利用Mosn的多协议机制，使用 X-Protocol 协议进行安全可靠通信。如下代码所示，RPC API 的接口与 Dapr 一致，通过 Grpc 接口 InvokeService 即可进行 RPC 调用。

```go
type DaprClient interface {
    // Invokes a method on a remote Dapr app.
    InvokeService(ctx context.Context, in *InvokeServiceRequest, opts ...grpc.CallOption) (*v1.InvokeResponse, error)
    ...
}
```

## 源码分析

为了便于理解，这里从外到内，再从内到外，由数据流转映射到源代码，也就是从Client发起请求，穿越一层一层的逻辑，到 Server 收到请求返回响应，再一层层的回到 Client 收到响应，一层层来分析 Layotto 的 RPC 流程，总共拆分成十步。另外因为 Grpc Client 和 Server 握手及交互相关的内容不是本文重点，所以分析的相对简略一些，其它步骤内容相对详细一些，大家也可以根据自己的情况直接从目录跳转到相应步骤。 

备注：本文基于 commit hash：1d2bed68c3b2372c34a12aeed41be125a4fdd15a

### 0x00 Layotto 初始化 RPC

Layotto 启动流程涉及众多本流程，在此只分析下跟 RPC 相关的及下述流程用的初始化，因为 Layotto 是建立在 Mosn 之上，所以从 Main 函数出发，urfave/cli 库会调用 Mosn 的 StageManager 初始化 Mosn, 进而在 Mosn NetworkFilter 中初始化 GrpcServer，具体流程如下。

```go
mosn.io/mosn/pkg/stagemanager.(*StageManager).runInitStage at stage_manager.go
=>
mosn.io/mosn/pkg/mosn.(*Mosn).initServer at mosn.go
=>
mosn.io/mosn/pkg/filter/network/grpc.(*grpcServerFilterFactory).Init at factory.go
=>
mosn.io/mosn/pkg/filter/network/grpc.(*Handler).New at factory.go
// 新建一个带有地址的 Grpc 服务器。同一个地址返回同一个服务器，只能启动一次
func (s *Handler) New(addr string, conf json.RawMessage, options ...grpc.ServerOption) (*registerServerWrapper, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    sw, ok := s.servers[addr]
    if ok {
        return sw, nil
    }
    ln, err := NewListener(addr)
    if err != nil {
        log.DefaultLogger.Errorf("create a listener failed: %v", err)
        return nil, err
    }
    // 调用 NewRuntimeGrpcServer
    srv, err := s.f(conf, options...)
    if err != nil {
        log.DefaultLogger.Errorf("create a registered server failed: %v", err)
        return nil, err
    }
    sw = &registerServerWrapper{
        server: srv,
        ln:     ln,
    }
    s.servers[addr] = sw
    return sw, nil
}
=
main.NewRunvtimeGrpcServer at main.go
=>
mosn.io/layotto/pkg/runtime.(*MosnRuntime).initRuntime at runtime.go
=>
mosn.io/layotto/pkg/runtime.(*MosnRuntime).initRpcs at runtime.go
=>
mosn.io/layotto/components/rpc/invoker/mosn.(*mosnInvoker).Init at mosninvoker.go
func (m *mosnInvoker) Init(conf rpc.RpcConfig) error {
    var config mosnConfig
    if err := json.Unmarshal(conf.Config, &config); err != nil {
        return err
    }

    // 初始化 RPC 调用前的 Filter
    for _, before := range config.Before {
        m.cb.AddBeforeInvoke(before)
    }

    // 初始化 RPC 调用后的 Filter
    for _, after := range config.After {
        m.cb.AddAfterInvoke(after)
    }

    if len(config.Channel) == 0 {
        return errors.New("missing channel config")
    }

    // 初始化与 Mosn 通信使用的通道、协议及对应端口
    channel, err := channel.GetChannel(config.Channel[0])
    if err != nil {
        return err
    }
    m.channel = channel
    return nil
}
...
// 完成一些列初始化后在 grpcServerFilter 中启动 Grpc Server
mosn.io/mosn/pkg/filter/network/grpc.(*grpcServerFilterFactory).Init at factory.go
func (f *grpcServerFilterFactory) Init(param interface{}) error {
    ...
    opts := []grpc.ServerOption{
        grpc.UnaryInterceptor(f.UnaryInterceptorFilter),
        grpc.StreamInterceptor(f.StreamInterceptorFilter),
    }
    // 经过上述初始化，完成 Grpc registerServerWrapper 的初始化
    sw, err := f.handler.New(addr, f.config.GrpcConfig, opts...)
    if err != nil {
        return err
    }
    // 启动 Grpc sever
    sw.Start(f.config.GracefulStopTimeout)
    f.server = sw
    log.DefaultLogger.Debugf("grpc server filter initialized success")
    return nil
}
...
// StageManager 在 runInitStage 之后进入 runStartStage 启动 Mosn
func (stm *StageManager) runStartStage() {
    st := time.Now()
    stm.SetState(Starting)
    for _, f := range stm.startupStages {
        f(stm.app)
    }

    stm.wg.Add(1)
    // 在所有启动阶段完成后启动 Mosn
    stm.app.Start()
    ...
}
```

### 0x01 Dubbo-go-sample client 发起请求

根据 [Dubbo Json Rpc Example](https://mosn.io/layotto/#/zh/start/rpc/dubbo_json_rpc) 例子运行如下命令

```shell
go run demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go -d '{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}'
```

使用 Layotto 对 App 提供的 Grpc API InvokeService 发起 RPC 调用，经过数据填充和连接建立等流程，最终通过 Grpc clientStream 中调用 SendMsg 向 Layotto 发送数据，具体流程如下。

```go

func main() {
    data := flag.String("d", `{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}`, "-d")
    flag.Parse()
    
    conn, err := grpc.Dial("localhost:34904", grpc.WithInsecure())
    if err != nil {
        log.Fatal(err)
    }

    cli := runtimev1pb.NewRuntimeClient(conn)
    ctx, cancel := context.WithCancel(context.TODO())
    defer cancel()
    // 通过 Grpc 接口 InvokeService 进行 RPC 调用
    resp, err := cli.InvokeService(
        ctx,
       // 使用 runtimev1pb.InvokeServiceRequest 发起 Grpc 请求
        &runtimev1pb.InvokeServiceRequest{
           // 要请求的 server 接口 ID
           Id: "org.apache.dubbo.samples.UserProvider",
            Message: &runtimev1pb.CommonInvokeRequest{
               // 要请求的接口对应的方法名
                Method:        "GetUser",
                ContentType:   "",
                Data:          &anypb.Any{Value: []byte(*data)},
                HttpExtension: &runtimev1pb.HTTPExtension{Verb: runtimev1pb.HTTPExtension_POST},
            },
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(resp.Data.GetValue()))
}
=>
mosn.io/layotto/spec/proto/runtime/v1.(*runtimeClient).InvokeService at runtime.pb.go
=>
google.golang.org/grpc.(*ClientConn).Invoke at call.go
=>
google.golang.org/grpc.(*clientStream).SendMsg at stream.go
=>
google.golang.org/grpc.(*csAttempt).sendMsg at stream.go
=>
google.golang.org/grpc/internal/transport.(*http2Client).Write at http2_client.go
```

### 0x02 Mosn EventLoop 读协程处理请求数据

上文说过 Layotto 的内核相当于是 Mosn，所以当网络连接数据到达时，会先到 Mosn 的 L4 网络层进行读写，具体流程如下。

```go
mosn.io/mosn/pkg/network.(*listener).accept at listener.go
=>
mosn.io/mosn/pkg/server.(*activeListener).OnAccept at handler.go
=>
mosn.io/mosn/pkg/server.(*activeRawConn).ContinueFilterChain at handler.go
=>
mosn.io/mosn/pkg/server.(*activeListener).newConnection at handler.go
=>
mosn.io/mosn/pkg/network.(*connection).Start at connection.go
=>
mosn.io/mosn/pkg/network.(*connection).startRWLoop at connection.go
func (c *connection) startRWLoop(lctx context.Context) {
    c.internalLoopStarted = true

    utils.GoWithRecover(func() {
       // 读协程
        c.startReadLoop()
    }, func(r interface{}) {
        c.Close(api.NoFlush, api.LocalClose)
    })

    if c.checkUseWriteLoop() {
        c.useWriteLoop = true
        utils.GoWithRecover(func() {
           // 写协程
            c.startWriteLoop()
        }, func(r interface{}) {
            c.Close(api.NoFlush, api.LocalClose)
        })
    }
}
```

在 startRWLoop 方法中我们可以看到会分别开启两个协程来分别处理该连接上的读写操作，即 startReadLoop 和 startWriteLoop，在 startReadLoop 中经过如下流转，把网络层读到的数据，由 filterManager 过滤器管理器把数据交由过滤器链进行处理，具体流程如下。

```go
mosn.io/mosn/pkg/network.(*connection).doRead at connection.go
=>
mosn.io/mosn/pkg/network.(*connection).onRead at connection.go
=>
mosn.io/mosn/pkg/network.(*filterManager).OnRead at filtermanager.go
=>
mosn.io/mosn/pkg/network.(*filterManager).onContinueReading at filtermanager.go
func (fm *filterManager) onContinueReading(filter *activeReadFilter) {
    var index int
    var uf *activeReadFilter

    if filter != nil {
        index = filter.index + 1
    }

    // filterManager遍历过滤器进行数据处理
    for ; index < len(fm.upstreamFilters); index++ {
        uf = fm.upstreamFilters[index]
        uf.index = index
        // 对没有初始化的过滤器调用其初始化方法 OnNewConnection，本例为func (f *grpcFilter) OnNewConnection() api.FilterStatus（向 Listener 发送 grpc 连接以唤醒 Listener 的 Accept）
        if !uf.initialized {
            uf.initialized = true

            status := uf.filter.OnNewConnection()

            if status == api.Stop {
                return
            }
        }

        buf := fm.conn.GetReadBuffer()

        if buf != nil && buf.Len() > 0 {
           // 通知相应过滤器处理
            status := uf.filter.OnData(buf)

            if status == api.Stop {
                return
            }
        }
    }
}
=>
mosn.io/mosn/pkg/filter/network/grpc.(*grpcFilter).OnData at filter.go
=>
mosn.io/mosn/pkg/filter/network/grpc.(*grpcFilter).dispatch at filter.go
func (f *grpcFilter) dispatch(buf buffer.IoBuffer) {
    if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
        log.DefaultLogger.Debugf("grpc get datas: %d", buf.Len())
    }
    // 发送数据唤醒连接读取
    f.conn.Send(buf)
    if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
        log.DefaultLogger.Debugf("read dispatch finished")
    }
}
```

### 0x03 Grpc Sever 作为 NetworkFilter 处理请求

第一阶段中从原始连接读取数据，会进入 Grpc Serve 处理，Serve 方法通过 net.Listener 监听连接，每次启动一个新的协程来处理新的连接（handleRawConn），建立一个基于Http2 的 Transport 进行传输层的 RPC 调用，具体流程如下。

```go
google.golang.org/grpc.(*Server).handleRawConn at server.go
func (s *Server) handleRawConn(lisAddr string, rawConn net.Conn) {
    // 校验服务状态
    if s.quit.HasFired() {
        rawConn.Close()
        return
    }
    rawConn.SetDeadline(time.Now().Add(s.opts.connectionTimeout))
    conn, authInfo, err := s.useTransportAuthenticator(rawConn)
    if err != nil {
       ...
    }
    // HTTP2 握手，创建 Http2Server 与客户端交换帧的初始化信息，帧和窗口大小等
    st := s.newHTTP2Transport(conn, authInfo)
    if st == nil {
        return
    }

    rawConn.SetDeadline(time.Time{})
    if !s.addConn(lisAddr, st) {
        return
    }
    // 创建一个协程进行流处理
    go func() {
        s.serveStreams(st)
        s.removeConn(lisAddr, st)
    }()
    ...
}
=>
google.golang.org/grpc.(*Server).serveStreams at server.go
=>
google.golang.org/grpc.(*Server).handleStream at server.go
func (s *Server) handleStream(t transport.ServerTransport, stream *transport.Stream, trInfo *traceInfo) {
    // 找到到需要调用的 FullMethod，此例为 spec.proto.runtime.v1.Runtime/InvokeService
    sm := stream.Method()
    if sm != "" && sm[0] == '/' {
        sm = sm[1:]
    }
    ...
    service := sm[:pos]
    method := sm[pos+1:]

    // 从注册的 service 列表中找到对应 serviceInfo 对象
    srv, knownService := s.services[service]
    if knownService {
        // 根据方法名找到单向请求的 md——MethodDesc，此 demo 为 mosn.io/layotto/spec/proto/runtime/v1._Runtime_InvokeService_Handler
        if md, ok := srv.methods[method]; ok {
            s.processUnaryRPC(t, stream, srv, md, trInfo)
            return
        }
        // 流式请求
        if sd, ok := srv.streams[method]; ok {
            s.processStreamingRPC(t, stream, srv, sd, trInfo)
            return
        }
    }
    ...
=>
google.golang.org/grpc.(*Server).processUnaryRPC at server.go
=>
mosn.io/layotto/spec/proto/runtime/v1._Runtime_InvokeService_Handler at runtime.pb.go
=>
google.golang.org/grpc.chainUnaryServerInterceptors at server.go
=>
// 服务端单向调用拦截器，用以调用 Mosn 的 streamfilter
mosn.io/mosn/pkg/filter/network/grpc.(*grpcServerFilterFactory).UnaryInterceptorFilter at factory.go
=>
google.golang.org/grpc.getChainUnaryHandler at server.go
// 递归生成链式UnaryHandler
func getChainUnaryHandler(interceptors []UnaryServerInterceptor, curr int, info *UnaryServerInfo, finalHandler UnaryHandler) UnaryHandler {
    if curr == len(interceptors)-1 {
        return finalHandler
    }

    return func(ctx context.Context, req interface{}) (interface{}, error) {
       // finalHandler就是mosn.io/layotto/spec/proto/runtime/v1._Runtime_InvokeService_Handler
        return interceptors[curr+1](ctx, req, info, getChainUnaryHandler(interceptors, curr+1, info, finalHandler))
    }
}
```

### 0x04 Layotto 发送 RPC 请求并写入 Local 虚拟连接

接上述 0x03 流程，从 Runtime_InvokeService_Handler 起，由 GRPC 默认 API 转换为 Dapr API，进入 Layotto 提供的对接 Mosn 的轻量 RPC 框架，具体流程如下。

```go
mosn.io/layotto/spec/proto/runtime/v1._Runtime_InvokeService_Handler at runtime.pb.go
=>
mosn.io/layotto/pkg/grpc/default_api.(*api).InvokeService at api.go
=>
mosn.io/layotto/pkg/grpc/dapr.(*daprGrpcAPI).InvokeService at dapr_api.go
=>
mosn.io/layotto/components/rpc/invoker/mosn.(*mosnInvoker).Invoke at mosninvoker.go
// 请求 Mosn 底座和返回响应
func (m *mosnInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (resp *rpc.RPCResponse, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("[runtime][rpc]mosn invoker panic: %v", r)
            log.DefaultLogger.Errorf("%v", err)
        }
    }()

    // 1. 如果超时时间为 0，设置默认 3000ms 超时
    if req.Timeout == 0 {
        req.Timeout = 3000
    }
    req.Ctx = ctx
    log.DefaultLogger.Debugf("[runtime][rpc]request %+v", req)
    // 2. 触发请求执行前的自定义逻辑
    req, err = m.cb.BeforeInvoke(req)
    if err != nil {
        log.DefaultLogger.Errorf("[runtime][rpc]before filter error %s", err.Error())
        return nil, err
    }
    // 3. 核心调用，下文会进行详细分析
    resp, err = m.channel.Do(req)
    if err != nil {
        log.DefaultLogger.Errorf("[runtime][rpc]error %s", err.Error())
        return nil, err
    }
    resp.Ctx = req.Ctx
    // 4. 触发请求返回后的自定义逻辑
    resp, err = m.cb.AfterInvoke(resp)
    if err != nil {
        log.DefaultLogger.Errorf("[runtime][rpc]after filter error %s", err.Error())
    }
    return resp, err
}
=>
mosn.io/layotto/components/rpc/invoker/mosn/channel.(*httpChannel).Do at httpchannel.go
func (h *httpChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
    // 1. 使用上一阶段设置的默认超时设置 context 超时
    timeout := time.Duration(req.Timeout) * time.Millisecond
    ctx, cancel := context.WithTimeout(req.Ctx, timeout)
    defer cancel()

    // 2. 创建连接得到，启动 readloop 协程进行 Layotto 和 Mosn 的读写交互（具体见下文分析）
    conn, err := h.pool.Get(ctx)
    if err != nil {
        return nil, err
    }
    
    // 3. 设置数据写入连接的超时时间
    hstate := conn.state.(*hstate)
    deadline, _ := ctx.Deadline()
    if err = conn.SetWriteDeadline(deadline); err != nil {
        hstate.close()
        h.pool.Put(conn, true)
        return nil, common.Error(common.UnavailebleCode, err.Error())
    }
    // 4. 因为初始化时配置的 Layotto 与 Mosn 交互使用的是 Http 协议，所以这里会构造 Http 请求
    httpReq := h.constructReq(req)
    defer fasthttp.ReleaseRequest(httpReq)

    // 借助 fasthttp 请求体写入虚拟连接
    if _, err = httpReq.WriteTo(conn); err != nil {
        hstate.close()
        h.pool.Put(conn, true)
        return nil, common.Error(common.UnavailebleCode, err.Error())
    }

    // 5. 构造 fasthttp.Response 结构体读取和解析 hstate 的返回，并设置读取超时时间
    httpResp := &fasthttp.Response{}
    hstate.reader.SetReadDeadline(deadline)

    // 在 Mosn 数据返回前这里会阻塞，readloop 协程读取 Mosn 返回的数据之后流程见下述 0x08 阶段
    if err = httpResp.Read(bufio.NewReader(hstate.reader)); err != nil {
        hstate.close()
        h.pool.Put(conn, true)
        return nil, common.Error(common.UnavailebleCode, err.Error())
    }
    h.pool.Put(conn, false)
    ...
}
=>
mosn.io/layotto/components/rpc/invoker/mosn/channel.(*connPool).Get at connpool.go
// Get is get wrapConn by context.Context
func (p *connPool) Get(ctx context.Context) (*wrapConn, error) {
    if err := p.waitTurn(ctx); err != nil {
        return nil, err
    }

    p.mu.Lock()
    // 1. 从连接池获取连接
    if ele := p.free.Front(); ele != nil {
        p.free.Remove(ele)
        p.mu.Unlock()
        wc := ele.Value.(*wrapConn)
        if !wc.isClose() {
            return wc, nil
        }
    } else {
        p.mu.Unlock()
    }

    // 2. 创建新的连接
    c, err := p.dialFunc()
    if err != nil {
        p.freeTurn()
        return nil, err
    }
    wc := &wrapConn{Conn: c}
    if p.stateFunc != nil {
        wc.state = p.stateFunc()
    }
    // 3. 启动 readloop 独立协程读取 Mosn 返回的数据
    if p.onDataFunc != nil {
        utils.GoWithRecover(func() {
            p.readloop(wc)
        }, nil)
    }
    return wc, nil
}
=>
```

上面第二步创建新的连接需要注意下，是调用了 init 阶段的 RegistChannel 初始化的协议中的 dialFunc func() (net.Conn, error)，因为配置里与 Mosn 交互用的是 Http 协议，所以这里是 newHttpChanel，目前还支持 Bolt，Dubbo 等，详见如下代码。

```go
mosn.io/layotto/components/rpc/invoker/mosn/channel.newHttpChannel at httpchannel.go
// newHttpChannel is used to create rpc.Channel according to ChannelConfig
func newHttpChannel(config ChannelConfig) (rpc.Channel, error) {
    hc := &httpChannel{}
    // 为减少连接创建开销的连接池，定义在 mosn.io/layotto/components/rpc/invoker/mosn/channel/connpool.go
    hc.pool = newConnPool(
        config.Size,
        // dialFunc
        func() (net.Conn, error) {
            _, _, err := net.SplitHostPort(config.Listener)
            if err == nil {
                return net.Dial("tcp", config.Listener)
            }
           //创建一对虚拟连接(net.Pipe)，Layotto 持有 local，Mosn 持有 remote, Layotto 向 local 写入，Mosn 会收到数据, Mosn 从 remote读取，执行 filter 逻辑并进行代理转发，再将响应写到 remote ,最后 Layotto 从 remote 读取，获得响应
            local, remote := net.Pipe()
            localTcpConn := &fakeTcpConn{c: local}
            remoteTcpConn := &fakeTcpConn{c: remote}
           // acceptFunc 是定义在 mosn.io/layotto/components/rpc/invoker/mosn/channel.go 中的闭包，闭包中监听了 remote 虚拟连接
            if err := acceptFunc(remoteTcpConn, config.Listener); err != nil {
                return nil, err
            }
            // the goroutine model is:
            // request goroutine --->  localTcpConn ---> mosn
            //        ^                                        |
            //        |                                        |
            //        |                                        |
            //         hstate <-- readloop goroutine     <------
            return localTcpConn, nil
        },
        // stateFunc
        func() interface{} {
            // hstate 是 readloop 协程与 request 协程通信的管道，是一对读写 net.Conn，请求协程从 reader net.Conn 中读数据，readloop 协程序往 writer net.Conn 写数据
            s := &hstate{}
            s.reader, s.writer = net.Pipe()
            return s
        },
        hc.onData,
        hc.cleanup,
    )
    return hc, nil
}
```

### 0x05 Mosn 读取 Remote 并执行 Filter 和代理转发

(1) 与 0x02 类似，filtermanager 执行过滤器处理阶段，这里会到 proxy 中进行代理转发，详见如下代码。

```go
...
mosn.io/mosn/pkg/network.(*filterManager).onContinueReading at filtermanager.go
=>
mosn.io/mosn/pkg/proxy.(*proxy).OnData at proxy.go
func (p *proxy) OnData(buf buffer.IoBuffer) api.FilterStatus {
    if p.fallback {
        return api.Continue
    }

    if p.serverStreamConn == nil {
        ...
        p.serverStreamConn = stream.CreateServerStreamConnection(p.context, proto, p.readCallbacks.Connection(), p)
    }
    //把数据分发到对应协议的解码器，在这里因为是 POST /org.apache.dubbo.samples.UserProvider HTTP/1.1，所以是 mosn.io/mosn/pkg/stream/http.(*serverStreamConnection).serve at stream.go
    p.serverStreamConn.Dispatch(buf)

    return api.Stop
}
=>
```

(2) serverStreamConnection.serve 监听并处理请求到 downstream OnReceive，详见如下代码。

```go
mosn.io/mosn/pkg/stream/http.(*serverStream).handleRequest at stream.go
func (s *serverStream) handleRequest(ctx context.Context) {
    if s.request != nil {
        // set non-header info in request-line, like method, uri
        injectCtxVarFromProtocolHeaders(ctx, s.header, s.request.URI())
        hasData := true
        if len(s.request.Body()) == 0 {
            hasData = false
        }

        if hasData {
           //在此进入 downstream OnReceive
            s.receiver.OnReceive(s.ctx, s.header, buffer.NewIoBufferBytes(s.request.Body()), nil)
        } else {
            s.receiver.OnReceive(s.ctx, s.header, nil, nil)
        }
    }
}
=>
mosn.io/mosn/pkg/proxy.(*downStream).OnReceive at downstream.go
func (s *downStream) OnReceive(ctx context.Context, headers types.HeaderMap, data types.IoBuffer, trailers types.HeaderMap) {
    ...
    var task = func() {
        ...

        phase := types.InitPhase
        for i := 0; i < 10; i++ {
            s.cleanNotify()

            phase = s.receive(s.context, id, phase)
            ...
            }
        }
    }

    if s.proxy.serverStreamConn.EnableWorkerPool() {
        if s.proxy.workerpool != nil {
            // use the worker pool for current proxy
            s.proxy.workerpool.Schedule(task)
        } else {
            // use the global shared worker pool
            pool.ScheduleAuto(task)
        }
        return
    }

    task()
    return

}
```

(3) 上述 ScheduleAuto 调度后，经过 downStream 的 reveive 的各个阶段处理，经过 upstreamRequest、http clientStream 等处理，最终从网络层的 connection.Write 发送数据并进入 WaitNotify 阶段阻塞，详见如下代码。

```go
mosn.io/mosn/pkg/sync.(*workerPool).ScheduleAuto at workerpool.go
=>
mosn.io/mosn/pkg/sync.(*workerPool).spawnWorker at workerpool.go
=>
mosn.io/mosn/pkg/proxy.(*downStream).receive at downstream.go
=>
InitPhase=>DownFilter=>MatchRoute=>DownFilterAfterRoute=>ChooseHost=>DownFilterAfterChooseHost=>DownRecvHeader=>DownRecvData
=>
mosn.io/mosn/pkg/proxy.(*downStream).receiveData at downstream.go
=>
mosn.io/mosn/pkg/proxy.(*upstreamRequest).appendData at upstream.go
=>
mosn.io/mosn/pkg/stream/http.(*clientStream).doSend at stream.go
=>
github.com/valyala/fasthttp.(*Request).WriteTo at http.go
=>
mosn.io/mosn/pkg/stream/http.(*streamConnection).Write at stream.go
>
mosn.io/mosn/pkg/network.(*connection).Write at connection.go
=>
mosn.io/mosn/pkg/proxy.(*downStream).receive at downstream.go
func (s *downStream) receive(ctx context.Context, id uint32, phase types.Phase) types.Phase {
    for i := 0; i <= int(types.End-types.InitPhase); i++ {
        s.phase = phase
        
        switch phase {
        ...
        case types.WaitNotify:
            s.printPhaseInfo(phase, id)
            if p, err := s.waitNotify(id); err != nil {
                return p
            }
        
            if log.Proxy.GetLogLevel() >= log.DEBUG {
            	log.Proxy.Debugf(s.context, "[proxy] [downstream] OnReceive send downstream response %+v", s.downstreamRespHeaders)
            }
        ...
} 
=>
func (s *downStream) waitNotify(id uint32) (phase types.Phase, err error) {
    if atomic.LoadUint32(&s.ID) != id {
        return types.End, types.ErrExit
    }

	if log.Proxy.GetLogLevel() >= log.DEBUG {
		log.Proxy.Debugf(s.context, "[proxy] [downstream] waitNotify begin %p, proxyId = %d", s, s.ID)
	}
	select {
	// 阻塞等待
	case <-s.notify:
	}
	return s.processError(id)
}
```

### 0x06 Dubbo-go-sample server 收到请求返回响应

这里就是 dubbo-go-sample server的处理，暂不展开，贴下日志信息，感兴趣的同学可以回去翻看源码。

```
[2022-04-18/21:03:18 github.com/apache/dubbo-go-samples/rpc/jsonrpc/go-server/pkg.(*UserProvider2).GetUser: user_provider2.go: 53] userID:"A003"
[2022-04-18/21:03:18 github.com/apache/dubbo-go-samples/rpc/jsonrpc/go-server/pkg.(*UserProvider2).GetUser: user_provider2.go: 56] rsp:&pkg.User{ID:"113", Name:"Moorse", Age:30, sex:0, Birth:703394193, Sex:"MAN"}
```

### 0x07 Mosn 框架处理响应并写回 Remote 虚拟连接

接上述 0x05 第三阶段，在 reveive 的循环阶段的 UpRecvData 阶段进入处理响应逻辑，经过一系列处理最终 Response 写回 0x04 中的 remote 虚拟连接，具体流程如下。

```go
mosn.io/mosn/pkg/proxy.(*downStream).receive at downstream.go
func (s *downStream) waitNotify(id uint32) (phase types.Phase, err error) {
    if atomic.LoadUint32(&s.ID) != id {
        return types.End, types.ErrExit
    }
    
    if log.Proxy.GetLogLevel() >= log.DEBUG {
        log.Proxy.Debugf(s.context, "[proxy] [downstream] waitNotify begin %p, proxyId = %d", s, s.ID)
    }
    // 返回响应
    select {
    case <-s.notify:
    }
    return s.processError(id)
}
=>
UpFilter
=>
UpRecvHeader
=>
func (s *downStream) receive(ctx context.Context, id uint32, phase types.Phase) types.Phase {
    for i := 0; i <= int(types.End-types.InitPhase); i++ {
        s.phase = phase

        switch phase {
        ...
        case types.UpRecvData:
            if s.downstreamRespDataBuf != nil {
            	s.printPhaseInfo(phase, id)
            	s.upstreamRequest.receiveData(s.downstreamRespTrailers == nil)
                if p, err := s.processError(id); err != nil {
              	   return p
              }
           }
        ...
}
=>
mosn.io/mosn/pkg/proxy.(*upstreamRequest).receiveData at upstream.go
=>
mosn.io/mosn/pkg/proxy.(*downStream).onUpstreamData at downstream.go
=>
mosn.io/mosn/pkg/proxy.(*downStream).appendData at downstream.go
=>
mosn.io/mosn/pkg/stream/http.(*serverStream).AppendData at stream.go
=>
mosn.io/mosn/pkg/stream/http.(*serverStream).endStream at stream.go
=>
mosn.io/mosn/pkg/stream/http.(*serverStream).doSend at stream.go
=>
github.com/valyala/fasthttp.(*Response).WriteTo at http.go
=>
github.com/valyala/fasthttp.writeBufio at http.go
=>
github.com/valyala/fasthttp.(*statsWriter).Write at http.go
=>
mosn.io/mosn/pkg/stream/http.(*streamConnection).Write at stream.go
```

### 0x08 Layotto 接收 RPC 响应并读取 Local 虚拟连接

上述 0x04 启动的 readloop 协程读IO被激活，从连接读取数Mosn 传回的数据，然后交给 hstate 管道中转处理再返回给请求协程，具体流程如下。

```go
mosn.io/layotto/components/rpc/invoker/mosn/channel.(*connPool).readloop at connpool.go
// readloop is loop to read connected then exec onDataFunc
func (p *connPool) readloop(c *wrapConn) {
    var err error

    defer func() {
        c.close()
        if p.cleanupFunc != nil {
            p.cleanupFunc(c, err)
        }
    }()

    c.buf = buffer.NewIoBuffer(defaultBufSize)
    for {
        // 从连接读取数据
        n, readErr := c.buf.ReadOnce(c)
        if readErr != nil {
            err = readErr
            if readErr == io.EOF {
                log.DefaultLogger.Debugf("[runtime][rpc]connpool readloop err: %s", readErr.Error())
            } else {
                log.DefaultLogger.Errorf("[runtime][rpc]connpool readloop err: %s", readErr.Error())
            }
        }

        if n > 0 {
            // 在onDataFunc 委托给 hstate 处理数据
            if onDataErr := p.onDataFunc(c); onDataErr != nil {
                err = onDataErr
                log.DefaultLogger.Errorf("[runtime][rpc]connpool onData err: %s", onDataErr.Error())
            }
        }

        if err != nil {
            break
        }

        if c.buf != nil && c.buf.Len() == 0 && c.buf.Cap() > maxBufSize {
            c.buf.Free()
            c.buf.Alloc(defaultBufSize)
        }
    }
}
=>
mosn.io/layotto/components/rpc/invoker/mosn/channel.(*httpChannel).onData at httpchannel.go
=>
mosn.io/layotto/components/rpc/invoker/mosn/channel.(*hstate).onData at httpchannel.go
=>
net.(*pipe).Write at pipe.go
=>
mosn.io/layotto/components/rpc/invoker/mosn/channel.(*httpChannel).Do at httpchannel.go
func (h *httpChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
    ...
    // 接上述0x04阶段，mosn 数据返回后，从 hstate 读取 readloop 协程从 mosn 返回的数据
    if err = httpResp.Read(bufio.NewReader(hstate.reader)); err != nil {
        hstate.close()
        h.pool.Put(conn, true)
        return nil, common.Error(common.UnavailebleCode, err.Error())
    }
    h.pool.Put(conn, false)

    // 获取 fasthttp 的数据部分，解析状态码，失败返回错误信息和状态码
    body := httpResp.Body()
    if httpResp.StatusCode() != http.StatusOK {
        return nil, common.Errorf(common.UnavailebleCode, "http response code %d, body: %s", httpResp.StatusCode(), string(body))
    }
    
    // 6. 将结果转换为 rpc.RPCResponse 返回
    rpcResp := &rpc.RPCResponse{
        ContentType: string(httpResp.Header.ContentType()),
        Data:        body,
        Header:      map[string][]string{},
    }
    httpResp.Header.VisitAll(func(key, value []byte) {
        rpcResp.Header[string(key)] = []string{string(value)}
    })
    return rpcResp, nil
```

### 0x09 Grpc Sever 处理数据帧返回给客户端

Grpc 并没有直接写入数据到连接，而是用协程异步 loop 循环从一个缓存结构里面获取帧然后写回到客户端，具体流程如下。

```go
google.golang.org/grpc/internal/transport.NewServerTransport at http2_server.go
func NewServerTransport(conn net.Conn, config *ServerConfig) (_ ServerTransport, err error) {
    ...
    // 协程异步loop循环
    go func() {
        t.loopy = newLoopyWriter(serverSide, t.framer, t.controlBuf, t.bdpEst)
        t.loopy.ssGoAwayHandler = t.outgoingGoAwayHandler
        if err := t.loopy.run(); err != nil {
            if logger.V(logLevel) {
                logger.Errorf("transport: loopyWriter.run returning. Err: %v", err)
            }
        }
        t.conn.Close()
        t.controlBuf.finish()
        close(t.writerDone)
    }()
    go t.keepalive()
    return t, nil
}
=>
google.golang.org/grpc/internal/transport.(*loopyWriter).run at controlbuf.go
=>
google.golang.org/grpc/internal/transport.(*bufWriter).Flush at http_util.go
=>
mosn.io/mosn/pkg/filter/network/grpc.(*Connection).Write at conn.go
=>
mosn.io/mosn/pkg/network.(*connection).Write at connection.go
=>
mosn.io/mosn/pkg/network.(*connection).writeDirectly at connection.go
=>
mosn.io/mosn/pkg/network.(*connection).doWrite at connection.go
```

### 0x10 dubbo-go-sample client 接收响应

接上述 0x01 发送数据之后会阻塞在 Client grpc 底层读IO中, Layotto经过上述一些列处理层层返回数据激活Client底层Read IO，具体流程如下。

```go
google.golang.org/grpc.(*ClientConn).Invoke at call.go
=>
google.golang.org/grpc.(*ClientConn).Invoke at call.go
=>
google.golang.org/grpc.(*clientStream).RecvMsg at stream.go
=>
google.golang.org/grpc.(*clientStream).withRetry at stream.go
=>
google.golang.org/grpc.(*csAttempt).recvMsg at stream.go
=>
google.golang.org/grpc.recvAndDecompress at rpc_util.go
=>
google.golang.org/grpc.recv at rpc_util.go
=>
google.golang.org/grpc.(*parser).recvMsg at rpc_util.go
=>
google.golang.org/grpc.(*csAttempt).recvMsg at stream.go
func (p *parser) recvMsg(maxReceiveMessageSize int) (pf payloadFormat, msg []byte, err error) {
    if _, err := p.r.Read(p.header[:]); err != nil {
        return 0, nil, err
    }
    ...
}
```

最终收到返回数据：
{"jsonrpc":"2.0","id":9527,"result":{"id":"113","name":"Moorse","age":30,"time":703394193,"sex":"MAN"}}

## 总结
Layotto RPC 处理流程涉及 GRPC、Dapr、Mosn 等相关的知识，整体流程较长，不过单纯看 Layotto 针对 Mosn 抽象的轻量 RPC 框架还是比较清晰和简单的，与 Mosn 集成的方式也比较新颖，值得进一步研读。至此 Layotto RPC 请求处理就分析完了，时间有限，没有进行一些更全面和深入的剖析，如有纰漏之处，欢迎指正，联系方式：rayo.wangzl@gmail.com。另外在此也希望大家能踊跃参与源码分析和开源社区来，一起学习，共同进步。


