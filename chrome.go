package main

import (
	"errors"
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

func chromeCall(method string, args ...*js.Object) ([]*js.Object, error) {
	parts := strings.Split(method, ".")
	var o *js.Object = js.Global
	for _, key := range parts[:len(parts)-1] {
		o = o.Get(key)
	}
	ch := make(chan []*js.Object, 1)
	var err error
	args = append(args, js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		lastErr := js.Global.Get("chrome").Get("runtime").Get("lastError")
		if lastErr != nil {
			err = errors.New(lastErr.Get("message").String())
		}
		ch <- args
		return nil
	}))
	o.Call(parts[len(parts)-1], args...)

	return <-ch, err
}

type eventListener struct {
	C   chan []*js.Object
	obj *js.Object
	fn  *js.Object
}

func newEventListener(eventName string) *eventListener {
	parts := strings.Split(method, ".")
	var o *js.Object = js.Global
	for _, key := range parts {
		o = o.Get(key)
	}
	l := &eventListener{
		C:   make(chan []*js.Object, 1),
		obj: o,
	}
	l.fn = js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		l.C <- args
	})
	o.Call("addListener", l.fn)
}
func (e *eventListener) Close() error {
	if e.obj == nil {
		return errors.New("already closed")
	}
	e.obj.Call("removeListener", e.fn)
	close(e.C)
	e.obj = nil
	e.fn = nil
	return nil
}
