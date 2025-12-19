package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockWiFiHandle реализует интерфейс WiFiHandle для тестирования
type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestNew(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	service := New(mockWiFi)

	assert.NotNil(t, service)
	assert.Equal(t, mockWiFi, service.WiFi)
}

func TestGetAddresses_Success(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)

	// Создаем тестовые интерфейсы
	interfaces := []*wifi.Interface{
		{Interface: &net.Interface{Name: "wlan0", HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}}},
		{Interface: &net.Interface{Name: "wlan1", HardwareAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}}},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	require.Len(t, addrs, 2)
	assert.Equal(t, net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}, addrs[0])
	assert.Equal(t, net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}, addrs[1])
	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_Error(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	expectedErr := errors.New("interface error")

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, expectedErr)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	assert.Nil(t, addrs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockWiFi.AssertExpectations(t)
}

func TestGetAddresses_NoInterfaces(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	require.NoError(t, err)
	assert.Empty(t, addrs)
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_Success(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)

	interfaces := []*wifi.Interface{
		{Interface: &net.Interface{Name: "wlan0", HardwareAddr: net.HardwareAddr{0x00, 0x11, 0x22, 0x33, 0x44, 0x55}}},
		{Interface: &net.Interface{Name: "wlan1", HardwareAddr: net.HardwareAddr{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}}},
		{Interface: &net.Interface{Name: "eth0", HardwareAddr: net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}}},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	require.Len(t, names, 3)
	assert.Equal(t, "wlan0", names[0])
	assert.Equal(t, "wlan1", names[1])
	assert.Equal(t, "eth0", names[2])
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_Error(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	expectedErr := errors.New("interface error")

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, expectedErr)

	service := New(mockWiFi)
	names, err := service.GetNames()

	assert.Nil(t, names)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockWiFi.AssertExpectations(t)
}

func TestGetNames_Empty(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)

	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := New(mockWiFi)
	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Empty(t, names)
	mockWiFi.AssertExpectations(t)
}
