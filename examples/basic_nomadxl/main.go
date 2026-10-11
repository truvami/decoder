package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/truvami/decoder/pkg/decoder/nomadxl/v1"
)

func main() {
	fmt.Println("initializing nomad XL decoder...")
	d := nomadxl.NewNomadXLv1Decoder()

	fmt.Println("decoding data...")
	data, err := d.Decode(context.Background(), "0000793000020152004b6076000c838c00003994", 103)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
