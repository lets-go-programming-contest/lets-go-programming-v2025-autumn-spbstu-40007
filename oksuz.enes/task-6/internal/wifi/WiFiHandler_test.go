package wifi_test

import (
	"errors"
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errUnexpectedType = errors.New("unexpected type")

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
		return nil, fmt.Errorf("%w: %T", errUnexpectedType, val)
	}

	return ifaces, nil
}
