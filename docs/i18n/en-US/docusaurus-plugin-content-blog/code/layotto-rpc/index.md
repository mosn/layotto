# Layotto Source Parsing — Processing RPC requests

> This paper is based on the Dubbo Json RPC as an example of the Layotto RPC processing.
>
> by：[Wang Zhilong](https://github.com/rayowang) | 21April 2022

- [overview](#overview)
- [source analysis](#source analysis)
  - [0x00 Layotto initialize RPC](#_0x00-layotto-initializ-rpc)
  - [0x01 Dubbo-go-sample client to request request] (#_0x01-dubbo-go-sample-client-request request)
  - [0x02 Mosn EventLoop Reader Processing Request Data](#_0x02-mosn-eventloop-read processing request)
  - [0x03 Grpc Sever as NetworkFilter to process requests](#_0x03-grpc-sever as -networkfilter-process requests)
  - [0x04 Layotto send RPC requests and write to Local Virtual Connections](#_0x04-layotto -rpc-request and write to -local-virtual connection)
  - [0x05 Mosn reads Remote and executes Filter and Proxy Forwarding](#_0x05-mosn-read-remote-remote--and --filter-and proxy forward)
  - [0x06 Dubbo-go-sample server received response to request return] (#_0x06-dubbo-go-sample-server-received response return)
  - [0x07 Mosn Framework handles response and writes back to Remote Virtual Connections](#_0x07-mosn-Framework handles response and -remote-virtual connection)
  - [0x08 Layotto receive RPC responses and read Local Virtual Connections](#_0x08-layotto-receive-rpc-response and read -local-virtual connection)
  - [0x09 Grpc Sever processed data frames returned to client](#_0x09-grpc-sever processed frame returned to client)
  - [0x10 Dubbo-go-sample client receiving response](#_0x10-dubbo-go-sample-client-receiving response)
- [summary](#summary)

## General description

Layotto has a clear and rich semiconductor API as a distributed prototype collection of prototype language distinguished from the network proxy service Mesh and using standard protocol API, which is part of the RPC API.Through RPC API app developers can interact with local Layotto instances of applications that also use the Sidecar architecture, thereby indirectly calling different service methods and using built-in capabilities to perform distributive tracking and diagnosis, traffic control, error handling, secure links, etc.and Layotto is based on the Grpc handler design, using the X-Protocol protocol for secure and reliable communications, except for Http/Grpc communications with other services.As shown in the following code, the RPC API interface is in line with Dapr and is available for RPC calls through the Grpc interface InvokeService.

```go
type DaprClient interface {
    // Invokes a method on a remote Dapr app.
    InvokeService(ctx context.Context, in *InvokeServiceRequest, opts ...grpc.CallOption) (*v1.InvokeResponse, error)
    ...
}
```

## Source analysis

For ease of understanding, from outside to inside, from inside to outside, from flow to source code, that is, from Client, through one layer of logic to the Server receiving a return response to requests, from another layer of return to client, and from one layer of analysis of Layotto RPC processes, split into 10 steps.Also, since the content of Gypc Client and Server handshakes and interactions is not the focus of this paper, the analysis is relatively brief and the other steps are relatively detailed and one can move directly from the directory to the corresponding step depending on his or her case.

Note：based on commit hash：1d2bed68c3b2372c34a12aeed41be125a4fdd15a

### 0x00 Layotto initialize RPC

Layotto starts the process involves a large number of processes in which only the initialization of the process related to RPC and described below is analyzed because Layotto is based on Mosn and is therefore starting from the Main function, urfave/cli library calls Mosn StageManager Mos, thus initializing GrpcServer in Mosn NetworkFilter as follows.

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

### 0x01 Dubbo-go-sample client request

Follow the example of [Dubbo Json Rpc Example](https://mosn.io/layotto/en-US/docs/start/rpc/dub_json_rpc)

```shell
go un demo/rpc/dubbo_json_rpc/dub_json_client/client.go -d '{"jsonrpc": "2.0", "method":"GetUser", "params":["A003"],"id":9527}'
```

Use Layotto for App Gypc API InvokeService initiate RPC calls, data filling and connecting processes leading to the dispatch of data to Layotto via SendMsg in Grpc clientStream, as follows.

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

### 0x02 Mosn EventLoop Reader Processing Request Data

The kernel from Layotto mentioned above is a mock-up of Mosn, so when network connection data arrives, it will first be read and written at the L4 network level in Mosn as follows.

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

In the startRWLoop method, we can see that two separate protocols will be opened to deal with reading and writing operations on the connection: startReadLoop and startWriteLoop; the following streams will be made in startReadLoop; the data read at the network level will be handled by the filterManager filter chain, as follows.

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

### 0x03 Grpc Sever processed requests as NetworkFilter

Reading data from the original connection in the first phase will enter the Grpc Serve handling, the Serve method will use the net.Listener listener, each time a new protocol is launched to handle the new connection (handleRawCon), and a RPC call based on Http2-based transport will be set out below.

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

### 0x04 Layotto send RPC requests and write to local virtual connections

The 0x03 process follows Runtime_InvokeService_Handler, converted from the GRPC Default API to Dapr API, entering the light RPC framework provided by Layotto in Mosn, as follows.

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

The creation of a new connection in the second step above requires attention by calling dialFunc func() in the protocol that initialized the init phase (net.Conn, error), because the configuration interacted with Mosn with Http protocols, this is newHttpChanel, which is currently supported by the Bolt, Dubbo et al.

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

### 0x05 Mosn read Remote and execute Filter and proxy forwarding

(1) Similar to 0x02, filtermanager executes the filter processing phase where proxy forwarding is made in proxy with the following code.

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

(2) ServerStreamConnection.serve listens and handles requests to downstream OnReceive, as described below.

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

(3) The above ScheduleAuto schedule, after processing the reveive of downstream Stream, processing upstam Request, as well as an application with an application from the network layer, eventually sending data from connection.Write and entering WaitNotify phases, as detailed below.

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

### 0x06 Dubbo-go-sample server received request return response

Here is a dubo-go-sample server handling, leave it now, post log messages and check the source code by interested classes.

```
[2022-04-18/21:03:03:18 github.com/apache/dub-go-samples/rpc/jsonrpc/go-server/pkg.(*UserProvider2).GetUser: user_provider2.go: 53] userID: "A003"
[2022-04-18/21:03:18 github.com/apache/dub-go-samples/rpc/jsonrpc/go-server/pkg. (*UserProvider2).GetUser: user_provider2.go: 56] rsp:&pkg.User{ID:"113", Name:"Moorse", Age:30, sex:0, Birth:703391943, Sex:"MAN"MAN"}
```

### 0x07 Mosn framework handles responses and writes back to Remote Virtual Connection

After the third phase of 0x05 above, the response logic goes into the UpRecvData phase of the reveive cycle phase through a series of final response writing back to the remote virtual connection at 0x04, as follows.

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

### 0x08 Layotto receive RPC responses and read Local Virtual Connection

Readloop Reading IO, activated by 0x04 above, is activated from connection read data from Mosn and then forwarded to the hstate pipe to return to the request process, as follows.

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

### 0x09 Grpc Sever processed data frames returned to clients

Grpc does not write data directly to connections, but uses a systray loop to fetch frames from a cache structure and write them back to the client, as follows.

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

### 0x10 dubbo-go-sample customer received response

The transmission of data from 0x01 above will be blocked in the client grpc bottom reading, and Layotto returns data from some of the processing layers above to enable ClientBottom Read IO, as follows.

```go
google.golang.org/grpc.(*ClientCon). Invoke at call.go
=>
google.golang.org/grpc.(*ClientCon). Invoke at call.go
=>
google.golang.org/grpc.(*clientStream). RecvMsg at stream. o
=>
google.golang.org/grpc.(*clientStream).withRetry at stream.go
=>
google.golang.org/grpc.(*csAtempt.recvMsg at stream.go
=>
google.golang.org/grpc.recvAndDecompress at rpc_util. o
=>
google.golang.org/grpc.recv at rpc_util.go
=>
google.golang.org/grpc.(*parser).recvMsg at rpc_util.go
=>
google.golang.org/grpc.(*csAttempt).recvMsg at stream. o
func (p *parser) recvMsg(maxReceiveMessageSize int) (pf payloadFormat, msg []byte, err error) LO
    if _, err := p. .Read(p.header[:]); err != nil {
        return 0, nil, err
    }
    ...
}
```

Last returned data：

```json
{"jsonrpc": "2.0", "id":9527, "result":{"id":"113", "name":"Moorse", "age":30,"time":703394193,"sex":"MAN"}}
```

## Summary

The Layotto RPC process involves knowledge related to GRPC, Dapr, Mosn and others, and the overall process is lengthy, although it is clearer and simpler simply to see Layotto for Mosn an abstract lightweight RPC framework and is more innovative and useful for further study.Here Layotto RPC requests are analyzed and time-limited without some more comprehensive and in-depth profiles, such as defects, welcome contact：rayo.wangzl@gmail.com.It is also hoped that there will be greater participation in source analysis and open source communities, learning together and making progress together.
