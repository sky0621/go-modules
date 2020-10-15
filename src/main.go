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
			ope, err := change.Doc.Ref.Get(ctx)
			if err != nil {
				log.Fatalln(err)
			}
			d := ope.Data()
			order, ok := d["order"]
			if ok {
				ods := strings.Split(order.(string), ":")
				if len(ods) > 0 {
					od := ods[0]
					switch od {
					case "add-school":
						time.Sleep(5 * time.Second)
					case "add-grade":
						time.Sleep(4 * time.Second)
					case "add-class":
						time.Sleep(3 * time.Second)
					case "add-teacher":
						time.Sleep(2 * time.Second)
					case "add-student":
						time.Sleep(1 * time.Second)
					}
				}
			}
			log.Printf("[operation-Data] %#+v", d)
		}
	}
}
