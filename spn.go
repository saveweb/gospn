package spn

import (
	"net/http"
	"time"
)

// Connector represent the necessary data to execute SPN requests
type Connector struct {
	AccessKey  string
	SecretKey  string
	HTTPClient *http.Client
}

// CaptureResponse represent the JSON response from SPN
// returned when a capture is executed
type CaptureResponse struct {
	URL       string `json:"url"`
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	StatusExt string `json:"status_ext"`
	Message   string `json:"message"`
}

// Init initialize the SPN connector that can be used
// to trigger archiving for an URL
func Init(accessKey, secretKey string) (Connector, error) {
	var connector Connector

	connector.AccessKey = accessKey
	connector.SecretKey = secretKey

	connector.HTTPClient = &http.Client{
		Timeout: time.Second * 15,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// TODO: test keys validity?

	return connector, nil
}
