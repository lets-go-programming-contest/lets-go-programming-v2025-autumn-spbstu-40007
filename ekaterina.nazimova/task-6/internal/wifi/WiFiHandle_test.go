package wifi_test

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errTypeAssertion = errors.New("type assertion failed")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	err := args.Error(1)

	if args.Get(0) == nil {
		if err != nil {
			return nil, fmt.Errorf("mock error: %w", err)
		}

		return nil, nil
	}

	ifaces, ok := args.Get(0).([]*wifi.Interface)
	if !ok {
		return nil, errTypeAssertion
	}

	if err != nil {
		return ifaces, fmt.Errorf("external error: %w", err)
	}

	return ifaces, nil
}
