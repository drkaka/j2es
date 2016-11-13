package uploader

import (
	"context"
	"errors"

	es "gopkg.in/olivere/elastic.v5"
)

var (
	client *es.BulkService
)

// PrepareUploader the upload client.
func PrepareUploader(urls []string) error {
	c, err := es.NewClient(
		es.SetURL(urls...),
		es.SetMaxRetries(10),
	)

	if err != nil {
		return errors.New("Failed to create ES client.")
	}

	client = es.NewBulkService(c)
	client.Timeout("60s")

	return nil
}

// AddRecord to add one record.
func AddRecord(index string, record interface{}) {
	req := es.NewBulkIndexRequest().Doc(record)
	req.Index(index)
	req.Type("log")
	client.Add(req)
}

// PostToES to post records to Elastic Search.
func PostToES() error {
	if _, err := client.Do(context.Background()); err != nil {
		return err
	}
	return nil
}
