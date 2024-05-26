package indexconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/agasyan/es7-test/pkg/docgen"
	"github.com/agasyan/es7-test/pkg/es"
	"github.com/agasyan/es7-test/pkg/metric"
	"github.com/nsqio/go-nsq"
)

// MessageHandler is a custom NSQ message handler.
type IndexerHandler struct {
	m       *metric.Metric
	es      *es.ESClient
	timeout time.Duration
	C       *nsq.Consumer
	name    string
}

type IndexerHandlerMSG struct {
	Action string          `json:"action"`
	Doc    docgen.Document `json:"doc"`
}

const (
	ActionIndex  = "INDEX"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

// HandleMessage is called when a new message is received.
func (h *IndexerHandler) HandleMessage(msg *nsq.Message) error {
	// Process the received message.
	var err error
	var act string
	found := true
	defer func(startTime time.Time) {
		if h.m != nil {
			h.m.SentMetrics(map[string]interface{}{
				"func":   fmt.Sprintf("IndexerHandler.HandleMessage_%v", h.name),
				"took":   time.Since(startTime).Seconds(),
				"isErr":  err != nil,
				"action": act,
				"found":  found,
			})
		}
	}(time.Now())

	var message IndexerHandlerMSG
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Println("Error unmarshaling message:", err)
		msg.Finish()
		return err
	}
	act = message.Action

	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	switch message.Action {
	case ActionIndex:
		err = h.es.Index(ctx, message.Doc)
	case ActionUpdate:
		err = h.es.Update(ctx, message.Doc)
	case ActionDelete:
		err = h.es.Delete(ctx, message.Doc)
	}
	if err != nil && (strings.Contains(strings.ToLower(err.Error()), "not found") || strings.Contains(strings.ToLower(err.Error()), "conflict")) {
		err = nil
		found = false
	}
	if err != nil {
		log.Println("Error es doc:", err)
		msg.RequeueWithoutBackoff(-1)
		return err
	}

	msg.Finish()
	return nil
}

func NewIndexerHandler(m *metric.Metric, esc *es.ESClient, timeout time.Duration, topicName, channelName, nsqd, keyNSQ string, numOfConsumer, maxInFlight int) (*IndexerHandler, error) {
	config := nsq.NewConfig()
	if maxInFlight > 0 {
		config.MaxInFlight = maxInFlight
	}
	consumer, err := nsq.NewConsumer(topicName, channelName, config)
	if err != nil {
		log.Println("Error creating NSQ consumer:", err)
		return nil, err
	}
	ih := &IndexerHandler{
		m:       m,
		es:      esc,
		C:       consumer,
		timeout: timeout,
		name:    keyNSQ,
	}
	consumer.AddConcurrentHandlers(ih, numOfConsumer)
	err = consumer.ConnectToNSQD(nsqd)
	if err != nil {
		log.Println("Error connecting to NSQD:", err)
		return nil, err
	}
	return ih, nil
}
