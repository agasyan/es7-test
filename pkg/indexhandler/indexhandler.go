package indexhandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/agasyan/es7-test/pkg/docgen"
	"github.com/agasyan/es7-test/pkg/indexconsumer"
	"github.com/agasyan/es7-test/pkg/metric"
	"github.com/brianvoe/gofakeit/v6"
	nsq "github.com/nsqio/go-nsq"
)

type PostHandler struct {
	dg      *docgen.Docgen
	nsqProd *nsq.Producer
	topic   string
	metric  *metric.Metric
}

type postReq struct {
	Count int `json:"count"`
}

// HandleRequest is a method of the PostHandler that handles POST requests.
func (ph *PostHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func(startTime time.Time) {
		if ph.metric != nil {
			ph.metric.SentMetrics(map[string]interface{}{
				"func":  "PostHandler.HandleRequest",
				"took":  time.Since(startTime).Seconds(),
				"isErr": strconv.FormatBool(err != nil),
			})
		}
	}(time.Now())

	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		err = fmt.Errorf("mnthod not allowed")
		return
	}

	var pr postReq
	err = json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusInternalServerError)
		return
	}

	docs := ph.dg.BulkGenerate(pr.Count)
	var messageBodies [][]byte
	for _, d := range docs {
		action := indexconsumer.ActionIndex
		randAct := gofakeit.RandomString([]string{indexconsumer.ActionIndex, indexconsumer.ActionUpdate, indexconsumer.ActionDelete})
		// to minimize update delete
		if randAct == indexconsumer.ActionDelete || randAct == indexconsumer.ActionUpdate {
			idToBeUpdated := ph.dg.GetExistKey(randAct)
			if idToBeUpdated > 0 {
				action = randAct
				d.ID = idToBeUpdated
			}
		}

		var mb []byte
		mb, err = json.Marshal(indexconsumer.IndexerHandlerMSG{
			Action: action,
			Doc:    d,
		})

		if err != nil {
			log.Printf("ERR marshall docs: %v", err)
			http.Error(w, "Error Marshall docs", http.StatusInternalServerError)
			return
		}
		messageBodies = append(messageBodies, mb)
	}

	err = ph.nsqProd.MultiPublish(ph.topic, messageBodies)
	if err != nil {
		http.Error(w, "Error sent nsq docs", http.StatusInternalServerError)
		return
	}

	// update the array after publish
	ph.dg.UpdateArr()

	// Send a response back

	fmt.Fprintf(w, "POST request received successfully!, with count doc %v", pr.Count)
}

func NewPostHandler(dg *docgen.Docgen, nsqdAddr, topic string, m *metric.Metric) (*PostHandler, error) {
	producer, err := nsq.NewProducer(nsqdAddr, nsq.NewConfig())
	if err != nil {
		return nil, fmt.Errorf("ERR create producer %v", err)
	}

	// Check if the producer is connected to NSQ
	if err := producer.Ping(); err != nil {
		return nil, fmt.Errorf("ERR ping producer %v", err)
	}

	return &PostHandler{
		dg:      dg,
		nsqProd: producer,
		metric:  m,
		topic:   topic,
	}, nil
}
