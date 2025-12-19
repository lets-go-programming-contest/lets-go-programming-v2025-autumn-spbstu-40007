package wifi

import (
	"errors"
	"fmt"
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/mock"
)

var testError = errors.New("test error")

type MockWiFiHandle struct {
	mock.Mock
}

func (m *MockWiFiHandle) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()

	// Получаем интерфейсы
	interfacesRaw := args.Get(0)
	if interfacesRaw == nil {
		// Если интерфейсов нет, возвращаем обернутую ошибку
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("mock returned error: %w", err)
		}

		return nil, nil
	}

	// Приводим тип
	interfaces, ok := interfacesRaw.([]*wifi.Interface)
	if !ok {
		// Ошибка приведения типа
		if err := args.Error(1); err != nil {
			return nil, fmt.Errorf("type assertion failed with error: %w", err)
		}
		return nil, fmt.Errorf("type assertion failed: expected []*wifi.Interface, got %T", interfacesRaw)
	}

	// Проверяем ошибку
	if err := args.Error(1); err != nil {
		return interfaces, fmt.Errorf("mock error: %w", err)
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

		mockHandle.AssertExpectations(t)
	})

	t.Run("success with empty interface list", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, testError)

		service := New(mockHandle)
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

		mockHandle.AssertExpectations(t)
	})

	t.Run("empty name list", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		mockHandle.On("Interfaces").Return(nil, testError)

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
	t.Parallel()

	mockHandle := new(MockWiFiHandle)
	service := New(mockHandle)

	if service.WiFi != mockHandle {
		t.Errorf("expected WiFi to be %v, got %v", mockHandle, service.WiFi)
	}
}
