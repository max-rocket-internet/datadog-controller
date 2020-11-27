package datadog

import (
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/max-rocket-internet/datadog-controller/api/v1beta1"
	"github.com/max-rocket-internet/datadog-controller/datadog/restclient"
	"github.com/max-rocket-internet/datadog-controller/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io/ioutil"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"strconv"
	"time"
)

type ApiKeyValidationResponse struct {
	Valid bool     `json:"valid"`
	Error []string `json:"errors"`
}

type MonitorDeleteResponse struct {
	DeletedMonitorId int      `json:"deleted_monitor_id"`
	Error            []string `json:"errors"`
}

type Config struct {
	datadogApiKey      string
	datadogAppKey      string
	DatadogHost        string
	datadogApiEndpoint string
	LogLevel           string
	apiBase            string
}

type Datadog struct {
	Log  logr.Logger
	Conf *Config
}

func newConfig() (*Config, error) {
	if err := utils.CheckRequiredEnvVars([]string{"DD_CLIENT_API_KEY", "DD_CLIENT_APP_KEY"}); err != nil {
		return nil, err
	}

	c := &Config{}
	c.datadogApiKey, _ = utils.GetEnvString("DD_CLIENT_API_KEY")
	c.datadogAppKey, _ = utils.GetEnvString("DD_CLIENT_APP_KEY")
	c.DatadogHost, _ = utils.GetEnvString("DATADOG_HOST", "datadoghq.eu")
	c.LogLevel, _ = utils.GetEnvString("LOG_LEVEL", "INFO")
	c.apiBase, _ = utils.GetEnvString("API_BASE", "/api/v1")
	c.datadogApiEndpoint = fmt.Sprintf("https://api.%s%s", c.DatadogHost, c.apiBase)

	return c, nil
}

var (
	httpUserAgent = "github/max-rocket-internet/datadog-controller/1.0"
	httpHeaders   = http.Header{}

	apiLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "datadog_controller",
		Subsystem: "api",
		Name:      "latency",
		Help:      "Latency of the Datadog API",
	}, []string{
		"endpoint",
		"response",
	})

	monitorEventCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "datadog_controller",
		Subsystem: "monitor",
		Name:      "event",
		Help:      "Count of monitor create/delete events",
	}, []string{
		"action",
	})
)

func (d Datadog) validateApiKey() error {
	d.Log.V(1).Info("Testing API token")

	response := ApiKeyValidationResponse{}

	results, _, err := d.apiRequest("GET", "/validate", nil)
	if err != nil {
		return fmt.Errorf("Error validating API key: %v", err)
	}

	err = json.Unmarshal(results, &response)
	if err != nil {
		return fmt.Errorf("Error unmarshalling API key validation response: %v", err)
	}

	if response.Valid {
		d.Log.V(1).Info("API key validated")
		return nil
	} else {
		return fmt.Errorf("API key invalid")
	}
}

func (d Datadog) DeleteMonitor(MonitorId int64) error {
	d.Log.V(1).Info(fmt.Sprintf("Deleting monitor %v", MonitorId))

	response := MonitorDeleteResponse{}

	results, responseCode, err := d.apiRequest("DELETE", fmt.Sprintf("/monitor/%v", MonitorId), nil)
	if err != nil {
		return err
	}

	if responseCode != 404 && responseCode != 200 {
		return fmt.Errorf("Error deleting monitor: %v", string(results))
	}

	err = json.Unmarshal(results, &response)
	if err != nil {
		return err
	}

	monitorEventCounter.WithLabelValues("deleted").Inc()

	return nil
}

func (d Datadog) CreateMonitor(MonitorSpec v1beta1.DatadogMonitorSpec) (int64, error) {
	d.Log.V(1).Info(fmt.Sprintf("Creating monitor '%v'", MonitorSpec.Name))

	requestBody, _ := json.Marshal(MonitorSpec)

	results, responseCode, err := d.apiRequest("POST", "/monitor", requestBody)
	if err != nil {
		monitorEventCounter.WithLabelValues("failed").Inc()
		return 0, err
	}

	if responseCode == 400 {
		monitorEventCounter.WithLabelValues("failed").Inc()
		return 0, fmt.Errorf("Error creating monitor '%v': %v", MonitorSpec.Name, string(results))
	}

	requestRespone := v1beta1.DatadogMonitorSpec{}
	err = json.Unmarshal(results, &requestRespone)
	if err != nil {
		monitorEventCounter.WithLabelValues("failed").Inc()
		return 0, err
	}

	monitorEventCounter.WithLabelValues("created").Inc()

	return requestRespone.Id, nil
}

func (d Datadog) UpdateMonitor(MonitorId int64, MonitorSpec v1beta1.DatadogMonitorSpec) error {
	d.Log.V(1).Info(fmt.Sprintf("Updating monitor '%v'", MonitorId))

	requestBody, _ := json.Marshal(MonitorSpec)

	results, responseCode, err := d.apiRequest("PUT", fmt.Sprintf("/monitor/%v", MonitorId), requestBody)
	if err != nil {
		monitorEventCounter.WithLabelValues("failed").Inc()
		return err
	}

	if responseCode == 400 {
		monitorEventCounter.WithLabelValues("failed").Inc()
		return fmt.Errorf("Error updating monitor '%v': %v", MonitorId, string(results))
	}

	requestRespone := v1beta1.DatadogMonitorSpec{}
	err = json.Unmarshal(results, &requestRespone)
	if err != nil {
		monitorEventCounter.WithLabelValues("failed").Inc()
		return err
	}

	monitorEventCounter.WithLabelValues("updated").Inc()

	return nil
}

func (d Datadog) apiRequest(RequestMethod string, RequestPath string, RequestBody []byte) ([]byte, int, error) {
	start := time.Now()
	resp, err := restclient.Do(RequestMethod, d.Conf.datadogApiEndpoint+RequestPath, RequestBody, httpHeaders)
	apiLatency.WithLabelValues(RequestPath, strconv.Itoa(resp.StatusCode)).Observe(time.Since(start).Seconds())
	if err != nil {
		d.Log.Error(err, fmt.Sprintf("Error making %v request for path %v", RequestMethod, RequestPath))
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	d.Log.V(1).Info(fmt.Sprintf("API response %v: %v (%v)", resp.StatusCode, RequestPath, RequestMethod))

	body, err := ioutil.ReadAll(resp.Body)

	// Need to handle error: {\"errors\":[\"Can not create duplicate monitors: Rate limit of 5 requests in 600 seconds reached. Please try again later.\"]}
	if resp.StatusCode == 429 {
		return nil, resp.StatusCode, fmt.Errorf("Error in apiRequest to %v (%v) %v: %v", RequestPath, RequestMethod, resp.StatusCode, string(body))
	}

	if err != nil {
		d.Log.Error(err, fmt.Sprintf("Error reading response body from %v request for path %v", RequestMethod, RequestPath))
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func New(logLevel string) (Datadog, error) {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	// How to set logLevel here?

	d := Datadog{}
	d.Log = ctrl.Log.WithName("datadog-api")

	config, err := newConfig()
	if err != nil {
		return d, err
	}

	httpHeaders.Add("DD-API-KEY", config.datadogApiKey)
	httpHeaders.Add("DD-APPLICATION-KEY", config.datadogAppKey)
	httpHeaders.Add("Content-Type", "application/json")
	httpHeaders.Add("User-Agent", httpUserAgent)

	d.Conf = config

	if err = d.validateApiKey(); err != nil {
		return d, err
	}

	return d, nil
}
