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
	log.Print("START__go-subscriber-fs")
	ctx := context.Background()
	project := os.Getenv("PUB_PROJECT")
	log.Printf("ENV:PROJECT:%s", project)
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	schoolsIter := client.Collection("school").Where("state", "==", "ACTIVE").Snapshots(ctx)
	defer schoolsIter.Stop()

	log.Print("LOOP_START!")
	for {
		log.Println("before next")
		school, err := schoolsIter.Next()
		log.Println("after next")
		if err != nil {
			log.Fatalln(err)
		}

		for _, change := range school.Changes {
			go func(changed firestore.DocumentChange) {
				log.Printf("change: %+v\n", changed)
				log.Printf("change.Data: %+v\n", changed.Doc.Data())
				if err := processTasksInSchool(ctx, changed); err != nil {
					log.Print(err)
				}
			}(change)
		}
	}
}

func processTasksInSchool(ctx context.Context, changed firestore.DocumentChange) error {
	operationIter := changed.Doc.Ref.Collection("operation").
		Where("operationSequence", ">", 0).OrderBy("operationSequence", firestore.Asc).Snapshots(ctx)
	defer operationIter.Stop()

	log.Print("[processTasksInSchool] LOOP_START!")
	for {
		log.Println("[processTasksInSchool] before next")
		operation, err := operationIter.Next()
		log.Println("[processTasksInSchool] after next")
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
						time.Sleep(5 * time.Second)
					case "add-grade":
						time.Sleep(15 * time.Second)
					case "add-class":
						time.Sleep(10 * time.Second)
					case "add-teacher":
						time.Sleep(5 * time.Second)
					case "add-student":
						time.Sleep(1 * time.Second)
					}
				}
			}
			log.Printf("[processTasksInSchool] %#+v", d)

			_, err := change.Doc.Ref.Delete(ctx)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	return nil
}
