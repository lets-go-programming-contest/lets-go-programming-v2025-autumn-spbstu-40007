package wifi_test

import (
	"errors"
	"testing"

	"github.com/mdlayher/wifi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	service "github.com/ami0-0/task-6/internal/wifi"
)

type MockScanner struct {
	mock.Mock
}

func (m *MockScanner) Interfaces() ([]*wifi.Interface, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*wifi.Interface), args.Error(1)
}

func TestGetInterfaceNames(t *testing.T) {
	m := new(MockScanner)
	svc := service.NewNetworkService(m)

	t.Run("success", func(t *testing.T) {
		fakeIfaces := []*wifi.Interface{{Name: "wlan0"}}
		m.On("Interfaces").Return(fakeIfaces, nil).Once()

		names, err := svc.GetInterfaceNames()
		assert.NoError(t, err)
		assert.Equal(t, []string{"wlan0"}, names)
		m.AssertExpectations(t)
	})

	t.Run("fail", func(t *testing.T) {
		m.On("Interfaces").Return(nil, errors.New("hardware error")).Once()
		_, err := svc.GetInterfaceNames()
		assert.Error(t, err)
	})
}