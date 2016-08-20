package serial

import (
	"fmt"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

type ParityBit string

const (
	ParityBitNo   ParityBit = "no"
	ParityBitOdd  ParityBit = "odd"
	ParityBitEven ParityBit = "even"
)

type ConnectionOptions struct {
	Persistent     bool
	Name           string
	BufferSize     int
	BitRate        int
	DataBits       int
	ParityBit      ParityBit
	StopBits       int
	CTSFlowControl bool
	ReceiveTimeout time.Duration
	SendTimeout    time.Duration
}

func (s *ConnectionOptions) fromJS(o *js.Object) error {
	s.Persistent = o.Get("persistent").Bool()
	s.Name = o.Get("name").String()
	s.BufferSize = o.Get("bufferSize").Int()
	s.BitRate = o.Get("bitrate").Int()
	switch o.Get("dataBits").String() {
	case "seven":
		s.DataBits = 7
	case "eight", "":
		s.DataBits = 8
	default:
		return fmt.Errorf("invalid dataBits value: %s", o.Get("dataBits").String())
	}
	switch o.Get("parityBit").String() {
	case "odd":
		s.ParityBit = ParityBitOdd
	case "even":
		s.ParityBit = ParityBitEven
	case "no", "":
		s.ParityBit = ParityBitNo
	default:
		return fmt.Errorf("invalid parityBit value: %s", o.Get("parityBit").String())
	}
	switch o.Get("stopBits").String() {
	case "two":
		s.StopBits = 2
	case "one", "":
		s.StopBits = 1
	default:
		return fmt.Errorf("invalid stopBits value: %s", o.Get("stopBits").String())
	}
	s.CTSFlowControl = o.Get("ctsFlowControl").Bool()
	s.ReceiveTimeout = time.Duration(float64(time.Millisecond) * o.Get("receiveTimeout").Float())
	s.SendTimeout = time.Duration(float64(time.Millisecond) * o.Get("sendTimeout").Float())
	return nil
}

func (s *ConnectionOptions) toJS(o *js.Object) error {
	o.Set("persistent", s.Persistent)
	o.Set("name", s.Name)
	o.Set("bufferSize", s.BufferSize)
	o.Set("bitrate", s.BitRate)
	switch s.DataBits {
	case 7:
		o.Set("dataBits", "seven")
	case 8, 0:
		o.Set("dataBits", "eight")
	default:
		return fmt.Errorf("invalid DataBits value: %d", s.DataBits)
	}
	switch s.ParityBit {
	case ParityBitEven:
		o.Set("parityBit", "even")
	case ParityBitOdd:
		o.Set("parityBit", "odd")
	case ParityBitNo, ParityBit(""):
		o.Set("parityBit", "no")
	default:
		return fmt.Errorf("invalid ParityBit value: %s", s.ParityBit)
	}
	switch s.StopBits {
	case 2:
		o.Set("stopBits", "two")
	case 1, 0:
		o.Set("stopBits", "one")
	default:
		return fmt.Errorf("invalid StopBits value: %d", s.StopBits)
	}
	o.Set("ctsFlowControl", s.CTSFlowControl)
	o.Set("receiveTimeout", s.ReceiveTimeout.Seconds()*1000)
	o.Set("sendTimeout", s.SendTimeout.Seconds()*1000)
	return nil
}
