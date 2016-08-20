package tcpserver

import "github.com/mastercactapus/chrome/internal/util"

// GetSockets is analogous to [chrome.sockets.tcpServer.getSockets](https://developer.chrome.com/apps/sockets_tcpServer#method-getSockets)
func GetSockets() ([]SocketInfo, error) {
	args, err := util.Call("chrome.sockets.tcpServer.getSockets")
	if err != nil {
		return nil, err
	}
	infos := make([]SocketInfo, args[0].Length())
	for i := range infos {
		err = infos[i].fromJS(args[0].Index(i))
		if err != nil {
			return nil, err
		}
	}
	return infos, nil
}
