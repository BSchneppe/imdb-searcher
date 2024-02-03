# IMDb Searcher with MeiliSearch

This is a simple IMDb searcher that uses MeiliSearch to index and search IMDb data. The IMDb data is seeded from the
[IMDb dataset](https://developer.imdb.com/non-commercial-datasets/).

## Prerequisites

Before you can use the IMDb Seeder and Searcher Client, you must have the following prerequisites installed and
configured:

- Go: The client is written in Go. You need Go 1.21 or later installed on your machine to compile and run the Go
  programs. Download it from the [Go official site](https://go.dev/).
- MeiliSearch: Although the client automatically sets up MeiliSearch using Docker, having MeiliSearch installed locally
  for development purposes can be helpful. Visit the [MeiliSearch documentation](https://www.meilisearch.com/docs/) for installation instructions.

## Usage

To use the IMDb Seeder and Searcher Client, follow these steps:

1. Configure the MeiliSearch server.

| Option             | Description                                     | Example                                            |
|--------------------|-------------------------------------------------|----------------------------------------------------|
| --meili-host       | The fully qualified Host with protocol and port | http://localhost:7700                              |
| --meili-master-key | API Key for accessing Meilisearch               | (optional) only required if MEILI_ENV!=development |
| --meili-timeout    | Duration to timeout connection after            | 5s                                                 |

2. Run the IMDb Seeder to index the IMDb data into MeiliSearch. See [Tests](./pkg/imdb-searcher/imdb-searcher_test.go) for some examples.

```bash
go run main.go --meili-host=http://localhost:7700
```
3. Use the [IMDb Searcher Client](./pkg/imdb-searcher/imdb-searcher.go) to search the data or alternatively use the
   MeiliSearch API against the index `imdb`.
## Local Development
There is a [compose file](compose.yml) that can be used to start a local MeiliSearch server and compile and run the IMDb Seeder.

Note: Not affiliated with IMDb. This is a personal project for educational purposes only.
