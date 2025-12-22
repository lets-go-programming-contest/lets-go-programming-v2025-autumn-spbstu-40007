package wifi

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWiFiHandle для тестирования
type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

// Test GetAddresses - успешное получение адресов
func TestGetAddresses_Success(t *testing.T) {
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

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	assert.NoError(t, err)
	assert.Len(t, addrs, 2)
	assert.Equal(t, hwAddr1, addrs[0])
	assert.Equal(t, hwAddr2, addrs[1])
	mockWiFi.AssertExpectations(t)
}

// Test GetAddresses - ошибка при получении интерфейсов
func TestGetAddresses_InterfacesError(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return(nil, errors.New("wifi error"))

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	assert.Error(t, err)
	assert.Nil(t, addrs)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockWiFi.AssertExpectations(t)
}

// Test GetAddresses - пустой список интерфейсов
func TestGetAddresses_Empty(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	assert.NoError(t, err)
	assert.Equal(t, []net.HardwareAddr{}, addrs)
	mockWiFi.AssertExpectations(t)
}

// Test GetAddresses - один интерфейс
func TestGetAddresses_Single(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)

	hwAddr, _ := net.ParseMAC("11:22:33:44:55:66")
	interfaces := []*wifi.Interface{
		{
			Name:         "wlan0",
			HardwareAddr: hwAddr,
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	addrs, err := service.GetAddresses()

	assert.NoError(t, err)
	assert.Len(t, addrs, 1)
	assert.Equal(t, hwAddr, addrs[0])
	mockWiFi.AssertExpectations(t)
}

// Test GetNames - успешное получение имён
func TestGetNames_Success(t *testing.T) {
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

	service := New(mockWiFi)
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Len(t, names, 3)
	assert.Equal(t, []string{"wlan0", "wlan1", "eth0"}, names)
	mockWiFi.AssertExpectations(t)
}

// Test GetNames - ошибка при получении интерфейсов
func TestGetNames_InterfacesError(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return(nil, fmt.Errorf("permission denied"))

	service := New(mockWiFi)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockWiFi.AssertExpectations(t)
}

// Test GetNames - пустой список
func TestGetNames_Empty(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	mockWiFi.On("Interfaces").Return([]*wifi.Interface{}, nil)

	service := New(mockWiFi)
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{}, names)
	mockWiFi.AssertExpectations(t)
}

// Test GetNames - один интерфейс
func TestGetNames_Single(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)

	interfaces := []*wifi.Interface{
		{
			Name: "wlan0",
		},
	}

	mockWiFi.On("Interfaces").Return(interfaces, nil)

	service := New(mockWiFi)
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Len(t, names, 1)
	assert.Equal(t, []string{"wlan0"}, names)
	mockWiFi.AssertExpectations(t)
}

// Test New - создание сервиса
func TestNew(t *testing.T) {
	mockWiFi := new(MockWiFiHandle)
	service := New(mockWiFi)

	assert.NotNil(t, service)
	assert.Equal(t, mockWiFi, service.WiFi)
}

// Test интеграция - несколько вызовов
func TestIntegration_MultipleOperations(t *testing.T) {
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

	service := New(mockWiFi)

	// GetAddresses
	addrs, err := service.GetAddresses()
	assert.NoError(t, err)
	assert.Len(t, addrs, 2)

	// GetNames
	names, err := service.GetNames()
	assert.NoError(t, err)
	assert.Len(t, names, 2)

	mockWiFi.AssertExpectations(t)
}
