package util

import (
	"errors"
	"strings"

	"github.com/gopherjs/gopherjs/js"
)

func Call(method string, args ...interface{}) ([]*js.Object, error) {
	parts := strings.Split(method, ".")
	var o *js.Object = js.Global
	for _, key := range parts[:len(parts)-1] {
		o = o.Get(key)
	}
	ch := make(chan []*js.Object, 1)
	var err error
	args = append(args, js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
		lastErr := js.Global.Get("chrome").Get("runtime").Get("lastError")
		if lastErr != js.Undefined {
			err = errors.New(lastErr.Get("message").String())
		}
		ch <- args
		return nil
	}))

	o.Call(parts[len(parts)-1], args...)

	return <-ch, err
}
