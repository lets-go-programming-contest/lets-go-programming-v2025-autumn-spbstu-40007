//nolint:all
package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	localwifi "task-6/internal/wifi"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	mock := &testWiFiHandle{}
	service := localwifi.New(mock)

	t.Run("interfaces ok", func(t *testing.T) {
		mac, _ := net.ParseMAC("00:11:22:33:44:55")

		mock.ifacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{HardwareAddr: mac},
			}, nil
		}

		addrs, err := service.GetAddresses()
		require.NoError(t, err)
		require.Len(t, addrs, 1)
		require.Equal(t, mac, addrs[0])
	})

	t.Run("interfaces error", func(t *testing.T) {
		mock.ifacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi down")
		}

		_, err := service.GetAddresses()
		require.Error(t, err)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	mock := &testWiFiHandle{}
	service := localwifi.New(mock)

	t.Run("interfaces ok", func(t *testing.T) {
		mock.ifacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
			}, nil
		}

		names, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"wlan0", "wlan1"}, names)
	})

	t.Run("interfaces error", func(t *testing.T) {
		mock.ifacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi down")
		}

		_, err := service.GetNames()
		require.Error(t, err)
	})
}
