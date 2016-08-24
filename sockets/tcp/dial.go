package tcp

import (
	"fmt"
	"net"
	"strconv"

	"github.com/mastercactapus/chrome/internal/util"
)

func Dial(network, address string) (net.Conn, error) {
	if network != "tcp" {
		return nil, fmt.Errorf("unsupported network type: %s", network)
	}

	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}
	args, err := util.Call("chrome.sockets.tcp.create")
	if err != nil {
		return nil, err
	}
	socketID := args[0].Get("socketId").Int()

	conn := NewConnection(socketID)

	args, err = util.Call("chrome.sockets.tcp.connect", socketID, host, port)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if args[0].Int() < 0 {
		conn.Close()
		return nil, fmt.Errorf("error code: %d", args[0].Int())
	}
	return conn, nil
}
