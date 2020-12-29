package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Server")
	b, err := downloadFile()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(b))
}

func downloadFile() ([]byte, error) {
	fmt.Println("Download Time!")
	bucket := "mybucket"
	object := "codes.yml"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := &storage.Query{Prefix: ""}

	var names []string
	it := client.Bucket(bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}

	fmt.Println(names)

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	fmt.Printf("Blob %v downloaded.\n", object)
	return data, nil
}
