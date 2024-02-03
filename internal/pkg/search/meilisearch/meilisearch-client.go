package meilisearch_client

import "github.com/meilisearch/meilisearch-go"

func InitMeiliSearchClient(clientConfig meilisearch.ClientConfig) (*meilisearch.Index, error) {
	client := meilisearch.NewClient(clientConfig)
	index := client.Index("imdb")
	_, err := index.UpdateFilterableAttributes(&[]string{"title_type", "year"})
	if err != nil {
		return nil, err
	}

	return index, nil
}
