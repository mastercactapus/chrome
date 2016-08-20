package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// chromeTCPSocket implements the net.Conn interface
type chromeTCPSocket struct {
	socketId int
	r        *io.PipeReader
	w        *io.PipeWriter
	eRecv    *eventListener
	eRecvErr *eventListener
	closeCh  chan struct{}
}

func newChromeSocket(socketId int) *chromeTCPSocket {
	r, w := io.Pipe()
	s := &chromeTCPSocket{
		socketId: socketId,
		r:        r,
		w:        w,
		eRecv:    newEventListener("chrome.sockets.tcp.onReceive"),
		eRecvErr: newEventListener("chrome.sockets.tcp.onReceiveError"),
		closeCh:  make(chan struct{}),
	}
	go s.loop()
	chromeCall("chrome.sockets.tcp.setPaused", socketId, false)
	return s
}
func (s *chromeTCPSocket) loop() {
	var o []*js.Object
mainLoop:
	for {
		select {
		case <-s.closeCh:
			break mainLoop
		case o = <-s.eRecv.C:
			if o[0].Get("socketId").Int() != s.socketId {
				continue
			}
			s.w.Write(js.Global.Get("Uint8Array").New(o[1]).Interface().([]byte))
		case o = <-s.eRecvErr.C:
			if o[0].Get("socketId").Int() != s.socketId {
				continue
			}
			err := fmt.Errorf("receive error, code: %d", o[0].Get("resultCode").Int())
			s.w.CloseWithError(err)
			s.Close()
			break mainLoop
		}
	}
	s.eRecv.Close()
	s.eRecvErr.Close()
}
func (s *chromeTCPSocket) Close() error {
	// wait for closeCh to be picked up before setting socketId to 0
	// so we know the we have broken out of the 'loop'
	//
	// we also want to do that first before we call tcp.close so
	// we don't trigger a receiveErr
	s.closeCh <- struct{}{}
	chromeCall("chrome.sockets.tcp.close", s.socketId)
	s.w.Close()
	s.socketId = 0
}

func (s *chromeTCPSocket) Write(p []byte) (int, error) {
	if s.socketId == 0 {
		return 0, errors.New("socket closed")
	}
	args, err := chromeCall("chrome.sockets.tcp.send", s.socketId, js.NewArrayBuffer(p))
	n := args[0].Get("bytesSent").Int()
	if err != nil {
		return n, err
	}
	if args[0].Get("resultCode").Int() < 0 {
		return n, fmt.Errorf("write failed, error code: %d", args[0].Get("resultCode").Int())
	}
	return n, nil
}

func (s *chromeTCPSocket) Read(p []byte) (int, error) {
	return s.r.Read(p)
}
func (s *chromeTCPSocket) LocalAddr() net.Addr {
	if s.socketId == 0 {
		return nil
	}
	args, err := chromeCall("chrome.sockets.tcp.getInfo", s.socketId)
	if err != nil {
		panic(err)
	}
	o := args[0]
	return &net.TCPAddr{
		IP:   net.ParseIP(o.Get("localAddress").String()),
		Port: o.Get("localPort").Int(),
	}
}
func (s *chromeTCPSocket) RemoteAddr() net.Addr {
	if s.socketId == 0 {
		return nil
	}
	args, err := chromeCall("chrome.sockets.tcp.getInfo", s.socketId)
	if err != nil {
		panic(err)
	}
	o := args[0]
	return &net.TCPAddr{
		IP:   net.ParseIP(o.Get("remoteAddress").String()),
		Port: o.Get("remotePort").Int(),
	}
}

func (s *chromeTCPSocket) SetDeadline(t time.Time) error {
	return nil
}
func (s *chromeTCPSocket) SetReadDeadline(t time.Time) error {
	return nil
}
func (s *chromeTCPSocket) SetWriteDeadline(t time.Time) error {
	return nil
}
