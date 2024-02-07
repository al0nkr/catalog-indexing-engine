package index

import(
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

func IndexDataHandler(w http.ResponseWriter, r *http.Request) {
	// Read the JSON data from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse the JSON data into a slice of maps
	var products []map[string]interface{}
	err = json.Unmarshal(body, &products)
	if err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}

	// Index the documents
	startTime := time.Now()
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "LqcRikNbtOTl0zOTwp746huTR8bR83LPI_xXeMhAnMo",
	})
	
	_, err = client.Index("products").AddDocuments(products)
	if err != nil {
		http.Error(w, "Error indexing documents", http.StatusInternalServerError)
		return
	}

	// Calculate the elapsed time
	elapsedTime := time.Since(startTime)

	// Prepare the response
	response := struct {
		Time       time.Duration `json:"time"`
		Error      string        `json:"error,omitempty"`
		AckMessage string        `json:"ack_message"`
	}{
		Time:       elapsedTime,
		AckMessage: "Data indexed successfully",
	}

	// Convert the response to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the response
	_, _ = w.Write(responseJSON)
}

