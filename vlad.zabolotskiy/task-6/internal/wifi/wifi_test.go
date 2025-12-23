package wifi_test

import (
	"errors"
	"net"
	"testing"

	service "github.com/se1lzor/task-6/internal/wifi" // <-- замени

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/require"
)

func TestWiFiService_Compact_WithTestifyMock(t *testing.T) {
	t.Parallel()

	mockHandle := new(MockWiFiHandle)
	svc := service.New(mockHandle)

	t.Run("GetAddresses ok/error", func(t *testing.T) {
		t.Parallel()

		hw, _ := net.ParseMAC("00:11:22:33:44:55")

		mockHandle.On("Interfaces").Return([]*wifi.Interface{
			{HardwareAddr: hw},
		}, nil).Once()

		got, err := svc.GetAddresses()
		require.NoError(t, err)
		require.Equal(t, []net.HardwareAddr{hw}, got)

		sentinel := errors.New("wifi fail")
		mockHandle.On("Interfaces").Return(([]*wifi.Interface)(nil), sentinel).Once()

		got, err = svc.GetAddresses()
		require.Error(t, err)
		require.Nil(t, got)
		require.ErrorIs(t, err, sentinel)
		require.Contains(t, err.Error(), "getting interfaces")

		mockHandle.AssertExpectations(t)
	})

	t.Run("GetNames ok/error", func(t *testing.T) {
		t.Parallel()

		mockHandle.On("Interfaces").Return([]*wifi.Interface{
			{Name: "wlan0"},
		}, nil).Once()

		got, err := svc.GetNames()
		require.NoError(t, err)
		require.Equal(t, []string{"wlan0"}, got)

		sentinel := errors.New("wifi fail")
		mockHandle.On("Interfaces").Return(([]*wifi.Interface)(nil), sentinel).Once()

		got, err = svc.GetNames()
		require.Error(t, err)
		require.Nil(t, got)
		require.ErrorIs(t, err, sentinel)
		require.Contains(t, err.Error(), "getting interfaces")

		mockHandle.AssertExpectations(t)
	})
}
