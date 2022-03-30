package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/malkev1ch/first-task/internal/config"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/malkev1ch/first-task/internal/service"
	mock_service "github.com/malkev1ch/first-task/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCat(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCat, id string)
	ctx := context.Background()
	testTable := []struct {
		name                string
		catID               string
		ctx                 context.Context
		mockBehavior        mockBehavior
		expectedCat         *model.Cat
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Unauthorized",
			catID:               "c9d1684e-afbf-4dd8-82bc-7932cc5b7b08",
			mockBehavior:        func(s *mock_service.MockCat, id string) {},
			ctx:                 ctx,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedRequestBody: `{"message":"invalid or expired jwt"}` + "\n",
		},
		{
			name:  "OK",
			catID: "9d9044a6-d8a8-4e8c-9132-e583d2ebd6c4",
			mockBehavior: func(s *mock_service.MockCat, id string) {
				s.EXPECT().Get(ctx, id).Return(&model.Cat{
					ID:         "9d9044a6-d8a8-4e8c-9132-e583d2ebd6c4",
					Name:       "New Cat",
					DateBirth:  time.Date(2018, 9, 22, 12, 42, 31, 0, time.UTC),
					Vaccinated: true,
					ImagePath:  "",
				}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"9d9044a6-d8a8-4e8c-9132-e583d2ebd6c4","name":"New Cat","dateBirth":"2018-09-22T12:42:31Z","vaccinated":true}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			mockCat := mock_service.NewMockCat(c)
			testCase.mockBehavior(mockCat, testCase.catID)
			services := &service.Service{Cat: mockCat}
			cfg := config.Config{}
			validator := NewValidator()
			handler := NewHandler(services, &cfg, validator)

			r := InitRouter(handler, &cfg)
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/cats/%s", testCase.catID), nil)

			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiIiwiZW1haWwiOiJxd2VydHlhc2RAZ21haWwuY29tIiwiZXhwIjoxNjUxMjQ2ODc4LCJqdGkiOiJhMjYwYTk2Mi03NThmLTQzNWEtYWE0NS0wNDkxYWI4NTg4ODgifQ.RgNsMV95kz16_YGTrHMAfURGG0MqOKkL29_r7f7n-70")

			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
