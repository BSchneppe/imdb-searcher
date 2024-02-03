package imdb_seeder

import (
	"compress/gzip"
	"crypto/tls"
	"html"
	meilisearch_client "imdb-seeder/internal/pkg/search/meilisearch"
	"imdb-seeder/internal/pkg/tsv_reader"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
)

var client = &http.Client{Transport: &http.Transport{
	TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	DisableCompression: false,
	DisableKeepAlives:  true,
	IdleConnTimeout:    20 * time.Second}}

func Seed(meilisearchConfig meilisearch.ClientConfig, logger *zap.Logger) {
	index, err := meilisearch_client.InitMeiliSearchClient(meilisearchConfig)
	if err != nil {
		logger.Fatal("failed to initialize meilisearch client", zap.Error(err))
	}

	taskIds := []int64{}

	titlemap := map[string]string{
		"movie":    "movie",
		"tvMovie":  "movie",
		"tvSeries": "series",
	}

	logger.Info("Writing titles to meilisearch...")

	req, _ := http.NewRequest("GET", "https://datasets.imdbws.com/title.basics.tsv.gz", nil)

	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal("failed to get data from imdb", zap.Error(err))
	}
	defer resp.Body.Close()

	gzreader, err := gzip.NewReader(resp.Body)
	if err != nil {
		logger.Fatal("failed to read gzip data", zap.Error(err))
	}
	defer gzreader.Close()
	parsertitle := tsv_reader.NewTabNewlineReader(gzreader)
	_, _ = parsertitle.Read()

	valueArgs := make([]map[string]interface{}, 0)
	var record []string
	var tsverr error
	var ok bool
	rowCount, insertCount := 0, 0
	for {
		record, tsverr = parsertitle.Read()
		if tsverr == io.EOF {
			break
		}
		if tsverr != nil {
			logger.Error("failed to read tsv record", zap.Error(tsverr))
			continue
		}
		typeOfMedia := record[1]
		if _, ok = titlemap[typeOfMedia]; ok && typeOfMedia != "" {
			rowCount++
			id := record[0]
			idWithoutPrefix := csvgetint(strings.TrimLeft(strings.TrimPrefix(id, "tt"), "0"))
			title := html.UnescapeString(record[2])
			year := csvgetint(record[5])
			imdbRecord := map[string]interface{}{
				"id":         idWithoutPrefix,
				"imdb_id":    id,
				"title":      title,
				"year":       year,
				"title_type": titlemap[typeOfMedia],
			}
			insertCount++
			valueArgs = append(valueArgs, imdbRecord)
		}
		if len(valueArgs) > 9998 {
			taskInfo, err := index.AddDocuments(valueArgs, "id")
			if err != nil {
				logger.Error("failed to add documents to meilisearch", zap.Error(err))
			}
			valueArgs = make([]map[string]interface{}, 0)
			taskIds = append(taskIds, taskInfo.TaskUID)
		}
	}

	if len(valueArgs) > 1 {
		taskInfo, err := index.AddDocuments(valueArgs, "id")
		if err != nil {
			logger.Error("failed to add documents to meilisearch", zap.Error(err))
		}
		taskIds = append(taskIds, taskInfo.TaskUID)
	}
	for _, id := range taskIds {
		task, _ := index.WaitForTask(id)
		if task.Status != "succeeded" {
			logger.Error("failed to add documents to meilisearch", zap.Any("task", task))
		}
	}
	logger.Info("Processed "+strconv.Itoa(rowCount+1)+" titles", zap.Int("inserted", insertCount))
}
func csvgetint(instr string) int {
	getint, err := strconv.Atoi(instr)
	if err != nil {
		return 0
	}
	return getint
}
