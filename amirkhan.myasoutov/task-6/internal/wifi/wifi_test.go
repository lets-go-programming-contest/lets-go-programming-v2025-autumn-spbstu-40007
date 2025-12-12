package wifi_test

import (
	"testing"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	service "github.com/ami0-0/task-6/internal/wifi"
)

type MockWiFiHandle struct { mock.Mock }
func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestGetNames(t *testing.T) {
	m := new(MockWiFiHandle)
	svc := service.New(m)

	m.On("Interfaces").Return([]*wifi.Interface{{Name: "wlan0"}}, nil).Once()

	names, err := svc.GetNames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"wlan0"}, names)
}