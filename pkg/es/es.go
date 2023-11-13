package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/agasyan/es7-test/pkg/docgen"
	es7 "github.com/elastic/go-elasticsearch/v7"
)

// ESClient represents an Elasticsearch client object.
type ESClient struct {
	client    *es7.Client
	indexName string
}

// NewElasticsearchClient creates a new instance of ElasticsearchClient.
func NewElasticsearchClient(addresses []string, indexName string) (*ESClient, error) {
	cfg := es7.Config{
		Addresses: addresses,
	}

	client, err := es7.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ESClient{client: client, indexName: indexName}, nil
}

// Index indexes a document in Elasticsearch.
func (ec *ESClient) Index(ctx context.Context, document docgen.Document) error {
	// Convert the struct to JSON
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return err
	}

	// Prepare the request with the provided context
	res, err := ec.client.Index(ec.indexName, strings.NewReader(string(documentJSON)),
		ec.client.Index.WithDocumentID(strconv.Itoa(document.ID)),
		ec.client.Index.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response status
	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Print the response
	fmt.Printf("Indexing Response for index %s, id %d: %s\n", ec.indexName, document.ID, body)

	return nil
}

// Update updates a document in Elasticsearch.
func (ec *ESClient) Update(ctx context.Context, document docgen.Document) error {
	// Convert the struct to JSON
	documentJSON, err := json.Marshal(map[string]interface{}{
		"doc": document, // Use "doc" field to update specific fields
	})
	if err != nil {
		return err
	}

	// Prepare the request with the provided context
	res, err := ec.client.Update(
		ec.indexName,
		strconv.Itoa(document.ID),
		strings.NewReader(string(documentJSON)),
		ec.client.Update.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response status
	if res.IsError() {
		return fmt.Errorf("error updating document: %s", res.String())
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Print the response
	fmt.Printf("Updating Response for index %s, id %d: %s\n", ec.indexName, document.ID, body)

	return nil
}

// Delete deletes a document in Elasticsearch.
func (ec *ESClient) Delete(ctx context.Context, document docgen.Document) error {
	// Prepare the request with the provided context
	res, err := ec.client.Delete(
		ec.indexName,
		strconv.Itoa(document.ID),
		ec.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check the response status
	if res.IsError() {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Print the response
	fmt.Printf("Deleting Response for index %s, id %d: %s\n", ec.indexName, document.ID, body)

	return nil
}

func ConstructWidthQuery(min, max int) interface{} {
	rangeQuery := map[string]interface{}{
		"range": map[string]interface{}{
			"width_image": map[string]interface{}{
				"gte": min,
				"lt":  max,
			},
		},
	}
	return rangeQuery
}

func ConstructHeightQuery(min, max int) interface{} {
	rangeQuery := map[string]interface{}{
		"range": map[string]interface{}{
			"height_image": map[string]interface{}{
				"gte": min,
				"lt":  max,
			},
		},
	}
	return rangeQuery
}

func ConstructAuthorQuery(name string) interface{} {
	q := map[string]interface{}{
		"match": map[string]interface{}{
			"author": name,
		},
	}
	return q
}

func ConstructTitleQuery(t string) interface{} {
	q := map[string]interface{}{
		"match": map[string]interface{}{
			"title": t,
		},
	}
	return q
}

func ConstructGenreQuery(t string) interface{} {
	q := map[string]interface{}{
		"match": map[string]interface{}{
			"genre": t,
		},
	}
	return q
}

func (ec *ESClient) Query(ctx context.Context, q []interface{}, size int) ([]docgen.Document, error) {
	boolQuery := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": q,
		},
	}
	// Serialize the combined query to JSON
	combinedJSON, err := json.Marshal(boolQuery)
	if err != nil {
		return []docgen.Document{}, err
	}

	res, err := ec.client.Search(
		ec.client.Search.WithContext(ctx),
		ec.client.Search.WithIndex(ec.indexName),
		ec.client.Search.WithBody(bytes.NewReader(combinedJSON)),
		ec.client.Search.WithTrackTotalHits(true),
		ec.client.Search.WithSize(size),
		ec.client.Search.WithSort("price:asc"),
	)
	if err != nil {
		return []docgen.Document{}, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return []docgen.Document{}, err
	}

	var documents []docgen.Document
	hits, found := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if found {
		for _, hit := range hits {
			source, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})
			document := docgen.Document{
				ID:          int(source["id"].(float64)),
				Title:       source["title"].(string),
				Author:      source["author"].(string),
				Genre:       source["genre"].(string),
				WidthImage:  int(source["width_image"].(float64)),
				HeightImage: int(source["height_image"].(float64)),
				ImageURL:    source["image_url"].(string),
				CreatedUnix: int64(source["created_unix"].(float64)),
			}
			documents = append(documents, document)
		}
	}

	return documents, nil
}

// GetInfo retrieves and prints the Elasticsearch server information.
func (ec *ESClient) GetInfo() error {
	info, err := ec.client.Info()
	if err != nil {
		return err
	}

	fmt.Printf("Elasticsearch Info: %+v\n", info)
	return nil
}
