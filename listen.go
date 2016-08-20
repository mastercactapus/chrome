package chrome

import (
	"errors"
	"fmt"
	"net"
)

// chromeListener implements the net.Listener interface
type chromeListener struct {
	socketId int
	e        *eventListener
}

// Listen is analogous to the net.Listen call, but for use in a chrome app
func Listen(net, laddr string) (net.Listener, error) {

	switch net {
	case "tcp":
		return chromeListenTCP(0, laddr)
	default:
		return nil, errors.New("unsuported network type: " + network)
	}
	return chromeListen(net, laddr)
}

// GetListeners will get all listening sockets for the given network type
func GetListeners(net string) ([]net.Listener, error) {
	switch net {
	case "tcp":
		args, err := chromeCall("chrome.sockets.tcpServer.getSockets")
		if err != nil {
			return nil, err
		}
		l := make([]net.Listener, args[0].Length())
		for i := range l {
			l[i] = &chromeListener{
				socketId: args[0].Index(i).Get("socketId").Int(),
				e:        newEventListener("chrome.sockets.tcpServer.onAccept"),
			}
		}
		return l, nil
	default:
		return nil, errors.New("unsuported network type: " + network)
	}
}

// bind a new tcpServer
func newChromeListenTCP(laddr string) (*chromeListener, error) {
	host, port, err := net.SplitHostPort(laddr)
	if err != nil {
		return nil, err
	}
	args, err := chromeCall("chrome.sockets.tcpServer.create")
	if err != nil {
		return nil, err
	}

	l := &chromeListener{
		socketId: args[0].Int(),
	}
	evt := newEventListener("chrome.sockets.tcpServer.onAccept")
	args, err = chromeCall("chrome.sockets.tcpServer.listen", l.socketId, host, port)
	if err != nil {
		evt.Close()
		return nil, err
	}
	code := args[0].Int()
	if code < 0 {
		evt.Close()
		return nil, fmt.Errorf("listen failed, error code: %d", code)
	}
	return l
}

func (c *chromeListener) Accept() (net.Conn, error) {
	for args := range c.e.C {
		socketId := args[0].Get("socketId").Int()
		if socketId != c.socketId {
			continue
		}

		clientSocketId := args[0].Get("clientSocketId").Int()
		return newChromeSocket(clientSocketId), nil
	}
	return nil, errors.New("closed")
}
func (c *chromeListener) Addr() net.Addr {
	args, err := chromeCall("chrome.sockets.tcpServer.getInfo", c.socketId)
	if err != nil {
		panic(err)
	}
	o := args[0]
	return &net.TCPAddr{
		IP:   net.ParseIP(o.Get("localAddress").String()),
		Port: o.Get("localPort").Int(),
	}
}
func (c *chromeListener) Close() error {
	if c.socketId == 0 {
		return errors.New("already closed")
	}
	c.e.Close()
	chromeCall("chrome.sockets.tcpServer.close", c.socketId)
	c.socketId = 0
}
