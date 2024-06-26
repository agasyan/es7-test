package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

// Define a struct that matches the JSON structure
type myData struct {
	Count   int  `json:"count"`
	IsIndex bool `json:"is_index"`
}

func NewCustomTargeter(IsIndexOnly bool, count int) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		data := myData{
			Count:   count,
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
	f := 50
	timeSec := 1
	durSec := 600
	rate := vegeta.Rate{Freq: f, Per: time.Duration(timeSec) * time.Second}
	IsIndexOnly := false
	count := 1
	duration := time.Duration(durSec) * time.Second
	targeter := NewCustomTargeter(IsIndexOnly, count)
	attacker := vegeta.NewAttacker()
	log.Printf("Starting Load Test for, %v, with freq:%v per timeInSec:%v\n", duration.Seconds(), f, timeSec)
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
