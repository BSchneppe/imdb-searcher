package imdb_searcher

import (
	"context"
	meilisearch_client "imdb-seeder/internal/pkg/search/meilisearch"
	"log"
	"os"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var searchClient *ImdbSearchClient

func TestMain(m *testing.M) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "getmeili/meilisearch:latest",
		ExposedPorts: []string{"7700/tcp"},
		Env:          map[string]string{"MEILI_ENV": "development"},
		WaitingFor:   wait.ForExposedPort(),
	}
	meiliC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start meilisearch: %s", err)
	}

	endpoint, err := meiliC.Endpoint(ctx, "http")
	if err != nil {
		log.Fatalf("Could not get endpoint: %s", err)
	}

	clientConfig := meilisearch.ClientConfig{
		Host: endpoint,
	}
	client, err := meilisearch_client.InitMeiliSearchClient(clientConfig)
	if err != nil {
		log.Fatalf("could not connect to meilisearch: %s", err)
	}
	res, err := client.AddDocuments([]map[string]interface{}{
		{"id": 38650, "imdb_id": "tt0038650", "title": "It's a Wonderful Life", "year": 1946, "title_type": "movie"},
		{"id": 63350, "imdb_id": "tt0063350", "title": "Night Of The Living Dead", "year": 1968, "title_type": "movie"},
		{"id": 140738, "imdb_id": "tt0140738", "title": "Flash Gordon", "year": 1954, "title_type": "series"},
		{"id": 123, "imdb_id": "tt0123", "title": "Flash Gordon", "year": 1953, "title_type": "series"},
		{"id": 124, "imdb_id": "tt01234", "title": "Flash Gordon", "year": 1954, "title_type": "movie"},
	}, "id")
	if err != nil {
		log.Fatalf("could not add documents: %s", err)
	}
	task, err := client.WaitForTask(res.TaskUID)
	if err != nil || task.Status != "succeeded" {
		log.Fatalf("could not add documents: %s,%v", err, task)
	}

	searchClient, err = NewSearchClient(SearchClientConfig{MeiliSearchConfig: clientConfig})
	if err != nil {
		log.Fatalf("could not create search client: %s", err)
	}

	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	if err := meiliC.Terminate(ctx); err != nil {
		log.Fatalf("Could not stop meilisearch: %s", err)
	}

	os.Exit(code)
}

func TestImdbSearchClient_GetClosestImdbTitleMovie(t *testing.T) {
	imdbTitle := searchClient.GetClosestImdbTitle("Night Of The Living Dead (1968) [BluRay] [720p] [YTS.AM]/Night.Of.The.Living.Dead.1968.720p.BluRay.x264-[YTS.AM].mp4")
	if imdbTitle.Title != "Night Of The Living Dead" {
		t.Errorf("expected title to be Night Of The Living Dead, got %s", imdbTitle.Title)
	}
}

func TestImdbSearchClient_GetClosestImdbTitleEmptyString(t *testing.T) {
	imdbTitle := searchClient.GetClosestImdbTitle("")
	if imdbTitle.Title != "" {
		t.Errorf("expected title to be empty, got %s", imdbTitle.Title)
	}
}

func TestImdbSearchClient_GetClosestImdbTitleSeries(t *testing.T) {
	imdbTitle := searchClient.GetClosestImdbTitle("FlashGordon.S01E01.ThePlanetOfDeath_512kb.mp4")
	if imdbTitle.Title != "Flash Gordon" {
		t.Errorf("expected title to be Flash Gordon, got %s", imdbTitle.Title)
	}
}

func TestImdbSearchClient_GetClosestImdbTitleSeriesWithYear(t *testing.T) {
	imdbTitle := searchClient.GetClosestImdbTitle("FlashGordon.S01E01.1954.ThePlanetOfDeath_512kb.mp4")
	if imdbTitle.Title != "Flash Gordon" || imdbTitle.Id != "tt0140738" {
		t.Errorf("expected title to be Flash Gordon, got %s ,Id:%s", imdbTitle.Title, imdbTitle.Id)
	}
}
