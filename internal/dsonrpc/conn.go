package dsonrpc

import (
	"encoding/json"
	"github.com/fpawel/goutils"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"io"
	"net"
	"sync"
)

type readWriteCloser struct {
	conn net.Conn
	w    *io.PipeWriter
	r    *io.PipeReader
	err  error      // last io error
	mu   sync.Mutex // protects err
}

func NewClient(conn net.Conn) *jsonrpc2.Client {
	c := &readWriteCloser{
		conn: conn,
	}
	c.initRead()
	c.initWrite()
	return jsonrpc2.NewClient(c)
}

func (x *readWriteCloser) initRead() {
	var w *io.PipeWriter
	x.r, w = io.Pipe()
	go func() {
		defer func() {
			_ = w.Close()
		}()

		var b16 []byte
		for x.error() == nil {
			b1 := []byte{0}
			if _, err := x.conn.Read(b1); err == nil {
				b16 = append(b16, b1...)
				continue
			}
			b, err := goutils.UTF8FromUTF16(b16)
			if err != nil {
				panic(err)
			}
			_, err = w.Write(b)
			x.setError(err)
			return
		}
	}()
}

func (x *readWriteCloser) initWrite() {
	var r *io.PipeReader
	r, x.w = io.Pipe()
	go func() {
		defer func() {
			_ = r.Close()
		}()

		dec := json.NewDecoder(r)
		var b json.RawMessage

		for err := dec.Decode(&b); err == nil && x.error() == nil; err = dec.Decode(&b) {
			b16 := goutils.UTF16FromString(string(b))
			_, err = x.conn.Write(b16)
			x.setError(err)
		}
		_ = r.Close()
	}()
}

func (x *readWriteCloser) error() (err error) {
	x.mu.Lock()
	err = x.err
	x.mu.Unlock()
	return
}

func (x *readWriteCloser) setError(err error) {
	x.mu.Lock()
	if x.err != nil {
		x.err = err
	}
	x.mu.Unlock()
}

func (x *readWriteCloser) Write(p []byte) (int, error) {
	if x.error() != nil {
		return 0, x.error()
	}
	return x.w.Write(p)
}

func (x *readWriteCloser) Read(p []byte) (int, error) {
	if x.error() != nil {
		return 0, x.error()
	}
	return x.r.Read(p)
}

func (x *readWriteCloser) Close() error {
	_ = x.w.Close()
	_ = x.r.Close()
	return x.conn.Close()
}
