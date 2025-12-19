package wifi

import (
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Run("success with multiple interfaces", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		mac1, _ := net.ParseMAC("00:11:22:33:44:55")
		mac2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
		ifaces := []*wifi.Interface{
			{HardwareAddr: mac1},
			{HardwareAddr: mac2},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(addrs) != 2 {
			t.Fatalf("expected 2 addresses, got %d", len(addrs))
		}
		if addrs[0].String() != mac1.String() {
			t.Errorf("address mismatch: %v != %v", addrs[0], mac1)
		}
		mockHandle.AssertExpectations(t)
	})

	t.Run("success with single interface", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		mac, _ := net.ParseMAC("01:23:45:67:89:ab")
		ifaces := []*wifi.Interface{{HardwareAddr: mac}}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(addrs) != 1 {
			t.Fatalf("expected 1 address, got %d", len(addrs))
		}
		mockHandle.AssertExpectations(t)
	})

	t.Run("success with empty interface list", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		ifaces := []*wifi.Interface{}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(addrs) != 0 {
			t.Fatalf("expected empty slice, got %v", addrs)
		}
		mockHandle.AssertExpectations(t)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, assert.AnError)

		service := New(mockHandle)
		addrs, err := service.GetAddresses()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if addrs != nil {
			t.Errorf("expected nil addrs, got %v", addrs)
		}
		if err.Error() != "getting interfaces: mock.AnError general error for testing" {
			t.Errorf("error message mismatch: %v", err)
		}
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("success with multiple names", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		ifaces := []*wifi.Interface{
			{Name: "wlan0"},
			{Name: "eth0"},
			{Name: "wlp2s0"},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := New(mockHandle)
		names, err := service.GetNames()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 3 {
			t.Fatalf("expected 3 names, got %d", len(names))
		}
		if names[0] != "wlan0" || names[2] != "wlp2s0" {
			t.Errorf("names mismatch: %v", names)
		}
		mockHandle.AssertExpectations(t)
	})

	t.Run("success with single name", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		ifaces := []*wifi.Interface{{Name: "eth0"}}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := New(mockHandle)
		names, err := service.GetNames()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 1 || names[0] != "eth0" {
			t.Fatalf("expected [eth0], got %v", names)
		}
		mockHandle.AssertExpectations(t)
	})

	t.Run("empty name list", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		ifaces := []*wifi.Interface{}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := New(mockHandle)
		names, err := service.GetNames()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 0 {
			t.Fatalf("expected empty slice, got %v", names)
		}
		mockHandle.AssertExpectations(t)
	})

	t.Run("error from Interfaces", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, assert.AnError)

		service := New(mockHandle)
		names, err := service.GetNames()

		if err == nil {
			t.Error("expected error, got nil")
		}
		if names != nil {
			t.Errorf("expected nil names, got %v", names)
		}
		mockHandle.AssertExpectations(t)
	})
}

func TestNew(t *testing.T) {
	mockHandle := new(MockWiFiHandle)
	service := New(mockHandle)

	if service.WiFi != mockHandle {
		t.Errorf("expected WiFi to be %v, got %v", mockHandle, service.WiFi)
	}
}
