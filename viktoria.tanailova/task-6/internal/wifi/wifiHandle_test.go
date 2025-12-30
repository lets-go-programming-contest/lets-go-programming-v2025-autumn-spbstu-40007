package wifi_test

import "github.com/mdlayher/wifi"

type fakeWiFiHandle struct {
	fn func() ([]*wifi.Interface, error)
}

func (f *fakeWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	return f.fn()
}
