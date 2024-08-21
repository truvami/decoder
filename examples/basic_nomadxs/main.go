package main

import (
	"encoding/json"
	"log"

	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func main() {
	log.Println("initializing tag S / L decoder...")
	d := tagsl.NewTagSLv1Decoder()

	// decode data
	log.Println("decoding data...")
	data, err := d.Decode("8002cdcd1300744f5e166018040b14341a", 1, "")
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