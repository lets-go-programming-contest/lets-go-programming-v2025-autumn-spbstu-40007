//nolint:all
package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
	localWifi "grigorii.smolianinov/task-6/internal/wifi"
)

func TestWiFiService_Coverage(t *testing.T) {
	mock := &WiFiMock{}
	service := localWifi.New(mock)

	t.Run("GetAddresses", func(t *testing.T) {
		hw, _ := net.ParseMAC("00:11:22:33:44:55")
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{{HardwareAddr: hw}}, nil
		}
		res, err := service.GetAddresses()
		require.NoError(t, err)
		require.Len(t, res, 1)

		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi fail")
		}
		_, err = service.GetAddresses()
		require.Error(t, err)
	})

	t.Run("GetNames", func(t *testing.T) {
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{{Name: "wlan0"}}, nil
		}
		res, err := service.GetNames()
		require.NoError(t, err)
		require.Equal(t, "wlan0", res[0])

		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi fail")
		}
		_, err = service.GetNames()
		require.Error(t, err)
	})
}
