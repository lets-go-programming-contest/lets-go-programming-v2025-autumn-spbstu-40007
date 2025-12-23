package wifi

import (
	"errors"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestGetAddresses_SingleInterface(t *testing.T) {
	mockHandle := new(MockWiFiHandle)

	mac, _ := net.ParseMAC("11:22:33:44:55:66")
	interfaces := []*wifi.Interface{
		{HardwareAddr: mac, Name: "wlan0"},
	}

	mockHandle.On("Interfaces").Return(interfaces, nil)

	service := New(mockHandle)
	addrs, err := service.GetAddresses()

	assert.NoError(t, err)
	assert.Len(t, addrs, 1)
	assert.Equal(t, mac, addrs[0])
	mockHandle.AssertExpectations(t)
}

func TestGetAddresses_Error(t *testing.T) {
	mockHandle := new(MockWiFiHandle)
	mockHandle.On("Interfaces").Return(nil, errors.New("wifi error"))

	service := New(mockHandle)
	addrs, err := service.GetAddresses()

	assert.Error(t, err)
	assert.Nil(t, addrs)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockHandle.AssertExpectations(t)
}

func TestGetNames_SingleInterface(t *testing.T) {
	mockHandle := new(MockWiFiHandle)

	interfaces := []*wifi.Interface{
		{Name: "eth0"},
	}

	mockHandle.On("Interfaces").Return(interfaces, nil)

	service := New(mockHandle)
	names, err := service.GetNames()

	assert.NoError(t, err)
	assert.Equal(t, []string{"eth0"}, names)
	mockHandle.AssertExpectations(t)
}

func TestGetNames_Error(t *testing.T) {
	mockHandle := new(MockWiFiHandle)
	mockHandle.On("Interfaces").Return(nil, errors.New("wifi error"))

	service := New(mockHandle)
	names, err := service.GetNames()

	assert.Error(t, err)
	assert.Nil(t, names)
	assert.Contains(t, err.Error(), "getting interfaces")
	mockHandle.AssertExpectations(t)
}

func TestWiFiHandleInterface(t *testing.T) {
	var handle WiFiHandle = new(MockWiFiHandle)
	assert.NotNil(t, handle)
}

func TestWiFiServiceStruct(t *testing.T) {
	mockHandle := new(MockWiFiHandle)
	service := WiFiService{WiFi: mockHandle}

	assert.NotNil(t, service.WiFi)
	assert.Equal(t, mockHandle, service.WiFi)
}
