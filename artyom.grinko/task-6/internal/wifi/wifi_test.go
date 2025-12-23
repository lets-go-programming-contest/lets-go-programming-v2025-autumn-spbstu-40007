package wifi_test

import (
	"fmt"
	"net"
	"testing"

	LeWiFi "task-6/internal/wifi"

	WiFi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var stubAddress = (func() net.HardwareAddr {
	address, _ := net.ParseMAC("00:00:00:00:00:01")

	return address
})()

type wifiStub struct {
	mock.Mock
}

func (wifiStub wifiStub) Interfaces() ([]*WiFi.Interface, error) {
	arguments := wifiStub.Called()
	return arguments.Get(0).([]*WiFi.Interface), arguments.Error(1)
}

func TestGetAddresses(t *testing.T) {
	t.Run("noInterfaces", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{}, fmt.Errorf("some erreur")).Once()
		wifiService := LeWiFi.New(wifi)
		addresses, err := wifiService.GetAddresses()
		require.Empty(t, addresses)
		require.NotNil(t, err)
	})

	t.Run("getSuccess", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{
			{HardwareAddr: stubAddress},
		}, nil).Once()
		wifiService := LeWiFi.New(wifi)
		addresses, err := wifiService.GetAddresses()
		require.NotEmpty(t, addresses)
		require.NotNil(t, err)
	})
}

func TestGetNames(t *testing.T) {
	t.Run("noInterfaces", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{}, fmt.Errorf("some erreur")).Once()
		wifiService := LeWiFi.New(wifi)
		addresses, err := wifiService.GetNames()
		require.Empty(t, addresses)
		require.NotNil(t, err)
	})

	t.Run("getSuccess", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{
			{HardwareAddr: stubAddress},
		}, nil).Once()
		wifiService := LeWiFi.New(wifi)
		names, err := wifiService.GetNames()
		require.NotEmpty(t, names)
		require.NotNil(t, err)
	})
}
