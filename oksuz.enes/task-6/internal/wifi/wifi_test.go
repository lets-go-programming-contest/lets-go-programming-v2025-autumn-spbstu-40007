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
		service := New(mockHandle)

		addr1, _ := net.ParseMAC("00:00:5e:00:53:01")
		ifaces := []*wifi.Interface{
			{HardwareAddr: addr1},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		addrs, err := service.GetAddresses()
		assert.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, addr1, addrs[0])
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("wifi error"))

		addrs, err := service.GetAddresses()
		assert.Error(t, err)
		assert.Nil(t, addrs)
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		ifaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("wifi error"))

		names, err := service.GetNames()
		assert.Error(t, err)
		assert.Nil(t, names)
		mockHandle.AssertExpectations(t)
	})
}
