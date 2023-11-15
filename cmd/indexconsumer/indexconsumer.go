package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agasyan/es7-test/pkg/es"
	"github.com/agasyan/es7-test/pkg/indexconsumer"
	"github.com/agasyan/es7-test/pkg/metric"
	"gopkg.in/yaml.v3"
)

const (
	esKube   = "kube"
	esVM     = "vm"
	consKube = "ckube"
	consVM   = "cvm"
)

// config represents the structure of your configuration.
type config struct {
	Server    serverConfig         `yaml:"server"`
	NR        nrConfig             `yaml:"nr"`
	ES        map[string]esConfig  `yaml:"es"`
	NSQConfig map[string]nsqConfig `yaml:"nsq"`
}

// serverConfig represents server configuration.
type serverConfig struct {
	AppName string `yaml:"app_name"`
	Port    int    `yaml:"port"`
	Env     string `yaml:"env"`
}

// nrConfig represents metric configuration.
type nrConfig struct {
	AccID      string `yaml:"acc_id"`
	LicenseKey string `yaml:"license_key"`
}

// nsqConfig represents nsq configuration.
type nsqConfig struct {
	NSQDAddress  string `yaml:"nsqd_address"`
	PublishTopic string `yaml:"publish_topic"`
	ConsumerName string `yaml:"consumer_name"`
	TimeoutMS    int    `yaml:"timeout_ms"`
}

// esConfig represents es configuration.
type esConfig struct {
	Host      string `yaml:"host"`
	IndexName string `yaml:"index_name"`
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
		m, err = metric.NewMetric(config.NR.AccID, config.NR.LicenseKey, config.Server.Env, config.Server.AppName)
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

	// New Consumer
	for kNSQ, vNSQ := range config.NSQConfig {

		var esc *es.ESClient
		switch kNSQ {
		case consVM:
			esc = esm[esVM]
		case consKube:
			esc = esm[esKube]
		}

		if esc != nil {
			ih, err := indexconsumer.NewIndexerHandler(
				m,
				esc,
				time.Duration(vNSQ.TimeoutMS)*time.Millisecond,
				vNSQ.PublishTopic,
				vNSQ.ConsumerName,
				vNSQ.NSQDAddress,
			)
			if err != nil {
				log.Fatalf("Error creating consumer: %v", err)
			}
			defer ih.C.Stop()
		}

	}

	http.HandleFunc("/health", handlerHealthCheck)
	fmt.Printf("Server is listening on :%d...", config.Server.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil)

}

func handlerHealthCheck(w http.ResponseWriter, _ *http.Request) {
	_, _ = io.WriteString(w, "Success")
}
