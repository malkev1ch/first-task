package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/malkev1ch/first-task/internal/config"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/malkev1ch/first-task/internal/service"
	mock_service "github.com/malkev1ch/first-task/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, input *model.CreateUser)
	ctx := context.Background()
	testTable := []struct {
		name                string
		ctx                 context.Context
		inputBody           string
		inputUser           *model.CreateUser
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockAuth, input *model.CreateUser) {
				s.EXPECT().SignUp(ctx, input).Return(&model.Tokens{
					RefreshToken: "qwerty",
					AccessToken:  "qwerty",
				}, nil)
			},
			ctx:       ctx,
			inputBody: `{"email":"qwerty@gmail.com", "password":"ZAQ!2wsxCDE#", "userName":"Some name"}`,
			inputUser: &model.CreateUser{
				UserName: "Some name",
				Email:    "qwerty@gmail.com",
				Password: "ZAQ!2wsxCDE#",
			},
			expectedStatusCode:  http.StatusCreated,
			expectedRequestBody: `{"accessToken":"qwerty","refreshToken":"qwerty"}` + "\n",
		},
		{
			name: "Service Error",
			mockBehavior: func(s *mock_service.MockAuth, input *model.CreateUser) {
				s.EXPECT().SignUp(ctx, input).Return(nil, errors.New("service error"))
			},
			ctx:       ctx,
			inputBody: `{"email":"qwerty@gmail.com", "password":"ZAQ!2wsxCDE#", "userName":"Some name"}`,
			inputUser: &model.CreateUser{
				UserName: "Some name",
				Email:    "qwerty@gmail.com",
				Password: "ZAQ!2wsxCDE#",
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"can't create user","error":"service error"}` + "\n",
		},
		{
			name:                "Invalid Password",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.CreateUser) {},
			ctx:                 ctx,
			inputBody:           `{"email":"qwerty@gmail.com", "password":"ZAQ!", "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'CreateUser.Password' Error:Field validation for 'Password' failed on the 'gt' tag"}` + "\n",
		},
		{
			name:                "Invalid email",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.CreateUser) {},
			ctx:                 ctx,
			inputBody:           `{"email":"qwertygmail.com", "password":"ZAQ!2wsxCDE#",  "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'CreateUser.Email' Error:Field validation for 'Email' failed on the 'email' tag"}` + "\n",
		},
		{
			name:                "Request without email",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.CreateUser) {},
			ctx:                 ctx,
			inputBody:           `{"password":"ZAQ!2wsxCDE#", "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'CreateUser.Email' Error:Field validation for 'Email' failed on the 'required' tag"}` + "\n",
		},
		{
			name:                "Request without password",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.CreateUser) {},
			ctx:                 ctx,
			inputBody:           `{"email":"qwerty@gmail.com", "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'CreateUser.Password' Error:Field validation for 'Password' failed on the 'required' tag"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init dependencies
			c := gomock.NewController(t)
			mockAuth := mock_service.NewMockAuth(c)
			testCase.mockBehavior(mockAuth, testCase.inputUser)
			services := &service.Service{Auth: mockAuth}
			cfg := config.Config{}
			validator := NewValidator()
			handler := NewHandler(services, &cfg, validator)

			// Init server
			r := InitRouter(handler, &cfg)

			// Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Set request headers
			req.Header.Set("Content-Type", "application/json")

			// Execute the request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestSignIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, input *model.AuthUser)
	ctx := context.Background()
	testTable := []struct {
		name                string
		ctx                 context.Context
		inputBody           string
		inputUser           *model.AuthUser
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockAuth, input *model.AuthUser) {
				s.EXPECT().SignIn(ctx, input).Return(&model.Tokens{
					RefreshToken: "qwerty",
					AccessToken:  "qwerty",
				}, nil)
			},
			ctx:       ctx,
			inputBody: `{"email":"qwerty@gmail.com", "password":"ZAQ!2wsxCDE#", "userName":"Some name"}`,
			inputUser: &model.AuthUser{
				Email:    "qwerty@gmail.com",
				Password: "ZAQ!2wsxCDE#",
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"accessToken":"qwerty","refreshToken":"qwerty"}` + "\n",
		},
		{
			name: "Service Error",
			mockBehavior: func(s *mock_service.MockAuth, input *model.AuthUser) {
				s.EXPECT().SignIn(ctx, input).Return(nil, errors.New("service error"))
			},
			ctx:       ctx,
			inputBody: `{"email":"qwerty@gmail.com", "password":"ZAQ!2wsxCDE#", "userName":"Some name"}`,
			inputUser: &model.AuthUser{
				Email:    "qwerty@gmail.com",
				Password: "ZAQ!2wsxCDE#",
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"authorisation failed","error":"service error"}` + "\n",
		},
		{
			name:                "Invalid Password",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.AuthUser) {},
			ctx:                 ctx,
			inputBody:           `{"email":"qwerty@gmail.com", "password":"ZAQ!", "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'AuthUser.Password' Error:Field validation for 'Password' failed on the 'gt' tag"}` + "\n",
		},
		{
			name:                "Invalid email",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.AuthUser) {},
			ctx:                 ctx,
			inputBody:           `{"email":"qwertygmail.com", "password":"ZAQ!2wsxCDE#",  "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'AuthUser.Email' Error:Field validation for 'Email' failed on the 'email' tag"}` + "\n",
		},
		{
			name:                "Request without email",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.AuthUser) {},
			ctx:                 ctx,
			inputBody:           `{"password":"ZAQ!2wsxCDE#", "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'AuthUser.Email' Error:Field validation for 'Email' failed on the 'required' tag"}` + "\n",
		},
		{
			name:                "Request without password",
			mockBehavior:        func(s *mock_service.MockAuth, input *model.AuthUser) {},
			ctx:                 ctx,
			inputBody:           `{"email":"qwerty@gmail.com", "userName":"Some name"}`,
			inputUser:           nil,
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields in json body or wrong values of fields","error":"Key: 'AuthUser.Password' Error:Field validation for 'Password' failed on the 'required' tag"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init dependencies
			c := gomock.NewController(t)
			mockAuth := mock_service.NewMockAuth(c)
			testCase.mockBehavior(mockAuth, testCase.inputUser)
			services := &service.Service{Auth: mockAuth}
			cfg := config.Config{}
			validator := NewValidator()
			handler := NewHandler(services, &cfg, validator)

			// Init server
			r := InitRouter(handler, &cfg)

			// Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBufferString(testCase.inputBody))

			// Set request headers
			req.Header.Set("Content-Type", "application/json")

			// Execute the request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestRefresh(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuth, input string)
	ctx := context.Background()
	testTable := []struct {
		name                string
		ctx                 context.Context
		inputBody           string
		inputService        string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockAuth, input string) {
				s.EXPECT().RefreshToken(ctx, input).Return(&model.Tokens{
					AccessToken:  "qwerty",
					RefreshToken: "qwerty",
				}, nil)
			},
			ctx:                 ctx,
			inputBody:           `{"refreshToken":"qwerty"}`,
			inputService:        "qwerty",
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"accessToken":"qwerty","refreshToken":"qwerty"}` + "\n",
		},
		{
			name: "Service Error",
			mockBehavior: func(s *mock_service.MockAuth, input string) {
				s.EXPECT().RefreshToken(ctx, input).Return(nil, errors.New("service error"))
			},
			ctx:                 ctx,
			inputBody:           `{"refreshToken":"qwerty"}`,
			inputService:        "qwerty",
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"failed refresh token","error":"service error"}` + "\n",
		},
		{
			name:                "Request without refresh token",
			mockBehavior:        func(s *mock_service.MockAuth, input string) {},
			ctx:                 ctx,
			inputBody:           `{"accessToken":"qwerty"}`,
			inputService:        "qwerty",
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"not enough fields","error":"Key: 'RefreshToken.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required' tag"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init dependencies
			c := gomock.NewController(t)
			mockAuth := mock_service.NewMockAuth(c)
			testCase.mockBehavior(mockAuth, testCase.inputService)
			services := &service.Service{Auth: mockAuth}
			cfg := config.Config{}
			validator := NewValidator()
			handler := NewHandler(services, &cfg, validator)

			// Init server
			r := InitRouter(handler, &cfg)

			// Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBufferString(testCase.inputBody))

			// Set request headers
			req.Header.Set("Content-Type", "application/json")

			// Execute the request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
