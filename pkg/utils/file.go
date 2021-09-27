package utils

import (
	"io"
	"sync"

	"mosn.io/pkg/buffer"
	"mosn.io/pkg/log"
)

type WrapFileStream struct {
	// Reader channel
	r        chan buffer.IoBuffer
	endRead  chan struct{}
	err      error
	fileOnce sync.Once
}

func NewFileReader() *WrapFileStream {
	return &WrapFileStream{
		r:       make(chan buffer.IoBuffer),
		endRead: make(chan struct{}),
	}
}
func (w *WrapFileStream) Read(b []byte) (int, error) {
	data, ok := <-w.r
	if !ok {
		if w.err == nil {
			return 0, io.EOF
		}
		return 0, w.err
	}
	n := copy(b, data.Bytes())
	data.Drain(n)
	w.endRead <- struct{}{}
	if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
		log.DefaultLogger.Debugf("file read data:  %d", n)
	}
	return n, nil
}

func (w *WrapFileStream) Write(data []byte) (n int, err error) {
	l := len(data)
	buf := buffer.NewIoBufferBytes(data)
	for buf.Len() > 0 {
		w.r <- buf
		<-w.endRead
	}
	return l, nil
}

func (w *WrapFileStream) closeChan() {
	close(w.r)
}
func (w *WrapFileStream) Close() error {
	w.fileOnce.Do(w.closeChan)
	return nil
}

func (w *WrapFileStream) CloseWithErr(err error) {
	w.err = err
	w.fileOnce.Do(w.closeChan)
}
