package channel

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/layotto/layotto/pkg/services/rpc"
	"github.com/layotto/layotto/pkg/services/rpc/invoker/mosn/transport_protocol"
	"mosn.io/api"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
)

func init() {
	RegistChannel("bolt", newXChannel)
	RegistChannel("boltv2", newXChannel)
	RegistChannel("dubbo", newXChannel)
}

func newXChannel(config ChannelConfig) (rpc.Channel, error) {
	proto := transport_protocol.GetProtocol(config.Protocol)
	if proto == nil {
		return nil, fmt.Errorf("protocol %s not found", proto)
	}
	m := &xChannel{
		listenerName: config.Listener,
		proto:        proto,
		reqIdMap:     map[uint64]chan api.XRespFrame{},
	}
	return m, nil
}

type xChannel struct {
	reqId        uint64
	listenerName string
	proto        transport_protocol.TransportProtocol

	connGuard sync.RWMutex
	conn      *memConn

	reqIdMapGuard sync.Mutex
	reqIdMap      map[uint64]chan api.XRespFrame
}

func (m *xChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	reqFrame := m.proto.ToFrame(req)
	id := atomic.AddUint64(&m.reqId, 1)
	reqFrame.SetRequestId(id)
	buf, encErr := m.proto.Encode(req.Ctx, reqFrame)
	if encErr != nil {
		return nil, encErr
	}
	defer buffer.PutIoBuffer(buf)

	respChan := make(chan api.XRespFrame, 1)
	m.reqIdMapGuard.Lock()
	m.reqIdMap[id] = respChan
	m.reqIdMapGuard.Unlock()

	var (
		conn *memConn
		err  error
	)

	for i := 0; i < 4; i++ {
		conn, err = m.getConn()
		if err != nil {
			m.deleteReqId(id)
			return nil, err
		}

		if _, err = conn.write(buf.Bytes()); err == nil {
			break
		}

		m.destroyConn(conn)
	}
	if err != nil {
		m.deleteReqId(id)
		return nil, err
	}

	timer := time.NewTimer(time.Duration(req.Timeout) * time.Millisecond)
	defer timer.Stop()

	select {
	case resp := <-respChan:
		return m.proto.FromFrame(resp)
	case <-timer.C:
		m.deleteReqId(id)
		return nil, ErrTimeout
	case <-conn.closeChan():
		m.deleteReqId(id)
		return nil, ErrConnClosed
	}
}

func (m *xChannel) deleteReqId(id uint64) {
	m.reqIdMapGuard.Lock()
	delete(m.reqIdMap, id)
	m.reqIdMapGuard.Unlock()
}

func (m *xChannel) acceptConn(conn *memConn) error {
	if err := acceptFunc(conn, m.listenerName); err != nil {
		return err
	}
	m.conn = conn
	return nil
}

func (m *xChannel) getConn() (*memConn, error) {
	m.connGuard.RLock()
	conn := m.conn
	m.connGuard.RUnlock()
	if conn != nil {
		return conn, nil
	}

	m.connGuard.Lock()
	defer m.connGuard.Unlock()
	if m.conn != nil {
		return m.conn, nil
	}

	conn = newMemConn(m)
	if err := m.acceptConn(conn); err != nil {
		return nil, err
	}

	return conn, nil
}

func (m *xChannel) destroyConn(oldConn *memConn) {
	m.connGuard.RLock()
	conn := m.conn
	m.connGuard.RUnlock()
	if oldConn != conn {
		return
	}

	m.connGuard.Lock()
	defer m.connGuard.Unlock()

	if oldConn != m.conn {
		return
	}
	m.conn = nil
}

// xprotocol write complete packet at once
func (m *xChannel) handleConnResp(b []byte) error {
	iframe, err := m.proto.Decode(context.TODO(), buffer.NewIoBufferBytes(b))
	if err != nil {
		return err
	}

	if iframe == nil {
		return errors.New("xprotocol decode return nil")
	}

	frame, ok := iframe.(api.XRespFrame)
	if !ok {
		return errors.New("xprotocol frame not XRespFrame")
	}

	reqID := frame.GetRequestId()

	m.reqIdMapGuard.Lock()
	notifyChan, ok := m.reqIdMap[reqID]
	if ok {
		delete(m.reqIdMap, reqID)
	}
	m.reqIdMapGuard.Unlock()

	if ok {
		notifyChan <- frame
	}

	return nil
}
func newMemConn(c *xChannel) *memConn {
	r, w := io.Pipe()
	return &memConn{wPipe: w, rPipe: r, closed: make(chan struct{}), mosnChannel: c}
}

type memConn struct {
	closeOnce sync.Once
	closed    chan struct{}
	rPipe     *io.PipeReader
	wPipe     *io.PipeWriter

	mosnChannel *xChannel
}

func (g *memConn) Read(b []byte) (n int, err error) {
	return g.rPipe.Read(b)
}

func (g *memConn) Write(b []byte) (n int, err error) {
	select {
	case <-g.closed:
		return 0, errors.New("mem conn closed")
	default:
		err := g.mosnChannel.handleConnResp(b)
		if err != nil {
			log.DefaultLogger.Errorf("[runtime][rpc][xchannel]Write err %s", err.Error())
		}
		return len(b), err
	}
}

func (g *memConn) Close() error {
	g.closeOnce.Do(func() {
		g.rPipe.Close()
		g.wPipe.Close()
		close(g.closed)
		log.DefaultLogger.Warnf("[runtime][rpc][xchannel]mem conn closed")
	})
	return nil
}

func (g *memConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (g *memConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (g *memConn) SetDeadline(t time.Time) error {
	return nil
}

func (g *memConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (g *memConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (g *memConn) write(b []byte) (int, error) {
	return g.wPipe.Write(b)
}

func (g *memConn) closeChan() chan struct{} {
	return g.closed
}
