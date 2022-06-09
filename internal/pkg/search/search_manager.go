package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	instance *Manager
)

func GetInstance() *Manager {
	return instance
}

type Manager struct {
	client      *elasticsearch.Client
	indexPrefix string
}

type Config struct {
	Url         string
	IndexPrefix string
}

func (manager *Manager) Setup(config Config) error {

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		fmt.Printf("elasticsearch new default client fail, %+v\n", err)
		return err
	}

	res, err := es.Info()
	if err != nil {
		fmt.Printf("elasticsearch get info fail, %+v\n", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		err = errors.New(res.String())
		fmt.Printf("elasticsearch get info res, is error, %+v\n", err)
		return err
	}

	instance = &Manager{client: es, indexPrefix: config.IndexPrefix}
	return nil
}

func (manager *Manager) CreateData(index string, id string, data interface{}) error {
	index = manager.indexPrefix + index

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(b),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), manager.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}

	return nil
}

func (manager *Manager) SearchData(index string, query map[string]interface{}, from int, size int) (int, []map[string]interface{}, error) {
	index = manager.indexPrefix + index

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return 0, nil, err
	}

	res, err := manager.client.Search(
		manager.client.Search.WithContext(context.Background()),
		manager.client.Search.WithIndex(index),
		manager.client.Search.WithBody(&buf),
		manager.client.Search.WithFrom(from),
		manager.client.Search.WithSize(size),
		manager.client.Search.WithTrackTotalHits(true),
		manager.client.Search.WithPretty(),
	)
	if err != nil {
		return 0, nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return 0, nil, errors.New(res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return 0, nil, err
	}

	repHits := r["hits"].(map[string]interface{})
	total := repHits["total"].(map[string]interface{})
	hits := repHits["hits"].([]interface{})
	count := int(total["value"].(float64))

	var sourceList []map[string]interface{}
	for _, hit := range hits {
		obj := hit.(map[string]interface{})
		source := obj["_source"].(map[string]interface{})
		sourceList = append(sourceList, source)
	}

	return count, sourceList, nil
}
