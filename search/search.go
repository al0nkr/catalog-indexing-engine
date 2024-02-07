package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/meilisearch/meilisearch-go"
	"google.golang.org/api/option"
)

func SearchHandler(client *meilisearch.Client) http.HandlerFunc {
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
		case "ids":
			output = getProductIDsFromHits(searchResponse.Hits)
		case "similar":
			similarProducts, _ := getSimilarProducts(client, query)
			output = similarProducts
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

func searchTopTitles(client *meilisearch.Client, query string) ([]string, error) {
	searchResponse, err := client.Index("products").Search(query, &meilisearch.SearchRequest{})
	if err != nil {
		return nil, err
	}

	var titles []string
	for _, hit := range searchResponse.Hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		title, ok := hitMap["product_title"].(string)
		if ok {
			titles = append(titles, title)
		}
	}

	return titles, nil
}

func getSimilarProducts(client *meilisearch.Client,query string) ([]string, error) {

	ctx := context.Background()
	client2, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyA2YUU940oZmaj8vOpVSLl5LeWudvQ0HnE"))
	if err != nil {
		return nil,err
	}
	defer client2.Close()

	// For text-only input, use the gemini-pro model
	model := client2.GenerativeModel("gemini-pro")
	prompt := fmt.Sprintf("provide 5 relevant and unique simple product search queries one might buy along with : %s {provide user queries only}", query)
	response, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	var similarQueries []string

	for _, cand := range response.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partStr := fmt.Sprintf("%v", part)
				lines := strings.Split(partStr, "\n")

				// Iterate through each line
				for _, line := range lines {
					similarQueries = append(similarQueries, line[2:])
				}
			}
		}
	}

	var similarProducts []string

	// Search for similar queries in the index and get top titles
	for _, query := range similarQueries {
		searchTitles, err := searchTopTitles(client, query)
		if err != nil {
			return nil, err
		}

		// Append top 2 titles to similarProducts
		for i, title := range searchTitles {
			if i >= 2 {
				break
			}
			similarProducts = append(similarProducts, title)
		}
	}

	if err != nil {
		return nil, err
	}

	return similarProducts, nil
}
