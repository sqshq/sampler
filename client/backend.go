package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sqshq/sampler/metadata"
	"io/ioutil"
	"net/http"
)

const (
	backendUrl       = "http://localhost/api/v1"
	installationPath = "/telemetry/installation"
	statisticsPath   = "/telemetry/statistics"
	crashPath        = "/telemetry/crash"
	registrationPath = "/license/registration"
	jsonContentType  = "application/json"
)

// Backend client is used to verify license and to send telemetry reports
// for analyses (anonymous usage data statistics and crash reports)
type BackendClient struct {
	client http.Client
}

func NewBackendClient() *BackendClient {
	return &BackendClient{
		client: http.Client{},
	}
}

func (c *BackendClient) ReportInstallation(statistics *metadata.Statistics) {

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(statistics)
	if err != nil {
		c.ReportCrash(err.Error(), statistics)
	}

	_, err = http.Post(backendUrl+installationPath, jsonContentType, buf)
	if err != nil {
		c.ReportCrash(err.Error(), statistics)
	}
}

func (c *BackendClient) ReportUsageStatistics(error string, statistics *metadata.Statistics) {
	// TODO
}

func (c *BackendClient) ReportCrash(error string, statistics *metadata.Statistics) {
	// TODO
}

func (c *BackendClient) RegisterLicenseKey(licenseKey string, statistics *metadata.Statistics) (*metadata.License, error) {

	req := struct {
		LicenseKey string
		Statistics *metadata.Statistics
	}{
		licenseKey,
		statistics,
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		c.ReportCrash(err.Error(), statistics)
	}

	response, err := http.Post(
		backendUrl+registrationPath, jsonContentType, buf)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, errors.New(string(body))
	}

	var license metadata.License
	json.NewDecoder(response.Body).Decode(&license)

	return &license, nil
}
