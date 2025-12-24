package wifi_test

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"

	wifisvc "github.com/ksuah/task-6/internal/wifi"
)

var errOperationFailed = errors.New("operation failed")

type mockWiFiHandler struct {
	interfaces []*wifi.Interface
	err        error
	callCount  int
}

func (m *mockWiFiHandler) Interfaces() ([]*wifi.Interface, error) {
	m.callCount++

	return m.interfaces, m.err
}

func TestWiFiService_GetAddresses_OK(t *testing.T) {
	t.Parallel()

	mockHandler := &mockWiFiHandler{
		interfaces: []*wifi.Interface{
			{
				Name:         "wlan0",
				HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			},
			{
				Name:         "wlan1",
				HardwareAddr: net.HardwareAddr{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
			},
		},
	}

	service := wifisvc.New(mockHandler)

	addresses, err := service.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, []net.HardwareAddr{
		{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
		{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
	}, addresses)
	require.Equal(t, 1, mockHandler.callCount)
}

func TestWiFiService_GetAddresses_InterfacesError(t *testing.T) {
	t.Parallel()

	mockHandler := &mockWiFiHandler{err: errOperationFailed}
	service := wifisvc.New(mockHandler)

	addresses, err := service.GetAddresses()
	require.Nil(t, addresses)
	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces:")
	require.ErrorIs(t, err, errOperationFailed)
	require.Equal(t, 1, mockHandler.callCount)
}

func TestWiFiService_GetNames_OK(t *testing.T) {
	t.Parallel()

	mockHandler := &mockWiFiHandler{
		interfaces: []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "wlan1"},
		},
	}

	service := wifisvc.New(mockHandler)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"wlan0", "wlan1"}, names)
	require.Equal(t, 1, mockHandler.callCount)
}

func TestWiFiService_GetNames_InterfacesError(t *testing.T) {
	t.Parallel()

	mockHandler := &mockWiFiHandler{err: errOperationFailed}
	service := wifisvc.New(mockHandler)

	names, err := service.GetNames()
	require.Nil(t, names)
	require.Error(t, err)
	require.ErrorContains(t, err, "getting interfaces:")
	require.ErrorIs(t, err, errOperationFailed)
	require.Equal(t, 1, mockHandler.callCount)
}
