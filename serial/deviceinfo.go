package serial

import "github.com/gopherjs/gopherjs/js"

type DeviceInfo struct {
	Path        string
	VendorID    int
	ProductID   int
	DisplayName string
}

func (s *DeviceInfo) fromJS(o *js.Object) error {
	s.Path = o.Get("path").String()
	s.VendorID = o.Get("vendorId").Int()
	s.ProductID = o.Get("productId").Int()
	s.DisplayName = o.Get("displayName").String()
	return nil
}
func (s *DeviceInfo) toJS(o *js.Object) error {
	o.Set("path", s.Path)
	o.Set("vendorId", s.VendorID)
	o.Set("productId", s.ProductID)
	o.Set("displayName", s.DisplayName)
	return nil
}
