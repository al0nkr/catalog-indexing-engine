package main

import (
	"github.com/al0nkr/catalog-indexing-engine/search"
	"github.com/al0nkr/catalog-indexing-engine/index"
	"github.com/go-chi/chi"
	"github.com/meilisearch/meilisearch-go"
	"net/http"
)

func main() {
	r := chi.NewRouter()

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "LqcRikNbtOTl0zOTwp746huTR8bR83LPI_xXeMhAnMo",
	})
	r.Post("/index", index.IndexDataHandler(client))
	http.ListenAndServe(":8080", r)
}
