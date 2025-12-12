package netinfo_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/mdlayher/wifi"
	"github.com/ami0-0/task-6/internal/netinfo"
)

type MockScanner struct { mock.Mock }
func (m *MockScanner) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestGetInterfaceNames(t *testing.T) {
	m := new(MockScanner)
	svc := netinfo.NewNetworkService(m)

	m.On("Interfaces").Return([]*wifi.Interface{{Name: "wlan0"}}, nil)

	names, err := svc.GetInterfaceNames()
	assert.NoError(t, err)
	assert.Equal(t, []string{"wlan0"}, names)
}