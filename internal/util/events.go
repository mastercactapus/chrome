package util

import (
	"errors"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

type EventListener struct {
	C    chan []*js.Object
	obj  *js.Object
	fn   *js.Object
	name string
}

func NewEventListener(eventName string) *EventListener {
	parts := strings.Split(eventName, ".")
	var o *js.Object = js.Global
	for _, key := range parts {
		o = o.Get(key)
	}
	l := &EventListener{
		C:    make(chan []*js.Object, 1),
		obj:  o,
		name: eventName,
	}
	l.fn = js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		l.C <- args
		return nil
	})
	if o == js.Undefined {
		panic(eventName + " is undefined. Check your manifest.json for permissions")
	}
	o.Call("addListener", l.fn)
	return l
}
func (e *EventListener) Close() error {
	if e.obj == nil {
		return errors.New("already closed")
	}
	e.obj.Call("removeListener", e.fn)
	close(e.C)
	e.obj = nil
	e.fn = nil
	return nil
}
