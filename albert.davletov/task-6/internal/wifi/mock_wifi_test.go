package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var interfaces []*wifi.Interface
	if val := args.Get(0); val != nil {
		interfaces = val.([]*wifi.Interface)
	}

	err := args.Error(1)
	if err != nil {
		return interfaces, fmt.Errorf("mock Interfaces error: %w", err)
	}

	return interfaces, nil
}
