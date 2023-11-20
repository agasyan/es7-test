package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/agasyan/es7-test/pkg/docgen"
	"github.com/agasyan/es7-test/pkg/es"
	"github.com/agasyan/es7-test/pkg/indexhandler"
	"github.com/agasyan/es7-test/pkg/metric"
	"gopkg.in/yaml.v3"
)

// config represents the structure of your configuration.
type config struct {
	Server    serverConfig        `yaml:"server"`
	NR        nrConfig            `yaml:"nr"`
	NSQConfig nsqConfig           `yaml:"nsq"`
	ES        map[string]esConfig `yaml:"es"`
}

// esConfig represents es configuration.
type esConfig struct {
	Host      string `yaml:"host"`
	IndexName string `yaml:"index_name"`
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
}

const (
	esKube = "kube"
	esVM   = "vm"
)

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
	m, err := metric.NewMetric(config.NR.AccID, config.NR.LicenseKey, config.Server.Env, config.Server.AppName)
	if err != nil {
		log.Fatalf("Error create metric service: %v", err)
	}

	// New Docgen
	d := docgen.Init()
	h, err := indexhandler.NewPostHandler(d, config.NSQConfig.NSQDAddress, config.NSQConfig.PublishTopic, m)
	if err != nil {
		log.Fatalf("Error create handler service: %v", err)
	}

	// Make map
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

	// init map and arr doc id
	if len(esm) > 0 {
		for _, v := range esm {
			if v != nil {
				ids, err := v.ScrollDocID(context.Background(), 100)
				if err != nil {
					log.Fatalf("Error scroll ES: %v", err)
				} else {
					d.InitMapArr(ids)
					// one time only
					break
				}

			}

		}
	}

	http.HandleFunc("/random-action-es", h.HandleRequest)
	fmt.Printf("Server is listening on :%d...", config.Server.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil)

}
