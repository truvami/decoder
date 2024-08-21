package main

import (
	"encoding/json"
	"log"

	"github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
)

func main() {
	log.Println("initializing tag S / L decoder...")
	d := nomadxs.NewNomadXSv1Decoder()

	// decode data
	log.Println("decoding data...")
	data, err := d.Decode("0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f09a71d2e", 1, "")
	if err != nil {
		panic(err)
	}

	// data to json
	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// print json
	log.Printf("result: %s\n", j)
}
