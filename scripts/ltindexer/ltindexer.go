package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	vegeta "github.com/tsenart/vegeta/lib"
)

// Define a struct that matches the JSON structure
type myData struct {
	Count   int  `json:"count"`
	IsIndex bool `json:"is_index"`
}

func NewCustomTargeter(IsIndexOnly bool) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		data := myData{
			Count:   gofakeit.Number(1, 5),
			IsIndex: IsIndexOnly,
		}
		jsonData, _ := json.Marshal(data)

		tgt.Method = "POST"
		tgt.Body = jsonData
		tgt.URL = "http://localhost:8080/random-action-es"

		return nil
	}
}

func main() {
	// lt 40 bau 20
	rate := vegeta.Rate{Freq: 40, Per: 1 * time.Second}
	IsIndexOnly := true
	duration := 600 * time.Second
	targeter := NewCustomTargeter(IsIndexOnly)
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
