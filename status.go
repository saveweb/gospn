package spn

import (
	"encoding/json"
	"net/http"
)

// UserStatus represent the data returned by the /save/status/user endpoint
type UserStatus struct {
	DailyCaptures      int `json:"daily_captures"`
	DailyCapturesLimit int `json:"daily_captures_limit"`
	Available          int `json:"available"`
	Processing         int `json:"processing"`
}

// CaptureStatus represent the date returned by the /save/status/{job_id} endpoint
type CaptureStatus struct {
	Timestamp   string   `json:"timestamp"`
	DurationSec float64  `json:"duration_sec"`
	OriginalURL string   `json:"original_url"`
	Status      string   `json:"status"`
	StatusExt   string   `json:"status_ext"`
	JobID       string   `json:"job_id"`
	Outlinks    []string `json:"outlinks"`
	Resources   []string `json:"resources"`
	Exception   string   `json:"exception"`
	Message     string   `json:"message"`
}

// GetCaptureStatus retrieve the informations about a SPN job
func (c Connector) GetCaptureStatus(jobID string) (captureStatus CaptureStatus, err error) {
	// Build request
	req, err := http.NewRequest("GET", "https://web.archive.org/save/status/"+jobID, nil)
	if err != nil {
		return captureStatus, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return captureStatus, err
	}

	json.NewDecoder(resp.Body).Decode(&captureStatus)

	if captureStatus.Outlinks == nil {
		captureStatus.Outlinks = []string{""}
	}

	if captureStatus.Resources == nil {
		captureStatus.Resources = []string{""}
	}

	return captureStatus, nil
}

// GetAvailableCaptureSlots retrieve the available capture slots for a given SPN account
func (c Connector) GetAvailableCaptureSlots() (availableSlots int, err error) {
	var userStatus UserStatus

	// Build request
	req, err := http.NewRequest("GET", "https://web.archive.org/save/status/user", nil)
	if err != nil {
		return availableSlots, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return availableSlots, err
	}

	json.NewDecoder(resp.Body).Decode(&userStatus)
	availableSlots = userStatus.Available

	return availableSlots, nil
}
