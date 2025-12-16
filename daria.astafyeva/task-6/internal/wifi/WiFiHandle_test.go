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

	var r0 []*wifi.Interface
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]*wifi.Interface)
	}

	var r1 error
	if ret.Error(1) != nil {
		r1 = fmt.Errorf("mock error: %w", ret.Error(1))
	}

	return r0, r1
}