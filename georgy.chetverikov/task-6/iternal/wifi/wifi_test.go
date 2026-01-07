package wifi_test

import (
	"errors"
	"net"
	"testing"

	service "github.com/falsefeelings/task-6/iternal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errWiFi = errors.New("wifi error")

func TestNew(t *testing.T) {
	t.Parallel()

	mockHandle := new(MockWiFiHandle)

	svc := service.New(mockHandle)
	assert.NotNil(t, svc)
	assert.Equal(t, mockHandle, svc.WiFi)
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)

		svc := service.New(mockHandle)

		addr1, _ := net.ParseMAC("00:00:00:00:00:01")
		addr2, _ := net.ParseMAC("00:00:00:00:00:02")

		interfaces := []*wifi.Interface{
			{HardwareAddr: addr1},
			{HardwareAddr: addr2},
		}

		mockHandle.On("Interfaces").Return(interfaces, nil).Once()

		addrs, err := svc.GetAddresses()
		require.NoError(t, err)
		assert.Equal(t, []net.HardwareAddr{addr1, addr2}, addrs)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)

		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errWiFi).Once()

		addrs, err := svc.GetAddresses()
		require.Error(t, err)
		assert.Nil(t, addrs)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)

		svc := service.New(mockHandle)

		interfaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		}

		mockHandle.On("Interfaces").Return(interfaces, nil).Once()

		names, err := svc.GetNames()
		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)

		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return(nil, errWiFi).Once()

		names, err := svc.GetNames()
		require.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "getting interfaces")
		mockHandle.AssertExpectations(t)
	})
}
