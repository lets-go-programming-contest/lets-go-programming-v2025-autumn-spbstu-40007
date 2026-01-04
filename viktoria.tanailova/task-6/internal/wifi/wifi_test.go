//nolint:all
package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"

	localwifi "task-6/internal/wifi"
)

var errWiFiNotAvailable = errors.New("wifi not available")

func TestWiFiService_GetAddresses_Variants(t *testing.T) {
	mock := &fakeWiFiHandle{}
	service := localwifi.New(mock)

	t.Run("returns hardware addresses", func(t *testing.T) {
		mac1, _ := net.ParseMAC("aa:bb:cc:dd:ee:01")
		mac2, _ := net.ParseMAC("aa:bb:cc:dd:ee:02")

		mock.fn = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{HardwareAddr: mac1},
				{HardwareAddr: mac2},
			}, nil
		}

		addrs, err := service.GetAddresses()
		assert.NoError(t, err)
		if assert.Len(t, addrs, 2) {
			assert.Equal(t, mac1, addrs[0])
			assert.Equal(t, mac2, addrs[1])
		}
	})

	t.Run("interfaces error", func(t *testing.T) {
		mock.fn = func() ([]*wifi.Interface, error) {
			return nil, errWiFiNotAvailable
		}

		_, err := service.GetAddresses()
		assert.Error(t, err)
	})
}

func TestWiFiService_GetNames_Variants(t *testing.T) {
	mock := &fakeWiFiHandle{}
	service := localwifi.New(mock)

	t.Run("returns interface names", func(t *testing.T) {
		mock.fn = func() ([]*wifi.Interface, error) {
			return []*wifi.Interface{
				{Name: "wlp3s0"},
				{Name: "wlan1"},
			}, nil
		}

		names, err := service.GetNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlp3s0", "wlan1"}, names)
	})

	t.Run("interfaces error", func(t *testing.T) {
		mock.fn = func() ([]*wifi.Interface, error) {
			return nil, errWiFiNotAvailable
		}

		_, err := service.GetNames()
		assert.Error(t, err)
	})
}
