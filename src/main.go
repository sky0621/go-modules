package main

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
)

func main() {
	ctx := context.Background()
	project := os.Getenv("PUB_PROJECT")
	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 常駐型
	for {
		if err := processRegions(ctx, client); err != nil {
			log.Print(err)
		}
		time.Sleep(3 * time.Second)
	}
}

func processRegions(ctx context.Context, client *firestore.Client) error {
	regions, err := client.Collection("region").Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	rwg := &sync.WaitGroup{}
	rch := make(chan struct{}, 3)

	for _, region := range regions {
		rwg.Add(1)
		rch <- struct{}{}
		go func(regionID string) {
			if err := processSchools(ctx, client, regionID, rch, rwg); err != nil {
				log.Print(err)
			}
		}(region.Ref.ID)
	}
	rwg.Wait()
	return nil
}

func processSchools(ctx context.Context, client *firestore.Client, regionID string, rch chan struct{}, rwg *sync.WaitGroup) error {
	defer func() {
		<-rch
		rwg.Done()
	}()

	schools, err := client.Collection("region").Doc(regionID).Collection("school").Where("state", "in", []string{"UNSYNCED", "SYNCED"}).Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	swg := &sync.WaitGroup{}
	sch := make(chan struct{}, 3)

	for _, school := range schools {
		swg.Add(1)
		sch <- struct{}{}
		go func(regionID, schoolID string) {
			if err := processTasksInSchool(ctx, client, regionID, schoolID, sch, swg); err != nil {
				log.Print(err)
			}
		}(regionID, school.Ref.ID)
	}
	swg.Wait()
	return nil
}

func processTasksInSchool(ctx context.Context, client *firestore.Client, regionID, schoolID string, sch chan struct{}, swg *sync.WaitGroup) error {
	defer func() {
		// to SYNCED
		_, err := client.Collection("region").Doc(regionID).Collection("school").Doc(schoolID).Set(ctx, map[string]interface{}{
			"state": "SYNCED",
		}, firestore.MergeAll)
		if err != nil {
			log.Print(err)
		}

		<-sch
		swg.Done()
	}()

	// to SYNCING
	_, err := client.Collection("region").Doc(regionID).Collection("school").Doc(schoolID).Set(ctx, map[string]interface{}{
		"state": "SYNCING",
	}, firestore.MergeAll)
	if err != nil {
		return err
	}

	tasks, err := client.Collection("region").Doc(regionID).Collection("school").
		Doc(schoolID).Collection("operation").OrderBy("operationSequence", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		d := task.Data()
		order, ok := d["order"]
		if ok {
			ods := strings.Split(order.(string), ":")
			if len(ods) > 0 {
				od := ods[0]
				switch od {
				case "add-school":
					time.Sleep(5 * time.Second)
				case "edit-school":
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
		log.Printf("%#+v", d)
	}
	return nil
}
