package spn

import (
	"encoding/json"
	"net/http"
	"time"
)

// UserStatus represent the data returned by the /save/status/user endpoint
type UserStatus struct {
	DailyCaptures      int `json:"daily_captures"`
	DailyCapturesLimit int `json:"daily_captures_limit"`
	Available          int `json:"available"`
	Processing         int `json:"processing"`
}

func (to *UserStatus) Update(from UserStatus) {
	to.DailyCaptures = from.DailyCaptures
	to.DailyCapturesLimit = from.DailyCapturesLimit
	to.Available = from.Available
	to.Processing = from.Processing
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

// GetUserStatus retrieve the user status for a given SPN account
func (c Connector) GetUserStatus() (userStatus UserStatus, err error) {

	// Build request
	req, err := http.NewRequest("GET", "https://web.archive.org/save/status/user", nil)
	if err != nil {
		return userStatus, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return userStatus, err
	}

	json.NewDecoder(resp.Body).Decode(&userStatus)

	return userStatus, nil
}

// Refresh the cached user status
func (c *Connector) refreshCachedUserStatus() {
	// Skip if we have enough slots and the status was updated less than 10s ago
	if c.cachedStatus.Available > 2 && time.Since(c.cachedStatusLastUpdated) < time.Second*10 {
		return
	}

	// Make sure we don't fetch the status too often (< 2s)
	wait := time.Second*2 - time.Since(c.cachedStatusLastUpdated)
	if wait > 0 {
		logger.Debug("Waiting before fetching user status", "wait", wait)
		time.Sleep(wait)
	}

	logger.Debug("Fetching user status")

	c.cachedStatusLastUpdated = time.Now()

	userStatus, err := c.GetUserStatus()
	if err != nil {
		logger.Error("Failed to fetch user status", "error", err)
		return
	}

	logger.Debug("User status fetched", "status", userStatus)
	c.cachedStatus.Update(userStatus)
	logger.Debug("cachedStatus updated", "cachedStatus", c.cachedStatus)
}

// Wait until a capture slot is available
func (c *Connector) GetAvailableCaptureSlot() (err error) {
	for {
		c.refreshCachedUserStatus()

		if c.cachedStatus.Available > 0 {
			c.cachedStatus.Available--
			c.cachedStatus.Processing++
			logger.Debug("AwaitAvailableSlot return", "cachedStatus", c.cachedStatus)
			return nil
		}

		logger.Debug("AwaitAvailableSlot waiting")
		time.Sleep(time.Second)
	}
}
