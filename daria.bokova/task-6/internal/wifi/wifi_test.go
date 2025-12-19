package wifi

import (
	"net"
	"testing"

	wifiPkg "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
)

func TestAddressRetrieval(t *testing.T) {
	t.Run("ShouldReturnMACAddresses", func(t *testing.T) {
		mock := new(WiFiMock)
		macAddr, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
		expectedInterfaces := []*wifiPkg.Interface{
			{HardwareAddr: macAddr},
		}

		mock.On("GetInterfaces").Return(expectedInterfaces, nil)

		serviceInstance := New(mock)
		addresses, err := serviceInstance.GetAddresses()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(addresses) != 1 {
			t.Fatalf("Expected 1 address, got %d", len(addresses))
		}

		if addresses[0].String() != macAddr.String() {
			t.Errorf("Expected address %v, got %v", macAddr, addresses[0])
		}

		mock.AssertExpectations(t)
	})

	t.Run("ShouldHandleInterfaceErrors", func(t *testing.T) {
		mock := new(WiFiMock)
		mock.On("GetInterfaces").Return(nil, assert.AnError)

		serviceInstance := New(mock)
		addresses, err := serviceInstance.GetAddresses()

		if err == nil {
			t.Error("Expected error but got none")
		}

		if addresses != nil {
			t.Error("Expected nil addresses on error")
		}

		mock.AssertExpectations(t)
	})
}

func TestNameRetrieval(t *testing.T) {
	t.Run("ShouldReturnInterfaceNames", func(t *testing.T) {
		mock := new(WiFiMock)
		expectedNames := []*wifiPkg.Interface{
			{Name: "eth0"},
			{Name: "wlan0"},
		}

		mock.On("GetInterfaces").Return(expectedNames, nil)

		serviceInstance := New(mock)
		names, err := serviceInstance.GetNames()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(names) != 2 {
			t.Fatalf("Expected 2 names, got %d", len(names))
		}

		if names[0] != "eth0" || names[1] != "wlan0" {
			t.Errorf("Expected names [eth0, wlan0], got %v", names)
		}

		mock.AssertExpectations(t)
	})

	t.Run("ShouldHandleNameRetrievalErrors", func(t *testing.T) {
		mock := new(WiFiMock)
		mock.On("GetInterfaces").Return(nil, assert.AnError)

		serviceInstance := New(mock)
		names, err := serviceInstance.GetNames()

		if err == nil {
			t.Error("Expected error but got none")
		}

		if names != nil {
			t.Error("Expected nil names on error")
		}

		mock.AssertExpectations(t)
	})
}

func TestServiceCreation(t *testing.T) {
	mock := new(WiFiMock)
	service := New(mock)

	if service.WiFi != mock {
		t.Errorf("Service should contain provided WiFi handle")
	}
}
