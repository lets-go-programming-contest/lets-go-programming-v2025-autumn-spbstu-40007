package wifi_test

import (
	"fmt"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

type MockWiFiHandle struct {
	mock.Mock
}

func (_m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	ret := _m.Called()

	var ifaces []*wifi.Interface
	if rf, ok := ret.Get(0).(func() []*wifi.Interface); ok {
		ifaces = rf()
	} else if v, ok := ret.Get(0).([]*wifi.Interface); ok {
		ifaces = v
	}

	var err error
	if rf, ok := ret.Get(1).(func() error); ok {
		err = rf()
	} else {
		if e := ret.Error(1); e != nil {
			err = fmt.Errorf("mock Interfaces: %w", e)
		}
	}

	return ifaces, err
}
