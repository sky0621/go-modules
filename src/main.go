package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

func main() {
	ctx := context.Background()
	project := os.Getenv("PUB_PROJECT")
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
	client.Collection("")
}
