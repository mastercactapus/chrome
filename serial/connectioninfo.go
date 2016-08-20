package serial

import "github.com/gopherjs/gopherjs/js"

type ConnectionInfo struct {
	ConnectionID int
	Paused       bool
	ConnectionOptions
}

func (s *ConnectionInfo) fromJS(o *js.Object) error {
	s.ConnectionID = o.Get("connectionId").Int()
	s.Paused = o.Get("paused").Bool()
	return s.ConnectionOptions.fromJS(o)
}

func (s *ConnectionInfo) toJS(o *js.Object) error {
	o.Set("connectionId", s.ConnectionID)
	o.Set("paused", s.Paused)
	return s.ConnectionOptions.toJS(o)
}
