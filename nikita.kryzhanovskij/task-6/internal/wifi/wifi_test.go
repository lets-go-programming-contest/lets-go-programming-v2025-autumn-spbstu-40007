package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	mac1, _ := net.ParseMAC("00:11:22:33:44:55")
	mac2, _ := net.ParseMAC("66:77:88:99:aa:bb")

	mockInterfaces := []*wifi.Interface{
		{
			Index:        1,
			Name:         "wlan0",
			HardwareAddr: mac1,
		},
		{
			Index:        2,
			Name:         "wlan1",
			HardwareAddr: mac2,
		},
	}

	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		addrs, err := service.GetAddresses()

		assert.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{mac1, mac2}, addrs)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error getting interfaces", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		expectedErr := errors.New("driver error")
		mockHandle.On("Interfaces").Return(nil, expectedErr)

		addrs, err := service.GetAddresses()

		assert.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	mockInterfaces := []*wifi.Interface{
		{Index: 1, Name: "eth0"},
		{Index: 2, Name: "wlan0"},
	}

	t.Run("success", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"eth0", "wlan0"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		service := New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errors.New("failed to list"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		mockHandle.AssertExpectations(t)
	})
}
