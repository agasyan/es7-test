package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Metric struct {
	licenseKey string
	accID      string
	env        string
	eventType  string
	client     *http.Client // Reuse HTTP client
}

const (
	url = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
)

func (m *Metric) SentMetrics(mv map[string]interface{}) error {
	req, err := http.NewRequest("POST", fmt.Sprintf(url, m.accID), nil)
	if err != nil {
		return fmt.Errorf("err create new req:%v", err)
	}

	if _, ok := mv["env"]; !ok {
		mv["env"] = m.env
	}

	mv["eventType"] = m.eventType

	// Set the authorization header.
	req.Header.Set("Api-Key", m.licenseKey)
	req.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(mv)
	if err != nil {
		return fmt.Errorf("err marshal json:%v", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(jsonBytes))
	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("err do req:%v", err)
	}
	defer resp.Body.Close()

	return nil
}

func NewMetric(accID, licenseKey, env, eventType string) (*Metric, error) {
	return &Metric{
		accID:      accID,
		licenseKey: licenseKey,
		env:        env,
		eventType:  eventType,
		client:     &http.Client{},
	}, nil
}
