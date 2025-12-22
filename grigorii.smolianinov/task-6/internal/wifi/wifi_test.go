package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiService_GetAddresses(t *testing.T) {
	mock := &WiFiMock{}
	service := wifi.New(mock)

	t.Run("Success", func(t *testing.T) {
		hw, _ := net.ParseMAC("00:11:22:33:44:55")
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{HardwareAddr: hw},
			}, nil
		}

		addrs, err := service.GetAddresses()
		assert.NoError(t, err)
		assert.Len(t, addrs, 1)
		assert.Equal(t, hw, addrs[0])
	})

	t.Run("Error", func(t *testing.T) {
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi error")
		}
		_, err := service.GetAddresses()
		assert.Error(t, err)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	mock := &WiFiMock{}
	service := wifi.New(mock)

	t.Run("Success Names", func(t *testing.T) {
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlan0"},
				{Name: "wlan1"},
			}, nil
		}

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0", "wlan1"}, names)
	})
}
