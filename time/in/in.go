package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Source struct {
	Interval string `json:"interval"`
}

type Version struct {
	Time time.Time `json:"time"`
}

type Request struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Response struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type Metadata []MetadataField

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

	destination := os.Args[1]

	err := os.MkdirAll(destination, 0755)
	if err != nil {
		fatal("creating destination", err)
	}

	file, err := os.Create(filepath.Join(destination, "input"))
	if err != nil {
		fatal("creating input file", err)
	}

	defer file.Close()

	var request Request

	err = json.NewDecoder(io.TeeReader(os.Stdin, file)).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	metadata := Metadata{
		{"interval", request.Source.Interval},
		{"now", time.Now().String()},
	}

	if !request.Version.Time.Equal(time.Time{}) {
		metadata = append(
			metadata,
			MetadataField{"previous", request.Version.Time.String()},
		)
	}

	json.NewEncoder(os.Stdout).Encode(Response{
		Version:  Version{Time: time.Now()},
		Metadata: metadata,
	})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
