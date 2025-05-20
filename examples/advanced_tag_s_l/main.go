package main

import (
	"log"

	"github.com/truvami/decoder/pkg/decoder"
	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func main() {
	log.Println("initializing tag S / L decoder...")
	d := tagsl.NewTagSLv1Decoder()

	// decode data
	log.Println("decoding data...")
	data, err := d.Decode("0002c420ff005ed85a12b4180719142607", 1, "")
	if err != nil {
		panic(err)
	}

	// check if decoded payload has the GNSS feature
	if !data.Is(decoder.FeatureGNSS) {
		panic("decoded payload does not have GNSS feature")
	}

	// cast to GNSS data
	gnssData, ok := data.Data.(decoder.UplinkFeatureGNSS)
	if !ok {
		panic("failed to cast to GNSS data")
	}

	// print GNSS data
	log.Printf("Latitude: %f\n", gnssData.GetLatitude())
	log.Printf("Longitude: %f\n", gnssData.GetLongitude())
	log.Printf("Altitude: %f\n", gnssData.GetAltitude())
}
