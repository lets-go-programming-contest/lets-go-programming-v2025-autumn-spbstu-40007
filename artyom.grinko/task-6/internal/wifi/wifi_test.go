package wifi_test

//nolint:gofumpt
import (
	"fmt"
	"net"
	"testing"

	LeWiFi "task-6/internal/wifi"

	WiFi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errSomeError = fmt.Errorf("some erreur")

var stubAddress = (func() net.HardwareAddr { //nolint:gochecknoglobals
	address, _ := net.ParseMAC("00:00:00:00:00:01")

	return address
})()

type wifiStub struct {
	mock.Mock
}

func (wifiStub wifiStub) Interfaces() ([]*WiFi.Interface, error) {
	arguments := wifiStub.Called()

	return arguments.Get(0).([]*WiFi.Interface), arguments.Error(1) //nolint:forcetypeassert
}

func TestGetAddresses(t *testing.T) {
	t.Run("noInterfaces", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{}, errSomeError)
		wifiService := LeWiFi.New(wifi)
		addresses, err := wifiService.GetAddresses()
		require.Empty(t, addresses)
		require.Error(t, err)
	})

	t.Run("getSuccess", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{
			{HardwareAddr: stubAddress},
		}, nil)
		wifiService := LeWiFi.New(wifi)
		addresses, err := wifiService.GetAddresses()
		require.NotEmpty(t, addresses)
		require.NoError(t, err)
	})
}

func TestGetNames(t *testing.T) {
	t.Run("noInterfaces", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{}, errSomeError)
		wifiService := LeWiFi.New(wifi)
		addresses, err := wifiService.GetNames()
		require.Empty(t, addresses)
		require.Error(t, err)
	})

	t.Run("getSuccess", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{
			{HardwareAddr: stubAddress},
		}, nil)
		wifiService := LeWiFi.New(wifi)
		names, err := wifiService.GetNames()
		require.NotEmpty(t, names)
		require.NoError(t, err)
	})
}
