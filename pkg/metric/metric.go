package metric

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Metric struct {
	licenseKey     string
	accID          string
	env            string
	eventType      string
	client         *http.Client // Reuse HTTP client
	metricBuffer   chan map[string]interface{}
	bufferMutex    sync.Mutex
	bufferInterval time.Duration
	bufferSize     int
}

const (
	url             = "https://insights-collector.newrelic.com/v1/accounts/%s/events"
	defaultInterval = 1 * time.Second
	defaultSize     = 100
)

func (m *Metric) SentMetrics(mv map[string]interface{}) {
	if _, ok := mv["env"]; !ok {
		mv["env"] = m.env
	}

	mv["eventType"] = m.eventType

	m.bufferMutex.Lock()

	if m.metricBuffer == nil {
		m.metricBuffer = make(chan map[string]interface{}, m.bufferSize)
		go m.sendMetricsPeriodically()
	}

	clonedMV := make(map[string]interface{})
	for k, v := range mv {
		clonedMV[k] = v
	}

	// If the buffer size is reached, trigger an immediate send
	if len(m.metricBuffer) >= m.bufferSize {
		// Unlock before triggering the immediate send
		m.bufferMutex.Unlock()
		go m.sendMetricsNow()
	} else {
		m.metricBuffer <- clonedMV
		m.bufferMutex.Unlock()
	}
}

func (m *Metric) sendMetricsNow() {
	m.bufferMutex.Lock()

	// If there are metrics in the buffer, send them
	if len(m.metricBuffer) > 0 {
		var metrics []map[string]interface{}
		for i := 0; i < len(m.metricBuffer); i++ {
			metrics = append(metrics, <-m.metricBuffer)
		}

		m.bufferMutex.Unlock() // Move the unlock here

		m.sendBatchMetrics(metrics)
	} else {
		m.bufferMutex.Unlock()
	}
}

func (m *Metric) sendMetricsPeriodically() {
	ticker := time.NewTicker(m.bufferInterval)
	defer ticker.Stop()

	for range ticker.C {
		go m.sendMetricsNow()
	}
}

func (m *Metric) sendBatchMetrics(metrics []map[string]interface{}) {
	for retry := 0; retry < 3; retry++ {
		req, err := http.NewRequest("POST", fmt.Sprintf(url, m.accID), nil)
		if err != nil {
			log.Printf("err create new req:%v\n", err)
			return
		}

		req.Header.Set("Api-Key", m.licenseKey)
		req.Header.Set("Content-Type", "application/json")

		jsonBytes, err := json.Marshal(metrics)
		if err != nil {
			log.Printf("err marshal json:%v\n", err)
			return
		}

		req.Body = io.NopCloser(bytes.NewReader(jsonBytes))
		resp, err := m.client.Do(req)
		if err != nil {
			log.Printf("err do req:%v, retry:%v\n", err, retry)

			// Retry if it's not the last attempt
			if retry < 2 {
				log.Println("Retrying...")
				continue
			}

			return
		}
		defer resp.Body.Close()

		break
	}
}

func NewMetric(accID, licenseKey, env, eventType string, bufferTime time.Duration, bufferSize int) (*Metric, error) {
	bi := defaultInterval
	if bufferTime > 0 {
		bi = bufferTime
	}

	bs := defaultSize
	if bufferSize > 0 {
		bs = bufferSize
	}

	return &Metric{
		accID:          accID,
		licenseKey:     licenseKey,
		env:            env,
		eventType:      eventType,
		client:         &http.Client{},
		bufferInterval: bi,
		bufferSize:     bs,
	}, nil
}
