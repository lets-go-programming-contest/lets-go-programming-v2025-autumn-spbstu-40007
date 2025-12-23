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

	val := args.Get(0)
	if val == nil {
		return nil, fmt.Errorf("mock error: %w", args.Error(1))
	}

	ifaces, ok := val.([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("unexpected type: %T", val)
	}

	return ifaces, nil
}
