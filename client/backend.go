package client

import (
	"github.com/sqshq/sampler/metadata"
	"net/http"
)

const (
	backendUrl             = "http://localhost:8080/api/v1"
	registrationPath       = "/registration"
	reportInstallationPath = "/report/installation"
	reportCrashPath        = "repost/crash"
)

type BackendClient struct {
	client http.Client
}

func NewBackendClient() *BackendClient {
	return &BackendClient{
		client: http.Client{},
	}
}

func (c *BackendClient) ReportInstallation(statistics *metadata.Statistics) {
	// TODO
}

func (c *BackendClient) ReportCrash() {
	// TODO
}

func (c *BackendClient) Register(key string) {
	// TODO
}
