package search

import (
	"net/http"
	"encoding/json"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

func searchHandler(client *meilisearch.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Missing query parameter", http.StatusBadRequest)
			return
		}

		startTime := time.Now()

		searchResponse, err := client.Index("products").Search(query, &meilisearch.SearchRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		elapsedTime := time.Since(startTime)

		var output []string
		outputType := r.URL.Query().Get("output")
		switch outputType {
		case "titles":
			output = getTitlesFromHits(searchResponse.Hits)
		case "productIDs":
			output = getProductIDsFromHits(searchResponse.Hits)
		default:
			http.Error(w, "Invalid output type", http.StatusBadRequest)
			return
		}

		response := struct {
			Output      []string      `json:"output"`
			ElapsedTime time.Duration `json:"elapsed_time"`
		}{
			Output:      output,
			ElapsedTime: elapsedTime,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func getTitlesFromHits(hits []interface{}) []string {
	var titles []string
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		title, ok := hitMap["product_title"].(string)
		if ok {
			titles = append(titles, title)
		}
	}
	return titles
}

func getProductIDsFromHits(hits []interface{}) []string {
	var productIDs []string
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		productID, ok := hitMap["product_id"].(string)
		if ok {
			productIDs = append(productIDs, productID)
		}
	}
	return productIDs
}

// func searchTitlesHandler(client *meilisearch.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		query := r.URL.Query().Get("query")
// 		if query == "" {
// 			http.Error(w, "Missing query parameter", http.StatusBadRequest)
// 			return
// 		}

// 		startTime := time.Now()

// 		searchResponse, err := client.Index("products").Search(query, &meilisearch.SearchRequest{})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		elapsedTime := time.Since(startTime)

// 		var titles []string
// 		for _, hit := range searchResponse.Hits {
// 			hitMap, ok := hit.(map[string]interface{})
// 			if !ok {
// 				continue
// 			}
// 			title, ok := hitMap["product_title"].(string)
// 			if ok {
// 				titles = append(titles, title)
// 			}
// 		}

// 		response := struct {
// 			Titles       []string      `json:"titles"`
// 			ElapsedTime  time.Duration `json:"elapsed_time"`
// 		}{
// 			Titles:       titles,
// 			ElapsedTime:  elapsedTime,
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 	}
// }