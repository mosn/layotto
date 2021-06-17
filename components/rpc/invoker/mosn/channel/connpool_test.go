package channel

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPut(t *testing.T) {
	active := 2
	p := newConnPool(
		active,
		func() (net.Conn, error) {
			p, _ := net.Pipe()
			return &fakeTcpConn{c: p}, nil
		},
		func() interface{} {
			return nil
		}, func(conn *wrapConn) error {
			return nil
		},
	)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	c1, err := p.Get(ctx)
	assert.Nil(t, err)
	c2, err := p.Get(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, p.free.Len())

	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
	assert.Equal(t, 0, p.free.Len())

	_, err = p.Get(ctx)
	t.Log(err)
	assert.Error(t, err)

	p.Put(c1, false)
	p.Put(c2, false)
	assert.Equal(t, active, p.free.Len())
}

type conns struct {
	sync.RWMutex
	conns []net.Conn
}

func (s *conns) add(c net.Conn) {
	s.Lock()
	s.conns = append(s.conns, c)
	s.Unlock()
}

func (s *conns) close() {
	s.Lock()
	for _, c := range s.conns {
		c.Close()
	}
	s.conns = nil
	s.Unlock()
}

func TestDeadconnRenew(t *testing.T) {
	conns := &conns{}
	active := 1
	p := newConnPool(
		active,
		func() (net.Conn, error) {
			p1, p2 := net.Pipe()
			conns.add(p2)
			return &fakeTcpConn{c: p1}, nil
		},
		func() interface{} {
			return nil
		}, func(conn *wrapConn) error {
			return nil
		},
	)

	c1, err := p.Get(context.TODO())
	assert.Nil(t, err)

	conns.close()
	p.Put(c1, false)
	c1.close()

	assert.Equal(t, active, p.free.Len())

	c2, err := p.Get(context.TODO())
	assert.Nil(t, err)
	assert.False(t, c1 == c2)
}

func TestPoolConcurrent(t *testing.T) {
	conns := &conns{}
	active := 5
	ch := make(chan error)
	p := newConnPool(
		active,
		func() (net.Conn, error) {
			p1, p2 := net.Pipe()
			conns.add(p2)
			return &fakeTcpConn{c: p1}, nil
		},
		func() interface{} {
			return nil
		}, func(conn *wrapConn) error {
			return <-ch
		},
	)

	actions := []string{
		"Get", "Put",
		"Get", "close",
		"Get", "Put",
		"readclose", "Get", "close",
		"Get", "readclose", "Put",
		"Get", "Put", "readclose",
		"Get", "Put",
		"Get", "close",
		"Get", "Put",
	}

	var wg sync.WaitGroup
	total := 10
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var c *wrapConn
			var err error
			for _, act := range actions {
				switch act {
				case "Get":
					ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
					c, err = p.Get(ctx)
					cancel()
					assert.Nil(t, err)
				case "Put":
					p.Put(c, false)
				case "close":
					p.Put(c, true)
				case "readclose":
					conns.close()
				}
			}
		}()
	}
	wg.Wait()

	assert.True(t, p.free.Len() <= active)
	t.Log(p.free.Len())
}
