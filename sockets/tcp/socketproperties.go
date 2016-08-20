package tcp

import "github.com/gopherjs/gopherjs/js"

type SocketProperties struct {
	Persistent bool
	Name       string
	BufferSize int
}

func (s *SocketProperties) fromJS(o *js.Object) error {
	s.Persistent = o.Get("persistent").Bool()
	s.Name = o.Get("name").String()
	s.BufferSize = o.Get("bufferSize").Int()
	return nil
}

func (s *SocketProperties) toJS(o *js.Object) error {
	o.Set("persistent", s.Persistent)
	o.Set("name", s.Name)
	o.Set("bufferSize", s.BufferSize)
	return nil
}
