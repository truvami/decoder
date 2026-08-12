package main

import (
	"context"
	"fmt"

	"github.com/truvami/decoder/pkg/decoder"
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

	if !data.Is(decoder.FeatureGNSS) {
		panic("decoded payload does not have GNSS feature")
	}

	gnssData, ok := data.Data.(decoder.UplinkFeatureGNSS)
	if !ok {
		panic("failed to cast to GNSS data")
	}

	fmt.Printf("Latitude: %f\n", gnssData.GetLatitude())
	fmt.Printf("Longitude: %f\n", gnssData.GetLongitude())
	fmt.Printf("Altitude: %f\n", gnssData.GetAltitude())
}
