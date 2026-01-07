package wifi_test

import (
	"errors"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var errBadInterfacesType = errors.New("mock: unexpected type for Interfaces() return[0]")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	var ifaces []*wifi.Interface

	v := args.Get(0)
	if v != nil {
		var ok bool

		ifaces, ok = v.([]*wifi.Interface)
		if !ok {
			return nil, errBadInterfacesType
		}
	}

	//nolint:wrapcheck
	return ifaces, args.Error(1)
}
