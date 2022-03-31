package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/malkev1ch/first-task/internal/model"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateCat(t *testing.T) {
	testTable := []struct {
		name               string
		requestBody        *model.CreateCat
		expectedStatusCode int
	}{
		{
			name: "OK",
			requestBody: &model.CreateCat{
				Name:       "Some Name",
				DateBirth:  time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
				Vaccinated: false,
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Bad request",
			requestBody: &model.CreateCat{
				DateBirth:  time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
				Vaccinated: false,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			tr := &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 20,
			}
			netClient := &http.Client{Transport: tr}
			postBody, _ := json.Marshal(testCase.requestBody)
			res, err := netClient.Post(fmt.Sprintf("http://127.0.0.1:%s/cats/", golangPort),
				"application/json", bytes.NewBuffer(postBody))
			if err != nil {
				logrus.Fatal(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Info(string(body))
			defer res.Body.Close()
			assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestGetCat(t *testing.T) {
	testTable := []struct {
		name               string
		catID              string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			catID:              "9d9044a6-d8a8-4e8c-9132-e583d2ebd6c4",
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/cats/%s", golangPort, testCase.catID))
			if err != nil {
				logrus.Fatal(err)
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Fatal(err)
			}
			logrus.Info(string(body))
			defer res.Body.Close()
			assert.Equal(t, testCase.expectedStatusCode, res.StatusCode)
		})
	}
}
