package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

func main() {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, os.Getenv("PUB_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	operationIter := client.Collection("operation").
		Where("sequence", ">", 0).OrderBy("sequence", firestore.Asc).Snapshots(ctx)
	defer operationIter.Stop()

	for {
		operation, err := operationIter.Next()
		if err != nil {
			log.Fatalln(err)
		}

		for _, change := range operation.Changes {
			d := change.Doc.Data()
			order, ok := d["order"]
			if ok {
				ods := strings.Split(order.(string), ":")
				if len(ods) > 0 {
					od := ods[0]
					switch od {
					case "add-school":
						time.Sleep(10 * time.Second)
					case "add-grade":
						time.Sleep(8 * time.Second)
					case "add-class":
						time.Sleep(6 * time.Second)
					case "add-teacher":
						time.Sleep(4 * time.Second)
					case "add-student":
						time.Sleep(2 * time.Second)
					}
				}
			}
			log.Printf("[operation-Data] %#+v", d)

			_, err := change.Doc.Ref.Delete(ctx)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
