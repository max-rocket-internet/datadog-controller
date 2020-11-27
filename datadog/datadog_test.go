package datadog

import (
	"bytes"
	"github.com/max-rocket-internet/datadog-controller/api/v1beta1"
	"github.com/max-rocket-internet/datadog-controller/datadog/mocks"
	"github.com/max-rocket-internet/datadog-controller/datadog/restclient"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func init() {
	restclient.Client = &mocks.MockClient{}
	os.Setenv("DD_CLIENT_API_KEY", "INVALID_API_KEY")
	os.Setenv("DD_CLIENT_APP_KEY", "INVALID_APP_KEY")
}

var (
	apiKeyValidResponseJson = `{"valid": true}`
)

// tests to do:
//   - Auth with bad API key
//   - Config errors
//   - Log level set correct from env and flag

func TestDeleteMonitor(t *testing.T) {
	responseJson := `{"deleted_monitor_id": 12345}`
	responseJsonBody := ioutil.NopCloser(bytes.NewReader([]byte(responseJson)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body := responseJsonBody

		if req.URL.Path == "/api/v1/validate" {
			body = ioutil.NopCloser(bytes.NewReader([]byte(apiKeyValidResponseJson)))
		}

		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}

	datadogApi, err := New("INFO")
	assert.Nil(t, err)

	err = datadogApi.DeleteMonitor(12345)
	assert.Nil(t, err)
}

func TestDeleteMonitorFail(t *testing.T) {
	responseJson := `{"deleted_monitor_id": 12345}`
	responseJsonBody := ioutil.NopCloser(bytes.NewReader([]byte(responseJson)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body := responseJsonBody

		if req.URL.Path == "/api/v1/validate" {
			body = ioutil.NopCloser(bytes.NewReader([]byte(apiKeyValidResponseJson)))
		}

		return &http.Response{
			StatusCode: 400,
			Body:       body,
		}, nil
	}

	datadogApi, err := New("INFO")
	assert.Nil(t, err)

	err = datadogApi.DeleteMonitor(12345)
	assert.NotNil(t, err)
}

func TestCreateMonitor(t *testing.T) {
	newMonitor := v1beta1.DatadogMonitorSpec{}
	newMonitor.Name = "test-create"
	newMonitor.Message = "test-message"
	newMonitor.Query = "test-query"
	newMonitor.Type = "query alert"

	responseJson := `{"id": 12345}`
	responseJsonBody := ioutil.NopCloser(bytes.NewReader([]byte(responseJson)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body := responseJsonBody

		if req.URL.Path == "/api/v1/validate" {
			body = ioutil.NopCloser(bytes.NewReader([]byte(apiKeyValidResponseJson)))
		}

		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}

	datadogApi, err := New("INFO")
	assert.Nil(t, err)

	monitorId, err := datadogApi.CreateMonitor(newMonitor)
	assert.EqualValues(t, monitorId, 12345)
	assert.Nil(t, err)
}

func TestCreateMonitorFail(t *testing.T) {
	newMonitor := v1beta1.DatadogMonitorSpec{}
	newMonitor.Name = "test-create"
	newMonitor.Message = "test-message"
	newMonitor.Query = "test-query"
	newMonitor.Type = "query alert"

	responseJson := `{"id": 12345}`
	responseJsonBody := ioutil.NopCloser(bytes.NewReader([]byte(responseJson)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body := responseJsonBody

		if req.URL.Path == "/api/v1/validate" {
			body = ioutil.NopCloser(bytes.NewReader([]byte(apiKeyValidResponseJson)))
		}

		return &http.Response{
			StatusCode: 400,
			Body:       body,
		}, nil
	}

	datadogApi, err := New("INFO")
	assert.Nil(t, err)

	_, err = datadogApi.CreateMonitor(newMonitor)
	assert.NotNil(t, err)
}

func TestUpdateMonitor(t *testing.T) {
	updatedMonitor := v1beta1.DatadogMonitorSpec{}
	updatedMonitor.Name = "test-create"
	updatedMonitor.Message = "test-message"
	updatedMonitor.Query = "test-query"
	updatedMonitor.Type = "query alert"

	responseJson := `{"id": 12345}`
	responseJsonBody := ioutil.NopCloser(bytes.NewReader([]byte(responseJson)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body := responseJsonBody

		if req.URL.Path == "/api/v1/validate" {
			body = ioutil.NopCloser(bytes.NewReader([]byte(apiKeyValidResponseJson)))
		}

		return &http.Response{
			StatusCode: 200,
			Body:       body,
		}, nil
	}

	datadogApi, err := New("INFO")
	assert.Nil(t, err)

	err = datadogApi.UpdateMonitor(12345, updatedMonitor)
	assert.Nil(t, err)
}

func TestUpdateMonitorFail(t *testing.T) {
	updatedMonitor := v1beta1.DatadogMonitorSpec{}
	updatedMonitor.Name = "test-create"
	updatedMonitor.Message = "test-message"
	updatedMonitor.Query = "test-query"
	updatedMonitor.Type = "query alert"

	responseJson := `{"id": 12345}`
	responseJsonBody := ioutil.NopCloser(bytes.NewReader([]byte(responseJson)))

	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		body := responseJsonBody

		if req.URL.Path == "/api/v1/validate" {
			body = ioutil.NopCloser(bytes.NewReader([]byte(apiKeyValidResponseJson)))
		}

		return &http.Response{
			StatusCode: 400,
			Body:       body,
		}, nil
	}

	datadogApi, err := New("INFO")
	assert.Nil(t, err)

	err = datadogApi.UpdateMonitor(12345, updatedMonitor)
	assert.NotNil(t, err)
}
