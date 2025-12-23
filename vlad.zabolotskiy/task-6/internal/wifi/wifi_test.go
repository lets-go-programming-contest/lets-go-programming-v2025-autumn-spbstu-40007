package wifi_test

import (
	"errors"
	"net"
	"testing"

	service "github.com/se1lzor/task-6/internal/wifi"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

var errWiFiFail = errors.New("wifi fail")

func TestWiFiService_GetAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		svc := service.New(mockHandle)

		hw, _ := net.ParseMAC("00:11:22:33:44:55")

		mockHandle.On("Interfaces").Return([]*wifi.Interface{
			{HardwareAddr: hw},
		}, nil).Once()

		got, err := svc.GetAddresses()
		require.NoError(t, err)
		require.Equal(t, []net.HardwareAddr{hw}, got)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return(([]*wifi.Interface)(nil), errWiFiFail).Once()

		got, err := svc.GetAddresses()
		require.Error(t, err)
		require.Nil(t, got)
		require.ErrorIs(t, err, errWiFiFail)
		require.Contains(t, err.Error(), "getting interfaces")

		mockHandle.AssertExpectations(t)
	})
}

func TestWiFiService_GetNames(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return([]*wifi.Interface{
			{Name: "wlan0"},
		}, nil).Once()

		got, err := svc.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"wlan0"}, got)

		mockHandle.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockHandle := new(MockWiFiHandle)
		svc := service.New(mockHandle)

		mockHandle.On("Interfaces").Return(([]*wifi.Interface)(nil), errWiFiFail).Once()

		got, err := svc.GetNames()
		require.Error(t, err)
		require.Nil(t, got)
		require.ErrorIs(t, err, errWiFiFail)
		require.Contains(t, err.Error(), "getting interfaces")

		mockHandle.AssertExpectations(t)
	})
}
