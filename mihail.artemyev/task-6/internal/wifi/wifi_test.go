package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	wifiPkg "task-6/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	errWiFi                   = errors.New("wifi error")
	errPermissionDenied       = errors.New("permission denied")
	errMockUnexpectedIfaceTyp = errors.New("mock returned unexpected type")
	errMockNilIfaces          = errors.New("mock returned nil interfaces without error")
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	if args.Get(0) == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock: %w", err)
		}

		return nil, errMockNilIfaces
	}

	got := args.Get(0)

	ifaces, ok := got.([]*wifi.Interface)
	if !ok {
		return nil, errMockUnexpectedIfaceTyp
	}

	if err := args.Error(1); err != nil {
		return ifaces, fmt.Errorf("mock: %w", err)
	}

	return ifaces, nil
}

func TestGetAddresses_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr1, _ := net.ParseMAC("00:11:22:33:44:55")
	hwAddr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

	interfaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: hwAddr1,
		},
		{
			Name:         "wlan1",
			HardwareAddr: hwAddr2,
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := wifiPkg.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Len(t, addrs, 2)
	assert.Equal(t, hwAddr1, addrs[0])
	assert.Equal(t, hwAddr2, addrs[1])
	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_InterfacesError(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return(nil, errWiFi)

	service := wifiPkg.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.Error(t, err)
	assert.Nil(t, addrs)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := wifiPkg.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Equal(t, []net.HardwareAddr{}, addrs)
	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_Single(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr, _ := net.ParseMAC("11:22:33:44:55:66")
	interfaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: hwAddr,
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := wifiPkg.New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Len(t, addrs, 1)
	assert.Equal(t, hwAddr, addrs[0])
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	interfaces := []*wifi.Interface{
		{
			Name: "wlan0",
		},
		{
			Name: "wlan1",
		},
		{
			Name: "eth0",
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := wifiPkg.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_InterfacesError(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return(nil, errPermissionDenied)

	service := wifiPkg.New(mockWiFi)
	names, err := service.GetNames()

	require.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_Empty(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := wifiPkg.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, []string{}, names)
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_Single(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	interfaces := []*wifi.Interface{
		{
			Name: "wlan0",
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := wifiPkg.New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Len(t, names, 1)
	assert.Equal(t, []string{"wlan0"}, names)
	mockWiFi.AssertExpectations(t)
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)
	service := wifiPkg.New(mockWiFi)

	require.NotNil(t, service)
	assert.Equal(t, mockWiFi, service.WiFi)
}

func TestIntegration_MultipleOperations(t *testing.T) {
	t.Parallel()

	mockWiFi := new(MockWiFiHandle)

	hwAddr1, _ := net.ParseMAC("00:11:22:33:44:55")
	hwAddr2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

	interfaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: hwAddr1,
		},
		{
			Name:         "wlan1",
			HardwareAddr: hwAddr2,
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := wifiPkg.New(mockWiFi)

	addrs, err := service.GetAddresses()
	require.NoError(t, err)
	assert.Len(t, addrs, 2)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Len(t, names, 2)

	mockWiFi.AssertExpectations(t)
}
