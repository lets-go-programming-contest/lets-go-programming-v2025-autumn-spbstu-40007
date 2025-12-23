package wifi_test

//nolint:gofumpt
import (
	"errors"
	"net"
	"testing"

	LeWiFi "task-6/internal/wifi"

	WiFi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errSomeError = errors.New("some erreur")

var stubAddress = (func() net.HardwareAddr { //nolint:gochecknoglobals
	address, _ := net.ParseMAC("00:00:00:00:00:01")

	return address
})()

type wifiStub struct {
	mock.Mock
}

func (wifiStub wifiStub) Interfaces() ([]*WiFi.Interface, error) { //nolint:govet
	arguments := wifiStub.Called()

	return arguments.Get(0).([]*WiFi.Interface), arguments.Error(1) //nolint:forcetypeassert,wrapcheck
}

func TestGetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("noInterfaces", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{}, errSomeError)
		wifiService := LeWiFi.New(wifi) //nolint:govet
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
		wifiService := LeWiFi.New(wifi) //nolint:govet
		addresses, err := wifiService.GetAddresses()
		require.NotEmpty(t, addresses)
		require.NoError(t, err)
	})
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	t.Run("noInterfaces", func(t *testing.T) {
		t.Parallel()

		wifi := wifiStub{}
		wifi.On("Interfaces").Return([]*WiFi.Interface{}, errSomeError)
		wifiService := LeWiFi.New(wifi) //nolint:govet
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
		wifiService := LeWiFi.New(wifi) //nolint:govet
		names, err := wifiService.GetNames()
		require.NotEmpty(t, names)
		require.NoError(t, err)
	})
}
