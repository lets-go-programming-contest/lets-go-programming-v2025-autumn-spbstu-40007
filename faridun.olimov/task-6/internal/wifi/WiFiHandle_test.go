package wifi

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
		return nil, args.Error(1)
	}

	if res, ok := args.Get(0).([]*wifi.Interface); ok {
		return res, args.Error(1)
	}

	return nil, fmt.Errorf("unexpected type: %w", args.Error(1))
}
