package controller

import (
	"bytes"
	"context"
	"friend_management/cmd/mocks"
	"friend_management/intenal/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConnectFriends(t *testing.T) {
	var jsonStr = []byte(`{"friends":["andy@example","john@example.com"]}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.BasicResponse
		mockError    error
	}{
		{
			name:        "retieve success",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: model.BasicResponse{
				Success: true,
			},
			expectedBody: "{\"success\":true}\n",
		},
		{
			name:         "retieve failed by incorrect input",
			bodyRequest:  bytes.NewBuffer([]byte("")),
			expectedBody: "Request body is invalid",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/connect", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("ConnectFriends", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := ConnectFriends(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetUser(t *testing.T) {
	var jsonStr = []byte(`{"email":"andy@example"}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.User
		mockError    error
	}{
		{
			name:        "retieve success",
			bodyRequest: bytes.NewBuffer(jsonStr),
			mockresponse: model.User{
				Email:   "a@gmail.com",
				Friends: []string{"b@gmail.com"},
			},
			expectedBody: "{\"email\":\"a@gmail.com\",\"friends\":[\"b@gmail.com\"],\"subscription\":null,\"blocked\":null}\n",
		},
		{
			name:         "retieve failed by incorrect input",
			bodyRequest:  bytes.NewBuffer([]byte("")),
			expectedBody: "{\"email\":\"\",\"friends\":null,\"subscription\":null,\"blocked\":null}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/friend", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("GetUser", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := GetUser(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	testCases := []struct {
		name         string
		expectedBody string
		mockresponse []model.User
		mockError    error
	}{
		{
			name:         "retieve success",
			expectedBody: "[{\"email\":\"a@gmail.com\",\"friends\":[\"b@gmail.com\"],\"subscription\":null,\"blocked\":null},{\"email\":\"b@gmail.com\",\"friends\":[\"a@gmail.com\"],\"subscription\":null,\"blocked\":null}]\n",
			mockresponse: []model.User{
				{
					Email:   "a@gmail.com",
					Friends: []string{"b@gmail.com"},
				},
				{
					Email:   "b@gmail.com",
					Friends: []string{"a@gmail.com"},
				},
			},
		},
		{
			name:         "Get users error",
			expectedBody: "null\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/friend/GetAll", nil)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("GetAllUsers", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := GetAllUsers(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestSendUpdate(t *testing.T) {
	var jsonStr = []byte(`{"sender":"a@gmail.com","text":"Hello World! b@gmail.com"}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.SendUpdateResponse
		mockError    error
	}{
		{
			name:         "retieve success",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":true,\"recipients\":[\"b@gmail.com\",\"c@gmail.com\",\"d@gmail.com\"]}\n",
			mockresponse: model.SendUpdateResponse{
				Success:    true,
				Recipients: []string{"b@gmail.com", "c@gmail.com", "d@gmail.com"},
			},
		},
		{
			name:         "retieve failed",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":false,\"recipients\":null}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/send", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("SendUpdate", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := SendUpdate(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestBlocked(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"a@gmail.com","target":"b@gmail.com"}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.BasicResponse
		mockError    error
	}{
		{
			name:         "retieve success",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":true}\n",
			mockresponse: model.BasicResponse{
				Success: true,
			},
		},
		{
			name:         "retieve failed",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":false}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/blocked", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("Blocked", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := Blocked(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestSubscription(t *testing.T) {
	var jsonStr = []byte(`{"requestor":"a@gmail.com","target":"b@gmail.com"}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.BasicResponse
		mockError    error
	}{
		{
			name:         "retieve success",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":true}\n",
			mockresponse: model.BasicResponse{
				Success: true,
			},
		},
		{
			name:         "retieve failed",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":false}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/subscribe", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("Subscription", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := Subscription(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestCommonFriends(t *testing.T) {
	var jsonStr = []byte(`{"friends":["a@gmail.com","b@gmail.com"]}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.FriendListResponse
		mockError    error
	}{
		{
			name:         "retieve success",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":true,\"friends\":[\"c@gmail.com\"],\"count\":1}\n",
			mockresponse: model.FriendListResponse{
				Success: true,
				Friends: []string{"c@gmail.com"},
				Count:   1,
			},
		},
		{
			name:         "retieve failed",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":false,\"friends\":null,\"count\":0}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/common", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("CommonFriends", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := CommonFriends(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestFriendList(t *testing.T) {
	var jsonStr = []byte(`{"email":"a@gmail.com"}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.FriendListResponse
		mockError    error
	}{
		{
			name:         "retieve success",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":true,\"friends\":[\"b@gmail.com\",\"c@gmail.com\"],\"count\":2}\n",
			mockresponse: model.FriendListResponse{
				Success: true,
				Friends: []string{"b@gmail.com", "c@gmail.com"},
				Count:   2,
			},
		},
		{
			name:         "retieve failed",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":false,\"friends\":null,\"count\":0}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/list", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("FriendList", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := FriendList(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestCreateNewUser(t *testing.T) {
	var jsonStr = []byte(`{"email":"a@gmail.com"}`)
	testCases := []struct {
		name         string
		bodyRequest  *bytes.Buffer
		expectedBody string
		mockresponse model.BasicResponse
		mockError    error
	}{
		{
			name:         "retieve success",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":true}\n",
			mockresponse: model.BasicResponse{
				Success: true,
			},
		},
		{
			name:         "retieve failed",
			bodyRequest:  bytes.NewBuffer(jsonStr),
			expectedBody: "{\"success\":false}\n",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/friend/addUser", tt.bodyRequest)
			require.NoError(t, err)
			chiCtx := chi.NewRouter()
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			serviceMock := new(mocks.Service)
			serviceMock.On("CreateNewUser", mock.Anything, mock.Anything).Return(tt.mockresponse, tt.mockError)
			handler := CreateNewUser(serviceMock)
			handler.ServeHTTP(rr, req)
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
