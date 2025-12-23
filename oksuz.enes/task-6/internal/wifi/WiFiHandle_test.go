package wifi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetWifiStatus(id int) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func TestWiFiHandle(t *testing.T) {
	mockRepo := new(MockRepository)

	mockRepo.On("GetWifiStatus", 1).Return("active", nil)
	assert.Equal(t, "WiFi is ON", WiFiHandle(mockRepo, 1))

	mockRepo.On("GetWifiStatus", 2).Return("inactive", nil)
	assert.Equal(t, "WiFi is OFF", WiFiHandle(mockRepo, 2))

	mockRepo.On("GetWifiStatus", 3).Return("", errors.New("db error"))
	assert.Equal(t, "Error: Database Connection", WiFiHandle(mockRepo, 3))
}
