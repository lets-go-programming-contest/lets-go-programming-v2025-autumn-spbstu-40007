package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	wifipkg "nikita.kryzhanovskij/task-6/internal/wifi"
)

var (
	errDriver       = errors.New("driver error")
	errFailedToList = errors.New("failed to list")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	mac1, _ := net.ParseMAC("00:11:22:33:44:55")
	mac2, _ := net.ParseMAC("66:77:88:99:aa:bb")

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		service := wifipkg.New(mockHandle)

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

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		addrs, err := service.GetAddresses()
		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{mac1, mac2}, addrs)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error getting interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		service := wifipkg.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errDriver)

		addrs, err := service.GetAddresses()
		require.Error(t, err)
		assert.Nil(t, addrs)

		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		service := wifipkg.New(mockHandle)

		mockInterfaces := []*wifi.Interface{
			{Index: 1, Name: "eth0"},
			{Index: 2, Name: "wlan0"},
		}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		names, err := service.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"eth0", "wlan0"}, names)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		service := wifipkg.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errFailedToList)

		names, err := service.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)

		mockHandle.AssertExpectations(t)
	})
}
