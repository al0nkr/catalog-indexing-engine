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

	http.ListenAndServe(":8080", r)
}
