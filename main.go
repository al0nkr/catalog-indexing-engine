package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/meilisearch/meilisearch-go"
	"main/index"
	"main/search"
	"net/http"
	"os/exec"
	"time"
)

func executeMeilisearch() {
	cmd := exec.Command("./meilisearch", "--master-key", "LqcRikNbtOTl0zOTwp746huTR8bR83LPI_xXeMhAnMo", "--http-payload-size-limit", "536870912")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Meilisearch command:", err)
		fmt.Println("Meilisearch stderr:", string(stdoutStderr))
		return
	}
	fmt.Println("Meilisearch output:", string(stdoutStderr))
}

func main() {

	go executeMeilisearch()
	r := chi.NewRouter()
	r.Post("/index", index.IndexDataHandler)

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "LqcRikNbtOTl0zOTwp746huTR8bR83LPI_xXeMhAnMo",
	})

	r.Get("/search", search.SearchHandler(client))

	server := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  10 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("Server started at :3000")
}
