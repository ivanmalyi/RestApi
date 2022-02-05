package appserver

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/ivanmalyi/RestApi/internal/app/model"
	"github.com/ivanmalyi/RestApi/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUserCreate(t *testing.T) {
	server := NewServer(teststore.New(),  sessions.NewCookieStore([]byte("secret")))
	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	} {
		{
			name: "valid",
			payload: map[string]string {
				"email":"ivan.malyi93@gmail.com",
				"password":"123456",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: "bad request",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string {
				"email":"invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc:=range testCases {
		t.Run(tc.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			b := &bytes.Buffer{}
			_ = json.NewEncoder(b).Encode(tc.payload)
			request := httptest.NewRequest(http.MethodPost, "/users", b)
			server.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, tc.expectedCode)
		})
	}
}

func TestServer_HandleSessionCreate(t *testing.T) {

	user := model.TestUser(t)
	store := teststore.New()
	_ = store.User().Create(user)
	server := NewServer(store, sessions.NewCookieStore([]byte("secret")))

	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	} {
		{
			name: "valid",
			payload: map[string]string {
				"email":user.Email,
				"password":user.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid payload",
			payload: "bad request",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string {
				"email":"invalid_email",
				"password":user.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string {
				"email":user.Email,
				"password":"invalid password",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc:=range testCases {
		t.Run(tc.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			b := &bytes.Buffer{}
			_ = json.NewEncoder(b).Encode(tc.payload)
			request := httptest.NewRequest(http.MethodPost, "/session", b)
			server.ServeHTTP(recorder, request)

			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}
