package main

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	vegeta "github.com/tsenart/vegeta/lib"
)

func NewCustomTargeter() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"

		baseURL := "http://localhost:8082/get/vm"

		// Create a URL object
		u, _ := url.Parse(baseURL)

		// Add query parameters
		q := u.Query()
		q.Add("genre", gofakeit.BookGenre())
		// Update the URL object with the new query parameters
		u.RawQuery = q.Encode()
		tgt.URL = u.String()

		return nil
	}
}

func main() {
	rate := vegeta.Rate{Freq: 2, Per: 1 * time.Second}
	duration := 3 * time.Second
	targeter := NewCustomTargeter()
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Load Test") {
		metrics.Add(res)
	}
	metrics.Close()

	rep, _ := vegeta.NewTextReporter(&metrics), &metrics

	// Create a buffer to capture the output
	var buffer bytes.Buffer

	rep.Report(&buffer)

	// Convert the buffer's contents to a string
	result := buffer.String()

	fmt.Printf("%+v  \n", result)
}
