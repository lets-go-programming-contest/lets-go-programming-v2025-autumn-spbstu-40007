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

	if args.Get(0) == nil {
		err := args.Error(1)
		return nil, fmt.Errorf("interface error: %w", err)
	}

	res, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, fmt.Errorf("type assertion failed: %w", args.Error(1))
	}

	return res, args.Error(1)
}
