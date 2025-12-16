package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)

		mac1, _ := net.ParseMAC("00:00:5e:00:53:01")
		mac2, _ := net.ParseMAC("00:00:5e:00:53:02")

		mockInterfaces := []*wifi.Interface{
			{HardwareAddr: mac1},
			{HardwareAddr: mac2},
		}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		assert.NoError(t, err)
		assert.Len(t, addrs, 2)
		assert.Equal(t, mac1, addrs[0])
		assert.Equal(t, mac2, addrs[1])

		mockHandle.AssertExpectations(t)
	})

	t.Run("error from interfaces", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		expectedErr := errors.New("system error")

		mockHandle.On("Interfaces").Return(nil, expectedErr)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		assert.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")

		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)

		mockInterfaces := []*wifi.Interface{
			{Name: "eth0"},
			{Name: "wlan0"},
		}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		service := New(mockHandle)
		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"eth0", "wlan0"}, names)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("failed"))

		service := New(mockHandle)
		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)

		mockHandle.AssertExpectations(t)
	})
}
