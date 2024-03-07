package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/meilisearch/meilisearch-go"
	nostr "github.com/nbd-wtf/go-nostr"
)

const name = "yabume-search-indexer"

const version = "0.0.4"

var revision = "HEAD"

func getEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("missing env var for %v", name)
	}
	return value
}

func main() {
	ctx := context.Background()

	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   getEnv("MEILISEARCH_URL"),
		APIKey: getEnv("MEILISEARCH_KEY"),
	})
	index := client.Index("events")

	relays := strings.Split(getEnv("RELAY_URL"), ",")
	pool := nostr.NewSimplePool(ctx)
	now := nostr.Now()
	sub := pool.SubMany(ctx, relays, nostr.Filters{
		nostr.Filter{Kinds: []int{0, 1, 42, 30023}, Since: &now},
	})
	for ev := range sub {
		task, err := index.AddDocuments(ev)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(task.TaskUID)
		}
	}
}
