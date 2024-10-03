package main

import (
	"encoding/json"
	"log"

	"github.com/truvami/decoder/pkg/decoder/nomadxl/v1"
)

func main() {
	log.Println("initializing nomad XL decoder...")
	d := nomadxl.NewNomadXLv1Decoder()

	// decode data
	log.Println("decoding data...")
	data, err := d.Decode("0000793000011702004B6089000C836AFFFFC13A", 103, "")
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
