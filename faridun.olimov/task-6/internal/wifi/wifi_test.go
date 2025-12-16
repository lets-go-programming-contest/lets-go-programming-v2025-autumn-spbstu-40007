package wifi_test

import (
	"errors"
	"net"
	"testing"

	"task-6/internal/wifi"

	mdwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errSystem = errors.New("system error")
	errFailed = errors.New("failed")
)

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mac, _ := net.ParseMAC("00:00:5e:00:53:01")
		mockInterfaces := []*mdwifi.Interface{{HardwareAddr: mac}}

		mockHandle.On("Interfaces").Return(mockInterfaces, nil)

		service := wifi.New(mockHandle)
		addrs, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Len(t, addrs, 1)
	})

	t.Run("error_interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errSystem)

		service := wifi.New(mockHandle)
		_, err := service.GetAddresses()

		require.Error(t, err)
	})
}
