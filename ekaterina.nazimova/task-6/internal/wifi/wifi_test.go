package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/UwUshkin/task-6/internal/wifi"
	mdlayher "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errWifiSys = errors.New("system error")

func TestWiFi(t *testing.T) {
	t.Parallel()

	hw, _ := net.ParseMAC("00:11:22:33:44:55")
	testIface := &mdlayher.Interface{
		Index:        1,
		Name:         "wlan0",
		HardwareAddr: hw,
	}

	t.Run("FullCoverage", func(t *testing.T) {
		t.Parallel()

		mockWiFi := new(MockWiFiHandle)
		service := wifi.New(mockWiFi)

		mockWiFi.On("Interfaces").Return([]*mdlayher.Interface{testIface}, nil).Once()
		mockWiFi.On("Interfaces").Return(nil, errWifiSys).Once()
		mockWiFi.On("Interfaces").Return([]*mdlayher.Interface{testIface}, nil).Once()
		mockWiFi.On("Interfaces").Return(nil, errWifiSys).Once()

		addresses, err := service.GetAddresses()
		require.NoError(t, err)
		require.NotNil(t, addresses)

		addresses, err = service.GetAddresses()
		require.Error(t, err)
		require.Nil(t, addresses)

		names, err := service.GetNames()
		require.NoError(t, err)
		require.NotNil(t, names)

		names, err = service.GetNames()
		require.Error(t, err)
		require.Nil(t, names)
	})
}
