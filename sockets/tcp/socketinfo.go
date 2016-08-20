package tcp

import "github.com/gopherjs/gopherjs/js"

type SocketInfo struct {
	SocketID     int
	Persistent   bool
	Name         string
	BufferSize   int
	Paused       bool
	Connected    bool
	LocalAddress string
	LocalPort    int
	PeerAddress  string
	PeerPort     int
}

func (s *SocketInfo) fromJS(o *js.Object) error {
	s.SocketID = o.Get("socketId").Int()
	s.Persistent = o.Get("persistent").Bool()
	s.Name = o.Get("name").String()
	s.BufferSize = o.Get("bufferSize").Int()
	s.Paused = o.Get("paused").Bool()
	s.Connected = o.Get("connected").Bool()
	s.LocalAddress = o.Get("localAddress").String()
	s.LocalPort = o.Get("localPort").Int()
	s.PeerAddress = o.Get("peerAddress").String()
	s.PeerPort = o.Get("peerPort").Int()
	return nil
}
func (s *SocketInfo) toJS(o *js.Object) error {
	o.Set("socketId", s.SocketID)
	o.Set("persistent", s.Persistent)
	o.Set("name", s.Name)
	o.Set("bufferSize", s.BufferSize)
	o.Set("paused", s.Paused)
	o.Set("connected", s.Connected)
	o.Set("localAddress", s.LocalAddress)
	o.Set("localPort", s.LocalPort)
	o.Set("peerAddress", s.PeerAddress)
	o.Set("peerPort", s.PeerPort)
	return nil
}
