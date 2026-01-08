package wifi_test

import (
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var interfaces []*wifi.Interface
	if args.Get(0) != nil {
		interfaces = args.Get(0).([]*wifi.Interface)
	}

	return interfaces, args.Error(1)
}
