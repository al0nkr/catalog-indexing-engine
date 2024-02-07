package main

import (
	"github.com/go-chi/chi"
	"github.com/meilisearch/meilisearch-go"
	"main/index"
	"main/search"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Post("/index", index.IndexDataHandler)

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "LqcRikNbtOTl0zOTwp746huTR8bR83LPI_xXeMhAnMo",
	})

	r.Get("/search", search.SearchHandler(client))
	http.ListenAndServe(":3000", r)
}
