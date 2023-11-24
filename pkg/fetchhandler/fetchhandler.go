package fetchhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/agasyan/es7-test/pkg/docgen"
	"github.com/agasyan/es7-test/pkg/es"
	"github.com/agasyan/es7-test/pkg/metric"
)

type FetchHandler struct {
	name    string
	es      *es.ESClient
	metric  *metric.Metric
	timeout time.Duration
}

func NewFetchHandler(es *es.ESClient, m *metric.Metric, name string, timeout time.Duration) (*FetchHandler, error) {
	return &FetchHandler{
		name:    name,
		timeout: timeout,
		metric:  m,
		es:      es,
	}, nil
}

type resp struct {
	Length int               `json:"length"`
	Docs   []docgen.Document `json:"docs,omitempty"`
}

// HandleRequest is a method of the FetchHandler that handles GET requests.
func (ph *FetchHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	var err error
	defer func(startTime time.Time) {
		if ph.metric != nil {
			ph.metric.SentMetrics(map[string]interface{}{
				"func":  fmt.Sprintf("FetchHandler.HandleRequest_%v", ph.name),
				"took":  time.Since(startTime).Seconds(),
				"isErr": strconv.FormatBool(err != nil),
			})
		}
	}(time.Now())

	// Parse the query parameters from the request
	queryParams := r.URL.Query()

	// Get specific query parameters
	width_min, _ := strconv.Atoi(queryParams.Get("width_min"))
	width_max, _ := strconv.Atoi(queryParams.Get("width_max"))
	height_min, _ := strconv.Atoi(queryParams.Get("height_min"))
	height_max, _ := strconv.Atoi(queryParams.Get("height_max"))
	title := queryParams.Get("title")
	genre := queryParams.Get("genre")
	author := queryParams.Get("author")

	var qArr []interface{}

	if width_min != 0 && width_max != 0 {
		qArr = append(qArr, es.ConstructWidthQuery(width_min, width_max))
	}

	if height_min != 0 && height_max != 0 {
		qArr = append(qArr, es.ConstructWidthQuery(height_min, height_max))
	}

	if title != "" {
		qArr = append(qArr, es.ConstructTitleQuery(title))
	}

	if genre != "" {
		qArr = append(qArr, es.ConstructGenreQuery(genre))
	}

	if author != "" {
		qArr = append(qArr, es.ConstructAuthorQuery(author))
	}
	size := 10
	szStr := queryParams.Get("size")
	if szStr != "" {
		szInt, _ := strconv.Atoi(szStr)
		if szInt != 0 {
			size = szInt
		}
	}

	var docs []docgen.Document
	ctx, cancel := context.WithTimeout(context.Background(), ph.timeout)
	defer cancel()
	docs, err = ph.es.Query(ctx, qArr, size)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal Server Error, err: %v", err), http.StatusInternalServerError)
		return
	}

	// You can perform further processing with the query parameters
	// For example, you might want to validate, sanitize, or use the parameters in your application logic.

	// Create a response struct
	response := resp{
		Docs:   docs,
		Length: len(docs),
	}

	// Marshal the response struct to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(responseJSON)
}
