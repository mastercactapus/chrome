package tcpserver

import "github.com/gopherjs/gopherjs/js"

type SocketInfo struct {
	SocketID     int
	Persistent   bool
	Name         string
	Paused       bool
	LocalAddress string
	LocalPort    int
}

func (s *SocketInfo) fromJS(o *js.Object) error {
	s.SocketID = o.Get("socketId").Int()
	s.Persistent = o.Get("persistent").Bool()
	s.Name = o.Get("name").String()
	s.Paused = o.Get("paused").Bool()
	s.LocalAddress = o.Get("localAddress").String()
	s.LocalPort = o.Get("localPort").Int()
	return nil
}

func (s *SocketInfo) toJS(o *js.Object) error {
	o.Set("socketId", s.SocketID)
	o.Set("persistent", s.Persistent)
	o.Set("name", s.Name)
	o.Set("paused", s.Paused)
	o.Set("localAddress", s.LocalAddress)
	o.Set("localPort", s.LocalPort)
	return nil
}
