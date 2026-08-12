package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func main() {
	fmt.Println("initializing tag S / L decoder...")
	d := tagsl.NewTagSLv1Decoder()

	fmt.Println("decoding data...")
	data, err := d.Decode(context.Background(), "0002c420ff005ed85a12b4180719142607", 1)
	if err != nil {
		panic(err)
	}

	j, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))
}
