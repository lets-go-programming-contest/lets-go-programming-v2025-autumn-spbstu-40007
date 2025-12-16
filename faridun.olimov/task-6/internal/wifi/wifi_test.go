package wifi

import (
	"errors"
	"net"
	"testing"

	mdwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

var (
	errSystem = errors.New("system error")
	errFailed = errors.New("failed")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockHandle := new(MockWiFiHandle)
		mac1, _ := net.ParseMAC("00:00:5e:00:53:01")

		mockInterfaces := []*mdwifi.Interface{
			{HardwareAddr: mac1},
		}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)
		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		assert.NoError(t, err)
		assert.Len(t, addrs, 1)
	})

	t.Run("error from interfaces", func(t *testing.T) {
		t.Parallel()
		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errSystem)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		assert.Error(t, err)
		assert.Nil(t, addrs)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockHandle := new(MockWiFiHandle)
		mockInterfaces := []*mdwifi.Interface{
			{Name: "wlan0"},
		}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)
		service := New(mockHandle)
		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()
		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errFailed)

		service := New(mockHandle)
		_, err := service.GetNames()

		assert.Error(t, err)
	})
}
