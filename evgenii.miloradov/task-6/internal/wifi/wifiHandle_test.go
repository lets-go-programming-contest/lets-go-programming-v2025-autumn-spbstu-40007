package wifi_test

import "github.com/mdlayher/wifi"

type testWiFiHandle struct {
	ifacesFunc func() ([]*wifi.Interface, error)
}

func (m *testWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	return m.ifacesFunc()
}
