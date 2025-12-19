package wifi_test

import (
	"errors"
	"fmt"
	"net"
	"testing"

	mdlayherwifi "github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"

	wifipkg "task-6/internal/wifi"
)

var (
	errTest          = errors.New("test error")
	ErrTypeAssertion = errors.New("type assertion failed")
	ErrMockReturned  = errors.New("mock returned error")
	ErrMock          = errors.New("mock error")
)

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*mdlayherwifi.Interface, error) {
	args := m.Called()

	interfacesRaw := args.Get(0)
	if interfacesRaw == nil {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrMockReturned, err)
		}

		return nil, nil
	}

	interfaces, ok := interfacesRaw.([]*mdlayherwifi.Interface)
	if !ok {
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("%w with error: %w", ErrTypeAssertion, err)
		}

		return nil, fmt.Errorf("%w: expected []*wifi.Interface, got %T",
			ErrTypeAssertion, interfacesRaw)
	}

	if err := args.Error(1); err != nil {
		return interfaces, fmt.Errorf("%w: %w", ErrMock, err)
	}

	return interfaces, nil
}

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success with multiple interfaces", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mac1, _ := net.ParseMAC("00:11:22:33:44:55")
		mac2, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")
		ifaces := []*mdlayherwifi.Interface{
			{HardwareAddr: mac1},
			{HardwareAddr: mac2},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := wifipkg.New(mockHandle)
		addrs, err := service.GetAddresses()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(addrs) != 2 {
			t.Fatalf("expected 2 addresses, got %d", len(addrs))
		}

		mockHandle.AssertExpectations(t)
	})

	t.Run("success with empty interface list", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		ifaces := []*mdlayherwifi.Interface{}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := wifipkg.New(mockHandle)
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
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errTest)

		service := wifipkg.New(mockHandle)
		addrs, err := service.GetAddresses()

		if err == nil {
			t.Error("expected error, got nil")
		}

		if addrs != nil {
			t.Errorf("expected nil addrs, got %v", addrs)
		}

		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success with multiple names", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		ifaces := []*mdlayherwifi.Interface{
			{Name: "wlan0"},
			{Name: "eth0"},
			{Name: "wlp2s0"},
		}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := wifipkg.New(mockHandle)
		names, err := service.GetNames()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(names) != 3 {
			t.Fatalf("expected 3 names, got %d", len(names))
		}

		mockHandle.AssertExpectations(t)
	})

	t.Run("empty name list", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		ifaces := []*mdlayherwifi.Interface{}

		mockHandle.On("Interfaces").Return(ifaces, nil)

		service := wifipkg.New(mockHandle)
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
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, errTest)

		service := wifipkg.New(mockHandle)
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
	t.Parallel()

	mockHandle := new(MockWiFiHandle)
	service := wifipkg.New(mockHandle)

	if service.WiFi != mockHandle {
		t.Errorf("expected WiFi to be %v, got %v", mockHandle, service.WiFi)
	}
}
