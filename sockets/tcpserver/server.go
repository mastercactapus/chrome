package tcpserver

import (
	"errors"
	"fmt"
	"net"

	"github.com/gopherjs/gopherjs/js"
	"github.com/mastercactapus/chrome/internal/util"
	"github.com/mastercactapus/chrome/sockets/tcp"
)

type Server struct {
	SocketID   int
	eAccept    *util.EventListener
	eAcceptErr *util.EventListener
	closeCh    chan struct{}
	closed     bool
}

func Listen(network, laddr string) (*Server, error) {
	if network != "tcp" {
		return nil, errors.New("unsupported network: " + network)
	}

	args, err := util.Call("chrome.sockets.tcpServer.create")
	if err != nil {
		return nil, err
	}

	s := NewServer(args[0].Get("socketId").Int())

	err = s.Listen(laddr)
	if err != nil {
		s.Close()
		return nil, err
	}

	return s, nil
}

func NewServer(socketID int) *Server {
	return &Server{
		SocketID:   socketID,
		eAccept:    util.NewEventListener("chrome.sockets.tcpServer.onAccept"),
		eAcceptErr: util.NewEventListener("chrome.sockets.tcpServer.onAcceptErr"),
		closeCh:    make(chan struct{}, 1),
	}
}

func (s *Server) Accept() (net.Conn, error) {
	return s.AcceptTCP()
}

func (s *Server) AcceptTCP() (*tcp.Connection, error) {
	if s.closed {
		return nil, errors.New("closed")
	}
	var args []*js.Object

	for {
		select {
		case args = <-s.eAccept.C:
			if args == nil {
				return nil, errors.New("closed")
			}
			if args[0].Get("socketId").Int() != s.SocketID {
				continue
			}
			clientSocketId := args[0].Get("clientSocketId").Int()
			// tcpServer pauses new sockets by default, so unpause it once constructed
			defer util.Call("chrome.sockets.tcp.setPaused", clientSocketId, false)
			return tcp.NewConnection(clientSocketId), nil
		case args = <-s.eAcceptErr.C:
			if args == nil {
				return nil, errors.New("closed")
			}
			if args[0].Get("socketId").Int() != s.SocketID {
				continue
			}
			return nil, fmt.Errorf("accept error code: %d", args[0].Get("resultCode").Int())
		}
	}
}

func (s *Server) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	util.Call("chrome.sockets.tcpServer.close", s.SocketID)
	s.eAccept.Close()
	s.eAcceptErr.Close()
	return nil
}

func (s *Server) Listen(laddr string) error {
	host, port, err := net.SplitHostPort(laddr)
	if err != nil {
		return err
	}
	args, err := util.Call("chrome.sockets.tcpServer.listen", s.SocketID, host, port)
	if err != nil {
		return err
	}
	if args[0].Int() < 0 {
		return fmt.Errorf("listen failed, error code: %d", args[0].Int())
	}
	return nil
}

func (s *Server) GetInfo() (*SocketInfo, error) {
	args, err := util.Call("chrome.sockets.tcpServer.getInfo", s.SocketID)
	if err != nil {
		return nil, err
	}
	info := new(SocketInfo)
	return info, info.fromJS(args[0])
}

func (s *Server) Addr() net.Addr {
	info, err := s.GetInfo()
	if err != nil {
		return nil
	}
	return &net.TCPAddr{
		IP:   net.ParseIP(info.LocalAddress),
		Port: info.LocalPort,
	}
}
