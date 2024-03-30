package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/agasyan/es7-test/pkg/docgen"
	"github.com/brianvoe/gofakeit/v6"
	vegeta "github.com/tsenart/vegeta/lib"
)

var (
	i      = 0
	j      = 0
	mutexI = &sync.Mutex{}
	mutexJ = &sync.Mutex{}
	a1     []string
	a2     []string
)

func NewCustomTargeter(baseURL string) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"
		mutexI.Lock()
		tgt.URL = a1[i]
		i++
		mutexI.Unlock()

		return nil
	}
}

func NewCustomTargeter2(baseURL string) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = "GET"

		mutexJ.Lock()
		tgt.URL = a2[j]
		j++
		mutexJ.Unlock()

		return nil
	}
}

func generateRandomUrlArr(count int, baseURL1 string, baseURL2 string, d *docgen.Docgen) ([]string, []string) {
	arrUrls1 := make([]string, 0, count)
	arrUrls2 := make([]string, 0, count)
	for i := 0; i < count; i++ {
		randQ := gofakeit.RandomString([]string{"price", "title", "catid", "width", "height"})
		// Create a URL object
		u1, _ := url.Parse(baseURL1)
		// Create a URL object
		u2, _ := url.Parse(baseURL2)

		q := u1.Query()
		switch randQ {
		case "price":
			fp_min := gofakeit.Number(1000, 10000)
			fp_max := gofakeit.Number(5000000, 10000000)
			q.Add("fp_min", fmt.Sprint(fp_min))
			q.Add("fp_max", fmt.Sprint(fp_max))
		case "title":
			n := d.GetRandomProduct().Name
			ptArr := strings.Split(n, " ")
			t := gofakeit.RandomString(ptArr)
			q.Add("title", strings.ToLower(t))
		case "cat_id":
			a := d.GetRandomProduct().CategoriesID
			q.Add("cat_id", fmt.Sprint(a))
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

		size := gofakeit.Number(20, 70)
		q.Add("size", fmt.Sprint(size))

		u1.RawQuery = q.Encode()
		u2.RawQuery = q.Encode()

		arrUrls1 = append(arrUrls1, u1.String())
		arrUrls2 = append(arrUrls2, u2.String())
	}

	return arrUrls1, arrUrls2
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
	// lt 75 bau 30
	f := 30
	timeSec := 1
	durSec := 600
	rate := vegeta.Rate{Freq: f, Per: time.Duration(timeSec) * time.Second}
	duration := time.Duration(durSec) * time.Second

	// Define your old and new URLs
	u1 := "http://localhost:8082/get/vm"
	u2 := "http://localhost:8082/get/kube"

	log.Println("generating random url")
	d := docgen.Init()
	a1, a2 = generateRandomUrlArr((f * timeSec * durSec), u1, u2, d)
	log.Printf("done generating random url len :%v ", len(a1))

	// Create wait group to wait for both attacks to finish
	var wg sync.WaitGroup
	wg.Add(2)

	// Run attack for the old URL concurrently
	go runAttack(NewCustomTargeter(u1), rate, duration, "URL 1 [VM]", &wg)

	// Run attack for the new URL concurrently
	go runAttack(NewCustomTargeter2(u2), rate, duration, "URL 2 [KUBE]", &wg)

	// Wait for both attacks to finish
	wg.Wait()

	log.Println("done attack")
}
