package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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

func ConstructPriceQuery(min, max int) interface{} {
	rangeQuery := map[string]interface{}{
		"range": map[string]interface{}{
			"final_price": map[string]interface{}{
				"gte": min,
				"lt":  max,
			},
		},
	}
	return rangeQuery
}

func ConstructTitleQuery(t string) interface{} {
	q := map[string]interface{}{
		"match": map[string]interface{}{
			"product_title": t,
		},
	}
	return q
}

func ConstructCategoryIDQuery(id int64) interface{} {
	q := map[string]interface{}{
		"match": map[string]interface{}{
			"category_id": id,
		},
	}
	return q
}

func (ec *ESClient) ScrollDocID(ctx context.Context, size int) ([]int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	// Serialize the combined query to JSON
	combinedJSON, err := json.Marshal(query)
	if err != nil {
		return []int{}, err
	}

	// Initial search request
	searchResp, err := ec.client.Search(
		ec.client.Search.WithIndex(ec.indexName),
		ec.client.Search.WithBody(bytes.NewReader(combinedJSON)),
		ec.client.Search.WithSize(size),
		ec.client.Search.WithSource("false"),
		ec.client.Search.WithScroll(time.Duration(1)*time.Minute),
	)
	if err != nil {
		return []int{}, err
	}
	defer searchResp.Body.Close()

	resArr := make([]int, 0)
	var result map[string]interface{}
	if err := json.NewDecoder(searchResp.Body).Decode(&result); err != nil {
		return []int{}, err
	}
	hits, _ := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		source, _ := hit.(map[string]interface{})
		docID, _ := source["_id"].(string)
		docIDInt, errConv := strconv.Atoi(docID)
		if errConv == nil {
			resArr = append(resArr, docIDInt)
		}

	}

	// Continue scrolling until there are no more results
	for {
		scrollResp, err := ec.client.Scroll(
			ec.client.Scroll.WithScrollID(result["_scroll_id"].(string)),
			ec.client.Scroll.WithScroll(time.Duration(1)*time.Minute), // Set the scroll duration
		)
		if err != nil {
			return []int{}, err
		}
		defer scrollResp.Body.Close()

		if scrollResp.IsError() {
			return []int{}, fmt.Errorf("scroll request failed with status code: %d", scrollResp.StatusCode)
		}

		// Process the scroll results
		if err := json.NewDecoder(scrollResp.Body).Decode(&result); err != nil {
			log.Fatalf("Error decoding scroll result: %s", err)
		}

		// Extract document IDs from the scroll results
		hits, _ := result["hits"].(map[string]interface{})["hits"].([]interface{})
		for _, hit := range hits {
			source, _ := hit.(map[string]interface{})
			docID, _ := source["_id"].(string)
			docIDInt, errConv := strconv.Atoi(docID)
			if errConv == nil {
				resArr = append(resArr, docIDInt)
			}
		}

		// Check if there are more results
		if len(hits) == 0 {
			break
		}
	}

	return resArr, nil
}

func (ec *ESClient) Query(ctx context.Context, q []interface{}, size int) ([]docgen.Document, error) {
	cq := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": q,
			},
		},
	}

	// Serialize the combined query to JSON
	combinedJSON, err := json.Marshal(cq)
	if err != nil {
		return []docgen.Document{}, err
	}

	res, err := ec.client.Search(
		ec.client.Search.WithContext(ctx),
		ec.client.Search.WithIndex(ec.indexName),
		ec.client.Search.WithBody(bytes.NewReader(combinedJSON)),
		ec.client.Search.WithTrackTotalHits(true),
		ec.client.Search.WithSize(size),
		ec.client.Search.WithSort("score_doc:asc"),
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
				ID:               int(source["id"].(float64)),
				ProductTitle:     source["product_title"].(string),
				ItemID:           int(source["item_id"].(float64)),
				CreatedUnix:      int64(source["created_unix"].(float64)),
				CreatedWeek:      source["created_week"].(string),
				CategoryID:       int64(source["category_id"].(float64)),
				CategoryName:     source["category_name"].(string),
				ProductCondition: int(source["product_condition"].(float64)),
				WidthImage:       int(source["width_image"].(float64)),
				HeightImage:      int(source["height_image"].(float64)),
				ImageURL:         source["image"].(string),
				Price:            source["price"].(float64),
				FinalPrice:       source["final_price"].(float64),
				DiscountPercent:  source["discount_percentage"].(float64),
				ScoreDoc:         source["score_doc"].(float64),
				CTR:              source["ctr"].(float64),
				PriceBid:         int(source["price_bid"].(float64)),
				IsCashback:       source["is_cashback"].(bool),
				IsFreeshipping:   source["is_freeshipping"].(bool),
				IsLocal:          source["is_local"].(bool),
				MinimumOrder:     int(source["minimum_order"].(float64)),
				ProvinceID:       int(source["province_id"].(float64)),
				ProvinceName:     source["province_name"].(string),
				CityID:           int(source["city_id"].(float64)),
				CityName:         source["city_name"].(string),
				DistrictID:       int(source["district_id"].(float64)),
				DistrictName:     source["district_name"].(string),
				Latitude:         source["lat"].(float64),
				Longitude:        source["long"].(float64),
				ShopName:         source["shop_name"].(string),
				ShopID:           int64(source["shop_id"].(float64)),
				ShopTier:         int(source["shop_tier"].(float64)),
				Rating:           source["rating"].(float64),
				Review:           int(source["review"].(float64)),
				WeightGram:       int(source["weight_gram"].(float64)),
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
