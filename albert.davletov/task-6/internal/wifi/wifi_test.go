package wifi_test

import (
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mywifi "task-6/internal/wifi"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)

		service := mywifi.New(mockWiFi)

		hwAddr1, _ := net.ParseMAC("00:11:22:33:44:55")
		hwAddr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

		testInterfaces := []*wifi.Interface{
			{
				Name:         "wlan0",
				HardwareAddr: hwAddr1,
			},
			{
				Name:         "wlan1",
				HardwareAddr: hwAddr2,
			},
		}

		mockWiFi.On("Interfaces").Return(testInterfaces, nil)

		addresses, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Len(t, addresses, 2)
		assert.Equal(t, []net.HardwareAddr{hwAddr1, hwAddr2}, addresses)

		mockWiFi.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := mywifi.New(mockWiFi)

		expectedErr := assert.AnError
		mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, expectedErr)

		addresses, err := service.GetAddresses()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
		assert.Contains(t, err.Error(), expectedErr.Error())
		assert.Nil(t, addresses)

		mockWiFi.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := mywifi.New(mockWiFi)

		hwAddr, _ := net.ParseMAC("00:11:22:33:44:55")
		testInterfaces := []*wifi.Interface{
			{
				Name:         "wlan0",
				HardwareAddr: hwAddr,
			},
			{
				Name:         "wlp3s0",
				HardwareAddr: hwAddr,
			},
			{
				Name:         "eth0",
				HardwareAddr: hwAddr,
			},
		}

		mockWiFi.On("Interfaces").Return(testInterfaces, nil)

		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Len(t, names, 3)
		assert.Equal(t, []string{"wlan0", "wlp3s0", "eth0"}, names)

		mockWiFi.AssertExpectations(t)
	})

	t.Run("failure", func(t *testing.T) {
		t.Parallel()
		mockWiFi := new(MockWiFiHandle)
		service := mywifi.New(mockWiFi)

		expectedErr := assert.AnError
		mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, expectedErr)

		names, err := service.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "getting interfaces")
		assert.Contains(t, err.Error(), expectedErr.Error())
		assert.Nil(t, names)

		mockWiFi.AssertExpectations(t)
	})
}
