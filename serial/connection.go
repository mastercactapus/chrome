package serial

import (
	"errors"
	"io"

	"github.com/gopherjs/gopherjs/js"
	"github.com/mastercactapus/chrome/internal/util"
)

type WriteError string

const (
	WriteErrorDisconnected WriteError = "disconnected"
	WriteErrorPending      WriteError = "pending"
	WriteErrorTimeout      WriteError = "timeout"
	WriteErrorSystemError  WriteError = "system_error"
	WriteErrorNone         WriteError = ""
)

func (w WriteError) Error() string {
	return string(w)
}

type ReadError string

const (
	ReadErrorDisconnected   ReadError = "disconnected"
	ReadErrorTimeout        ReadError = "timeout"
	ReadErrorDeviceLost     ReadError = "device_lost"
	ReadErrorBreak          ReadError = "break"
	ReadErrorFrameError     ReadError = "frame_error"
	ReadErrorOverrun        ReadError = "overrun"
	ReadErrorBufferOverflow ReadError = "buffer_overflow"
	ReadErrorParityError    ReadError = "parity_error"
	ReadErrorSystemError    ReadError = "system_error"
)

func (r ReadError) Error() string {
	return string(r)
}

type Connection struct {
	ConnectionID int
	r            *io.PipeReader
	w            *io.PipeWriter
	eRecv        *util.EventListener
	eRecvErr     *util.EventListener
	closed       bool
}

func Connect(path string, opts *ConnectionOptions) (*Connection, error) {
	var args []*js.Object
	var err error
	if opts == nil {
		args, err = util.Call("chrome.serial.connect", path)
	} else {
		jsOpts := new(js.Object)
		opts.toJS(jsOpts)
		args, err = util.Call("chrome.serial.connect", path, jsOpts)
	}

	if err != nil {
		return nil, err
	}

	info := new(ConnectionInfo)
	err = info.fromJS(args[0])
	if err != nil {
		return nil, err
	}
	return NewConnection(info.ConnectionID), nil
}

func NewConnection(connectionID int) *Connection {
	r, w := io.Pipe()
	c := &Connection{
		ConnectionID: connectionID,
		r:            r,
		w:            w,
		eRecv:        util.NewEventListener("chrome.serial.onReceive"),
		eRecvErr:     util.NewEventListener("chrome.serial.OnReceiveError"),
	}
	go c.loop()
	return c
}

func (c *Connection) loop() {
	var o []*js.Object
mainLoop:
	for {
		select {
		case o = <-c.eRecv.C:
			if o == nil {
				break mainLoop
			}
			if o[0].Get("connectionId").Int() != c.ConnectionID {
				continue
			}
			c.w.Write(js.Global.Get("Uint8Array").New(o[0].Get("data")).Interface().([]byte))
		case o = <-c.eRecvErr.C:
			if o == nil {
				break mainLoop
			}
			if o[0].Get("connectionId").Int() != c.ConnectionID {
				continue
			}
			c.w.CloseWithError(ReadError(o[0].Get("error").String()))
			c.Close()
			break mainLoop
		}
	}
}

// GetInfo is analogous to [chrome.serial.getInfo](https://developer.chrome.com/apps/serial#method-getInfo)
func (c *Connection) GetInfo() (*ConnectionInfo, error) {
	args, err := util.Call("chrome.serial.getInfo", c.ConnectionID)
	if err != nil {
		return nil, err
	}
	info := new(ConnectionInfo)
	return info, info.fromJS(args[0])
}

// Update is analogous to [chrome.serial.update](https://developer.chrome.com/apps/serial#method-update)
func (c *Connection) Update(opt *ConnectionOptions) error {
	if c.closed {
		return errors.New("closed")
	}
	jsOpts := new(js.Object)
	opt.toJS(jsOpts)
	args, err := util.Call("chrome.serial.update", c.ConnectionID, jsOpts)
	if err != nil {
		return err
	}
	if args[0].Bool() {
		return nil
	}

	return errors.New("update failed")
}

// Flush is analogous to [chrome.serial.flush](https://developer.chrome.com/apps/serial#method-flush)
func (c *Connection) Flush() error {
	if c.closed {
		return errors.New("closed")
	}
	args, err := util.Call("chrome.serial.flush", c.ConnectionID)
	if err != nil {
		return err
	}
	if args[0].Bool() {
		return nil
	}

	return errors.New("flush failed")
}
func (c *Connection) Write(p []byte) (int, error) {
	if c.closed {
		return 0, errors.New("closed")
	}
	args, err := util.Call("chrome.serial.send", c.ConnectionID, js.NewArrayBuffer(p))
	n := args[0].Get("bytesSent").Int()
	if err != nil {
		return n, err
	}

	err = WriteError(args[0].Get("error").String())
	if err != WriteErrorNone {
		return n, err
	}
	return n, nil
}
func (c *Connection) Read(p []byte) (int, error) {
	return c.r.Read(p)
}
func (c *Connection) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	c.eRecv.Close()
	c.eRecvErr.Close()
	util.Call("chrome.serial.disconnect", c.ConnectionID)
	c.w.Close()
	return nil
}
