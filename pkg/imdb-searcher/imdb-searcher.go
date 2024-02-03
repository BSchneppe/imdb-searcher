package imdb_searcher

import (
	"fmt"
	"imdb-seeder/internal/pkg/search"
	meilisearch_client "imdb-seeder/internal/pkg/search/meilisearch"
	"regexp"
	"strconv"

	"go.uber.org/zap"

	"github.com/meilisearch/meilisearch-go"
	parsetorrentname "github.com/middelink/go-parse-torrent-name"
)

type ImdbSearchClient struct {
	index  *meilisearch.Index
	logger *zap.Logger
}
type ImdbMinimalTitle struct {
	Id    string
	Type  string
	Title string
}

type SearchClientConfig struct {
	MeiliSearchConfig meilisearch.ClientConfig
	Logger            *zap.Logger
}

func NewSearchClient(searchClientConfig SearchClientConfig) (*ImdbSearchClient, error) {
	if searchClientConfig.Logger == nil {
		searchClientConfig.Logger, _ = zap.NewProduction()
	}
	index, err := meilisearch_client.InitMeiliSearchClient(searchClientConfig.MeiliSearchConfig)
	if err != nil {
		return nil, err
	}
	searchClient := &ImdbSearchClient{index: index, logger: searchClientConfig.Logger}
	return searchClient, nil
}

var seasonRegex = regexp.MustCompile(`(\.|\s)S(\d+)`)

func (searchClient *ImdbSearchClient) GetClosestImdbTitle(fileName string) ImdbMinimalTitle {
	if len(fileName) < 2 {
		return ImdbMinimalTitle{}
	}
	parsedFile, err := parsetorrentname.Parse(fileName)
	if err != nil {
		return ImdbMinimalTitle{}
	}
	var imdbMinimal ImdbMinimalTitle
	if parsedFile.Title == "" {
		return ImdbMinimalTitle{}
	}
	titleSearch := parsedFile.Title
	var filters interface{}
	if parsedFile.Year != 0 {
		filters = "year < " + strconv.Itoa(parsedFile.Year+2) + " AND year > " + strconv.Itoa(parsedFile.Year-2)
	} else {
		filters = nil
	}
	if seasonRegex.MatchString(fileName) {
		if filters == nil {
			filters = "title_type = series"
		} else {
			filters = filters.(string) + " AND title_type = series"
		}

	}

	searchRes, err := searchClient.index.Search(search.NormalizeString(titleSearch),
		&meilisearch.SearchRequest{
			Limit:                1,
			AttributesToSearchOn: []string{"title"},
			Filter:               filters,
		})
	if err != nil {
		searchClient.logger.Error(fmt.Sprintf("Error searching for %s", titleSearch), zap.Error(err))
		return ImdbMinimalTitle{}
	}
	var hit map[string]interface{}
	for _, result := range searchRes.Hits {
		hit = result.(map[string]interface{})
		imdbMinimal.Id = hit["imdb_id"].(string)
		imdbMinimal.Type = hit["title_type"].(string)
		imdbMinimal.Title = hit["title"].(string)
		break
	}

	return imdbMinimal

}
