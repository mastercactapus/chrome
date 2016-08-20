package serial

import "github.com/mastercactapus/chrome/internal/util"

// GetDevices is analogous to [chrome.serial.getDevices](https://developer.chrome.com/apps/serial#method-getDevices)
func GetDevices() ([]DeviceInfo, error) {
	args, err := util.Call("chrome.serial.getDevices")
	if err != nil {
		return nil, err
	}
	infos := make([]DeviceInfo, args[0].Length())
	for i := range infos {
		err = infos[i].fromJS(args[0].Index(i))
		if err != nil {
			return nil, err
		}
	}
	return infos, nil
}

// GetConnections is analogous to [chrome.serial.getConnections](https://developer.chrome.com/apps/serial#method-getConnections)
func GetConnections() ([]ConnectionInfo, error) {
	args, err := util.Call("chrome.serial.getConnections")
	if err != nil {
		return nil, err
	}
	conns := make([]ConnectionInfo, args[0].Length())
	for i := range conns {
		err = conns[i].fromJS(args[0].Index(i))
		if err != nil {
			return nil, err
		}
	}
	return conns, nil
}
