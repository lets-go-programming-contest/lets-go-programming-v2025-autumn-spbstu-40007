package wifi

import (
	"net"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	t.Run("Успешное получение адресов", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		mac, _ := net.ParseMAC("00:11:22:33:44:55")
		expectedIface := &wifi.Interface{HardwareAddr: mac}

		mockHandle.On("Interfaces").Return([]*wifi.Interface{expectedIface}, nil)

		service := New(mockHandle)
		addresses, err := service.GetAddresses()

		require.NoError(t, err)
		assert.Len(t, addresses, 1)
		assert.Equal(t, mac, addresses[0])
		mockHandle.AssertExpectations(t)
	})

	t.Run("Ошибка при получении интерфейсов", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		expectedErr := assert.AnError

		mockHandle.On("Interfaces").Return(nil, expectedErr)

		service := New(mockHandle)
		addresses, err := service.GetAddresses()

		require.Error(t, err)
		assert.Nil(t, addresses)
		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Run("Успешное получение имён", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		expectedIface := &wifi.Interface{Name: "wlan0"}

		mockHandle.On("Interfaces").Return([]*wifi.Interface{expectedIface}, nil)

		service := New(mockHandle)
		names, err := service.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
		mockHandle.AssertExpectations(t)
	})

	t.Run("Ошибка при получении интерфейсов", func(t *testing.T) {
		mockHandle := new(MockWiFiHandle)
		expectedErr := assert.AnError

		mockHandle.On("Interfaces").Return(nil, expectedErr)

		service := New(mockHandle)
		names, err := service.GetNames()

		require.Error(t, err)
		assert.Nil(t, names)
		mockHandle.AssertExpectations(t)
	})
}
