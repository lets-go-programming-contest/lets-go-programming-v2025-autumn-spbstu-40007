package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestWiFiService(t *testing.T) {
	mock := &WiFiMock{}
	service := wifi.New(mock)

	t.Run("GetAddresses Success", func(t *testing.T) {
		hw, _ := net.ParseMAC("00:11:22:33:44:55")
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{{HardwareAddr: hw}}, nil
		}
		res, err := service.GetAddresses()
		assert.NoError(t, err)
		assert.Equal(t, hw, res[0])
	})

	t.Run("GetAddresses Error", func(t *testing.T) {
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi error")
		}
		_, err := service.GetAddresses()
		assert.Error(t, err)
	})

	t.Run("GetNames Success", func(t *testing.T) {
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{{Name: "wlan0"}}, nil
		}
		res, err := service.GetNames()
		assert.NoError(t, err)
		assert.Contains(t, res, "wlan0")
	})

	t.Run("GetNames Error", func(t *testing.T) {
		mock.InterfacesFunc = func() ([]*wifi.Interface, error) {
			return nil, errors.New("wifi error")
		}
		_, err := service.GetNames()
		assert.Error(t, err)
	})
}
