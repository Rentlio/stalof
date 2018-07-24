package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		displayUsage()
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		decorateError("Google storage client error", err)
		return
	}

	bucket := client.Bucket(args[1])

	it := bucket.Objects(ctx, &storage.Query{})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			decorateError("Storage iterator error", err)
		}
		object := bucket.Object(attrs.Name)
		r, err := object.NewReader(ctx)
		if err != nil {
			decorateError("Object reader error", err)
		}
		defer r.Close()
		data, err := ioutil.ReadAll(r)
		if err != nil {
			decorateError("Error reading object data", err)
		}

		// now split data by lines
		lines := strings.Split(string(data), "\n")
		handleLogLines(lines)
	}
}

func handleLogLines(lines []string) {
	for _, line := range lines {
		handleLogLine(line)
	}
}

func handleLogLine(line string) {
	if line == "" {
		return
	}
	logEntry := stackDriverLog{}
	err := json.Unmarshal([]byte(line), &logEntry)
	if err != nil {
		decorateError("Log entry unmarshaling error", err)
		return
	}
	fmt.Println(logEntry.TextPayload)
}

func displayUsage() {
	fmt.Println("Usage: stalof my-bucket-name")
}

func decorateError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", message, err.Error())
}

// struct where we unmarshall stackdriver logs
// we use just textPayload property because we are only
// interested in original text logs.
type stackDriverLog struct {
	TextPayload string `json:"textPayload"`
}
