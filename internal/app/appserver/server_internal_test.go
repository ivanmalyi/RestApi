package appserver

import (
	"github.com/ivanmalyi/RestApi/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUserCreate(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/users", nil)

	server := NewServer(teststore.New())
	server.ServeHTTP(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusOK)
}
