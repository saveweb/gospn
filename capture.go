package spn

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Capture execute a capture via https://web.archive.org/save
// and return the response. Options for the capture can be specified
// when calling the method
func (c Connector) Capture(URL string, options ...string) (captureResponse CaptureResponse, err error) {
	// Check if there are available slots for capture, if not then we wait a bit
	for {
		freeSlots, err := c.GetAvailableCaptureSlots()
		if err != nil {
			return captureResponse, err
		}

		if freeSlots > 0 {
			break
		} else {
			time.Sleep(time.Second)
			continue
		}
	}

	// Build request body
	dataString := "url=" + URL
	for _, option := range options {
		dataString = dataString + "&" + option
	}

	// Build request
	req, err := http.NewRequest("POST", "https://web.archive.org/save", bytes.NewBufferString(dataString))
	if err != nil {
		return captureResponse, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "LOW "+c.AccessKey+":"+c.SecretKey)

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return captureResponse, err
	}

	json.NewDecoder(resp.Body).Decode(&captureResponse)

	return captureResponse, nil
}
