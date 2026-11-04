package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/truvami/decoder/pkg/decoder/nomadxs/v1"
)

func main() {
	fmt.Println("initializing nomad XS decoder...")
	d := nomadxs.NewNomadXSv1Decoder()

	fmt.Println("decoding data...")
	data, err := d.Decode(context.Background(), "0002c420ff005ed85a12b4180719142607240001ffbaffc2fc6f00d71d2e", 1)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
