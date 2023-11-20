package main

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	vegeta "github.com/tsenart/vegeta/lib"
)

func NewCustomTargeter(baseURL string) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"
		tgt.URL = baseURL

		// Create a URL object
		u, _ := url.Parse(baseURL)

		// Add query parameters
		q := u.Query()
		randQ := gofakeit.RandomString([]string{"genre", "title", "author", "width", "height"})
		switch randQ {
		case "genre":
			q.Add("genre", gofakeit.BookGenre())
		case "title":
			q.Add("title", gofakeit.BookTitle())
		case "author":
			q.Add("author", gofakeit.FirstName())
		case "width":
			w_min := gofakeit.Number(200, 600)
			w_max := gofakeit.Number(600, 2000)
			q.Add("width_min", fmt.Sprint(w_min))
			q.Add("width_max", fmt.Sprint(w_max))
		case "height":
			min := gofakeit.Number(200, 600)
			max := gofakeit.Number(600, 2000)
			q.Add("height_min", fmt.Sprint(min))
			q.Add("height_max", fmt.Sprint(max))
		}

		size := gofakeit.Number(50, 70)
		q.Add("size", fmt.Sprint(size))

		// Update the URL object with the new query parameters
		u.RawQuery = q.Encode()
		tgt.URL = u.String()

		return nil
	}
}

func runAttack(targeter vegeta.Targeter, rate vegeta.Rate, duration time.Duration, name string, wg *sync.WaitGroup) {
	defer wg.Done()

	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, name) {
		metrics.Add(res)
	}
	metrics.Close()

	rep, _ := vegeta.NewTextReporter(&metrics), &metrics

	// Create a buffer to capture the output
	var buffer bytes.Buffer
	rep.Report(&buffer)

	// Convert the buffer's contents to a string
	result := buffer.String()

	fmt.Printf("%s Report:\n%s\n", name, result)
}

func main() {
	rate := vegeta.Rate{Freq: 5, Per: 1 * time.Second}
	duration := 120 * time.Second

	// Define your old and new URLs
	oldURL := "http://localhost:8082/get/vm"
	newURL := "http://localhost:8082/get/vm"

	// Create wait group to wait for both attacks to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Run attack for the old URL concurrently
	go runAttack(NewCustomTargeter(oldURL), rate, duration, "Old URL", &wg)

	// Run attack for the new URL concurrently
	go runAttack(NewCustomTargeter(newURL), rate, duration, "New URL", &wg)

	// Wait for both attacks to finish
	wg.Wait()
}
