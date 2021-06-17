package channel

import (
	"container/list"
	"context"
	"errors"
	"net"
	"sync"

	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/utils"
)

var (
	connpoolTimeout = errors.New("connection pool timeout")
)

type wrapConn struct {
	net.Conn
	buf   buffer.IoBuffer
	state interface{}

	closeOnce sync.Once
	closeChan chan struct{}
}

func (w *wrapConn) Close() error {
	w.closeOnce.Do(func() {
		w.Conn.Close()
		close(w.closeChan)
	})
	return nil
}

// im-memory fake conn pool
func newConnPool(
	maxActive int,
	dialFunc func() (net.Conn, error),
	stateFunc func() interface{},
	onDataFunc func(*wrapConn) error) *connPool {

	p := &connPool{
		maxActive:  maxActive,
		dialFunc:   dialFunc,
		stateFunc:  stateFunc,
		onDataFunc: onDataFunc,
		sema:       make(chan struct{}, maxActive),
		free:       list.New(),
	}
	return p
}

type connPool struct {
	maxActive  int
	dialFunc   func() (net.Conn, error)
	stateFunc  func() interface{}
	onDataFunc func(*wrapConn) error

	sema chan struct{}
	mu   sync.Mutex
	free *list.List
}

func (p *connPool) Get(ctx context.Context) (*wrapConn, error) {
	if err := p.waitTurn(ctx); err != nil {
		return nil, err
	}

	p.mu.Lock()
	// get free conn
	if ele := p.free.Front(); ele != nil {
		p.free.Remove(ele)
		p.mu.Unlock()
		wc := ele.Value.(*wrapConn)
		return wc, nil
	}
	p.mu.Unlock()

	// create new conn
	c, err := p.dialFunc()
	if err != nil {
		p.freeTurn()
		return nil, err
	}
	wc := &wrapConn{Conn: c, closeChan: make(chan struct{})}
	if p.stateFunc != nil {
		wc.state = p.stateFunc()
	}
	if p.onDataFunc != nil {
		utils.GoWithRecover(func() {
			p.readloop(wc)
		}, nil)
	}
	return wc, nil
}

func (p *connPool) Put(c *wrapConn, close bool) {
	if close {
		c.Close()
		p.freeTurn()
		return
	}

	p.mu.Lock()
	if p.free.Len() < p.maxActive {
		p.free.PushBack(c)
		p.mu.Unlock()
	} else {
		p.mu.Unlock()
		c.Close()
	}
	p.freeTurn()
}

func (p *connPool) readloop(c *wrapConn) {
	defer c.Close()

	c.buf = buffer.NewIoBuffer(16 * 1024)
	for {
		_, err := c.buf.ReadOnce(c)
		if err != nil {
			log.DefaultLogger.Errorf("[runtime][rpc]connpool readloop err: %s", err.Error())
			return
		}
		if err = p.onDataFunc(c); err != nil {
			log.DefaultLogger.Errorf("[runtime][rpc]connpool onData err: %s", err.Error())
			return
		}
	}
}

func (p *connPool) waitTurn(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return connpoolTimeout
	case p.sema <- struct{}{}:
		return nil
	}
}

func (p *connPool) freeTurn() {
	<-p.sema
}
