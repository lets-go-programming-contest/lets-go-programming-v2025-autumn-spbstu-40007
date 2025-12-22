package wifi_test

import (
	"github.com/mdlayher/wifi"
)

// WiFiMock имитирует WiFiHandle для тестов
type WiFiMock struct {
	InterfacesFunc func() ([]*wifi.Interface, error)
}

func (m *WiFiMock) Interfaces() ([]*wifi.Interface, error) {
	return m.InterfacesFunc()
}
