package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agasyan/es7-test/pkg/es"
	"github.com/agasyan/es7-test/pkg/fetchhandler"
	"github.com/agasyan/es7-test/pkg/metric"
	"gopkg.in/yaml.v3"
)

const (
	esVM   = "vm"
	esKube = "kube"
)

// config represents the structure of your configuration.
type config struct {
	Server serverConfig             `yaml:"server"`
	NR     nrConfig                 `yaml:"nr"`
	ES     map[string]esConfig      `yaml:"es"`
	H      map[string]handlerConfig `yaml:"handler"`
}

// serverConfig represents server configuration.
type serverConfig struct {
	AppName string `yaml:"app_name"`
	Port    int    `yaml:"port"`
	Env     string `yaml:"env"`
}

// nrConfig represents metric configuration.
type nrConfig struct {
	AccID        string `yaml:"acc_id"`
	LicenseKey   string `yaml:"license_key"`
	BufferSize   int    `yaml:"buffer_size"`
	BufferTimeMS int    `yaml:"buffer_time_ms"`
}

// esConfig represents es configuration.
type esConfig struct {
	Host      string `yaml:"host"`
	IndexName string `yaml:"index_name"`
}

type handlerConfig struct {
	TimeoutMS int `yaml:"timeout_ms"`
}

func main() {
	// Read the YAML file
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Parse the YAML data into the Config struct
	var config config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	// New Metric
	var m *metric.Metric
	if config.NR.AccID != "" && config.NR.LicenseKey != "" {
		m, err = metric.NewMetric(config.NR.AccID, config.NR.LicenseKey, config.Server.Env, config.Server.AppName, time.Duration(config.NR.BufferTimeMS)*time.Millisecond, config.NR.BufferSize)
		if err != nil {
			log.Fatalf("Error create metric service: %v", err)
		}
	}

	// NEW ES
	esm := make(map[string]*es.ESClient, len(config.ES))
	for kES, vES := range config.ES {
		if kES == esVM || kES == esKube {
			esc, err := es.NewElasticsearchClient([]string{vES.Host}, vES.IndexName)
			if err != nil {
				log.Fatalf("Error create ES service: %v", err)
			}
			esm[kES] = esc
		}
	}

	var hArr []string
	for k, v := range esm {
		if value, ok := config.H[k]; ok {
			h, err := fetchhandler.NewFetchHandler(
				v,
				m,
				k,
				time.Duration(value.TimeoutMS)*time.Millisecond)
			if err != nil {
				log.Fatalf("Error create handler service: %v", err)
			}
			http.HandleFunc(fmt.Sprintf("/get/%s", k), h.HandleRequest)
			hArr = append(hArr, k)
		}
	}

	if len(hArr) > 0 {
		fmt.Printf("Server is listening on :%d...", config.Server.Port)
		http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil)
	}
}
