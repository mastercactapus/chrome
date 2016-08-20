package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/mastercactapus/chrome/internal/util"
)

type Connection struct {
	SocketID int
	closed   bool
	r        *io.PipeReader
	w        *io.PipeWriter
	eRecv    *util.EventListener
	eRecvErr *util.EventListener
	closeCh  chan struct{}
}

func NewConnection(socketID int) *Connection {
	r, w := io.Pipe()
	s := &Connection{
		SocketID: socketID,
		r:        r,
		w:        w,
		eRecv:    util.NewEventListener("chrome.sockets.tcp.onReceive"),
		eRecvErr: util.NewEventListener("chrome.sockets.tcp.onReceiveError"),
		closeCh:  make(chan struct{}),
	}
	go s.loop()
	return s
}

func (c *Connection) loop() {
	var o []*js.Object
mainLoop:
	for {
		select {
		case <-c.closeCh:
			break mainLoop
		case o = <-c.eRecv.C:
			if o[0].Get("socketId").Int() != c.SocketID {
				continue
			}
			c.w.Write(js.Global.Get("Uint8Array").New(o[1]).Interface().([]byte))
		case o = <-c.eRecvErr.C:
			if o[0].Get("socketId").Int() != c.SocketID {
				continue
			}
			err := fmt.Errorf("receive error, code: %d", o[0].Get("resultCode").Int())
			c.w.CloseWithError(err)
			c.Close()
			break mainLoop
		}
	}
	c.eRecv.Close()
	c.eRecvErr.Close()
}

func (c *Connection) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true

	// wait for closeCh to be picked up before we call tcp.close
	c.closeCh <- struct{}{}
	util.Call("chrome.sockets.tcp.close", c.SocketID)
	c.w.Close()
	return nil
}

func (c *Connection) Read(p []byte) (int, error) {
	return c.r.Read(p)
}

func (c *Connection) Write(p []byte) (int, error) {
	if c.SocketID == 0 {
		return 0, errors.New("socket closed")
	}
	args, err := util.Call("chrome.sockets.tcp.send", c.SocketID, js.NewArrayBuffer(p))
	n := args[0].Get("bytesSent").Int()
	if err != nil {
		return n, err
	}
	if args[0].Get("resultCode").Int() < 0 {
		return n, fmt.Errorf("write failed, error code: %d", args[0].Get("resultCode").Int())
	}
	return n, nil
}

func (c *Connection) GetInfo() (*SocketInfo, error) {
	args, err := util.Call("chrome.sockets.tcp.getInfo", c.SocketID)
	if err != nil {
		return nil, err
	}
	info := new(SocketInfo)
	return info, info.fromJS(args[0])
}

func (c *Connection) LocalAddr() net.Addr {
	s, err := c.GetInfo()
	if err != nil {
		return nil
	}
	return &net.TCPAddr{
		IP:   net.ParseIP(s.LocalAddress),
		Port: s.LocalPort,
	}
}
func (c *Connection) RemoteAddr() net.Addr {
	s, err := c.GetInfo()
	if err != nil {
		return nil
	}
	return &net.TCPAddr{
		IP:   net.ParseIP(s.PeerAddress),
		Port: s.PeerPort,
	}
}

// SetDeadline does nothing as it is unsupported via the chrome.* APIs
func (c *Connection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline does nothing as it is unsupported via the chrome.* APIs
func (c *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline does nothing as it is unsupported via the chrome.* APIs
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}
