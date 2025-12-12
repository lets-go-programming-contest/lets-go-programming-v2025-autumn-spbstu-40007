package wifi_test

import (
	"errors"
	"testing"

	service "github.com/ami0-0/task-6/internal/wifi"
	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errDriver = errors.New("driver error")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	if val, ok := args.Get(0).([]*wifi.Interface); ok {
		return val, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		m := new(MockWiFiHandle)
		svc := service.New(m)

		fakeIfaces := []*wifi.Interface{{Name: "wlan0"}, {Name: "eth0"}}
		m.On("Interfaces").Return(fakeIfaces, nil).Once()

		names, err := svc.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "eth0"}, names)
		m.AssertExpectations(t)
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()
		m := new(MockWiFiHandle)
		svc := service.New(m)

		fakeIfaces := []*wifi.Interface{}
		m.On("Interfaces").Return(fakeIfaces, nil).Once()

		names, err := svc.GetNames()
		require.NoError(t, err)
		assert.Empty(t, names)
	})

	t.Run("interface error", func(t *testing.T) {
		t.Parallel()
		m := new(MockWiFiHandle)
		svc := service.New(m)

		m.On("Interfaces").Return(nil, errDriver).Once()

		_, err := svc.GetNames()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
		m.AssertExpectations(t)
	})
}