package wifi_test

import (
	"errors"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface
	if v := args.Get(0); v != nil {
		var ok bool
		ifaces, ok = v.([]*wifi.Interface)
		if !ok {
			return nil, errors.New("mock: unexpected type for Interfaces() return[0]")
		}
	}

	//nolint:wrapcheck
	return ifaces, args.Error(1)
}
