package storage

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mastercactapus/chrome/internal/util"
)

var (
	Sync    = &areaSync{area: "sync"}
	Local   = &areaLocal{area: "local"}
	Managed = &areaManaged{area: "managed"}
)

type areaSync struct{ area }
type areaLocal struct{ area }
type areaManaged struct{ area }

type area string

func (s areaLocal) QuotaBytes() int {
	return s.getIntProp("QUOTA_BYTES")
}

func (s areaSync) QuotaBytes() int {
	return s.getIntProp("QUOTA_BYTES")
}
func (s areaSync) QuotaBytesPerItem() int {
	return s.getIntProp("QUOTA_BYTES_PER_ITEM")
}
func (s areaSync) MaxItems() int {
	return s.getIntProp("MAX_ITEMS")
}
func (s areaSync) MaxWriteOperationsPerHour() int {
	return s.getIntProp("MAX_WRITE_OPERATIONS_PER_HOUR")
}
func (s areaSync) MaxWriteOperationsPerMinute() int {
	return s.getIntProp("MAX_WRITE_OPERATIONS_PER_MINUTE")
}
func (s areaSync) MaxSustainedWriteOperationsPerMinute() int {
	return s.getIntProp("MAX_SUSTAINED_WRITE_OPERATIONS_PER_MINUTE")
}

func (s area) getIntProp(name string) int {
	return js.Global.Get("chrome").Get("storage").Get(string(s)).Get(name).Int()
}

func (s area) Get(key string) (string, error) {
	return s.GetWithDefault(key, "")
}
func (s area) GetWithDefault(key, def string) (string, error) {
	res, err := util.Call("chrome.storage."+string(s)+".get", key)
	if err != nil {
		return "", err
	}

	val := res[0].Get(key)
	if val == js.Undefined {
		return def, nil
	}

	return val.String(), nil
}
func (s area) GetMany(keys []string) ([]string, error) {
	res, err := util.Call("chrome.storage."+string(s)+".get", keys)
	if err != nil {
		return nil, err
	}

	vals := make([]string, len(keys))
	for i, key := range keys {
		v := res[0].Get(key)
		if v == js.Undefined {
			vals[i] = ""
		} else {
			vals[i] = v.String()
		}
	}

	return vals, nil
}
func (s area) GetManyWithDefaults(keys map[string]string) (map[string]string, error) {
	res, err := util.Call("chrome.storage."+string(s)+".get", keys)
	if err != nil {
		return nil, err
	}

	vals := make(map[string]string, len(keys))
	for key, def := range keys {
		v := res[0].Get(key)
		if v == js.Undefined {
			vals[key] = def
		} else {
			vals[key] = v.String()
		}
	}

	return vals, nil
}

func (s area) Set(key, val string) error {
	return s.SetMany(map[string]string{key: val})
}
func (s area) SetMany(data map[string]string) error {
	o := new(js.Object)
	for key, val := range data {
		o.Set(key, val)
	}
	_, err := util.Call("chrome.storage."+string(s)+".set", o)
	return err
}

func (s area) Remove(key string) error {
	_, err := util.Call("chrome.storage."+string(s)+".remove", key)
	return err
}
func (s area) RemoveMany(keys []string) error {
	_, err := util.Call("chrome.storage."+string(s)+".remove", keys)
	return err
}

func (s area) Clear() error {
	_, err := util.Call("chrome.storage." + string(s) + ".clear")
	return err
}
func (s area) GetBytesInUse(key string) (int, error) {
	res, err := util.Call("chrome.storage."+string(s)+".getBytesInUse", key)
	if err != nil {
		return -1, err
	}
	return res[0].Int(), nil
}
func (s area) GetBytesInUseMany(keys []string) (int, error) {
	res, err := util.Call("chrome.storage."+string(s)+".getBytesInUse", keys)
	if err != nil {
		return -1, err
	}
	return res[0].Int(), nil
}
