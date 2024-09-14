package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolosHGW/goph-keeper/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthHandler struct {
	mock.Mock
}

func (m *MockAuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	w.WriteHeader(http.StatusOK)
}

func TestNewRouter_RegisterRoute(t *testing.T) {
	mockAuthHandler := new(MockAuthHandler)

	mockAuthHandler.On("RegisterUser", mock.Anything, mock.Anything).Return()

	handlers := &handler.Handlers{
		RegisterHandler: mockAuthHandler,
	}

	r := NewRouter(handlers)

	req, _ := http.NewRequest(http.MethodPost, "/api/user/register", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockAuthHandler.AssertCalled(t, "RegisterUser", mock.Anything, mock.Anything)
}

func TestNewRouter_UnknownRoute(t *testing.T) {
	mockAuthHandler := new(MockAuthHandler)

	mockAuthHandler.On("RegisterUser", mock.Anything, mock.Anything).Return()

	handlers := &handler.Handlers{
		RegisterHandler: mockAuthHandler,
	}

	r := NewRouter(handlers)

	req, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
