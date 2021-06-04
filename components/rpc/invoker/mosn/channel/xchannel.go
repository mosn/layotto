package channel

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/layotto/layotto/components/rpc"
	"github.com/layotto/layotto/components/rpc/invoker/mosn/transport_protocol"
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
	proto.Init(config.Ext)
	m := &xChannel{
		listenerName: config.Listener,
		proto:        proto,
	}
	return m, nil
}

type xChannel struct {
	listenerName string
	proto        transport_protocol.TransportProtocol

	streamGuard sync.RWMutex
	xconn       *xConn
}

func (m *xChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	duration := time.Duration(req.Timeout) * time.Millisecond
	deadline := time.Now().Add(duration)
	timer := time.NewTimer(duration)
	defer timer.Stop()

	conn, err := m.getConn()
	if err != nil {
		return nil, err
	}
	respChan, closeChan, err := conn.WriteRequest(req, deadline)
	if err != nil {
		m.destroyConn(conn)
		return nil, err
	}

	select {
	case resp := <-respChan:
		return m.proto.FromFrame(resp)
	case <-timer.C:
		m.destroyConn(conn)
		return nil, ErrTimeout
	case <-closeChan:
		m.destroyConn(conn)
		return nil, ErrConnClosed
	}
}

func (m *xChannel) getConn() (*xConn, error) {
	m.streamGuard.RLock()
	sc := m.xconn
	m.streamGuard.RUnlock()
	if sc != nil {
		return sc, nil
	}

	m.streamGuard.Lock()
	defer m.streamGuard.Unlock()
	if m.xconn != nil {
		return m.xconn, nil
	}

	local, remote := net.Pipe()
	localTcpConn := &fakeTcpConn{c: local}
	remoteTcpConn := &fakeTcpConn{c: remote}
	if err := acceptFunc(remoteTcpConn, m.listenerName); err != nil {
		return nil, err
	}

	m.xconn = newXConn(m.proto, localTcpConn)
	go m.xconn.ReadResponse()
	return m.xconn, nil
}

func (m *xChannel) destroyConn(oldConn *xConn) {
	m.streamGuard.RLock()
	sc := m.xconn
	m.streamGuard.RUnlock()
	if oldConn != sc {
		return
	}

	m.streamGuard.Lock()
	defer m.streamGuard.Unlock()

	if oldConn != m.xconn {
		return
	}
	m.xconn.Close()
	m.xconn = nil
}

func newXConn(proto transport_protocol.TransportProtocol, conn net.Conn) *xConn {
	return &xConn{
		proto:     proto,
		conn:      conn,
		closeChan: make(chan struct{}),
		respChans: map[uint64]chan api.XRespFrame{},
	}
}

type xConn struct {
	reqId uint64
	proto transport_protocol.TransportProtocol

	closeOnce sync.Once
	closeChan chan struct{}

	guard     sync.Mutex
	conn      net.Conn
	respChans map[uint64]chan api.XRespFrame
}

func (s *xConn) WriteRequest(req *rpc.RPCRequest, ddl time.Time) (chan api.XRespFrame, chan struct{}, error) {
	// encode request
	frame := s.proto.ToFrame(req)
	id := atomic.AddUint64(&s.reqId, 1)
	frame.SetRequestId(id)
	buf, encErr := s.proto.Encode(req.Ctx, frame)
	if encErr != nil {
		return nil, nil, encErr
	}
	respChan := make(chan api.XRespFrame, 1)

	s.guard.Lock()
	defer s.guard.Unlock()

	// set timeout
	if err := s.conn.SetWriteDeadline(ddl); err != nil {
		return nil, nil, err
	}

	// register response channel
	s.respChans[id] = respChan

	// write packet
	if _, err := s.conn.Write(buf.Bytes()); err != nil {
		delete(s.respChans, id)
		return nil, nil, err
	}

	return respChan, s.closeChan, nil
}

func (s *xConn) ReadResponse() {
	defer s.Close()
	buf := buffer.NewIoBuffer(16 * 1024)

	for {
		if _, err := buf.ReadOnce(s.conn); err != nil {
			log.DefaultLogger.Errorf("[runtime][rpc]xchannel read err: %s", err.Error())
			return
		}

		for {
			iframe, err := s.proto.Decode(context.TODO(), buf)
			if err != nil {
				log.DefaultLogger.Errorf("[runtime][rpc]xchannel decode err: %s", err.Error())
				return
			}

			if iframe == nil {
				break
			}

			frame, ok := iframe.(api.XRespFrame)
			if !ok {
				log.DefaultLogger.Errorf("[runtime][rpc]xchannel type not XRespFrame")
				return
			}

			reqID := frame.GetRequestId()
			s.guard.Lock()
			notifyChan, ok := s.respChans[reqID]
			if ok {
				delete(s.respChans, reqID)
			}
			s.guard.Unlock()
			if ok {
				notifyChan <- frame
			}
		}
	}
}

func (s *xConn) Close() {
	s.closeOnce.Do(func() {
		close(s.closeChan)
		s.conn.Close()
	})
}
