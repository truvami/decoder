package main

import (
	"log"

	"github.com/truvami/decoder/pkg/decoder/tagsl/v1"
)

func main() {
	log.Println("initializing tag S / L decoder...")
	d := tagsl.NewTagSLv1Decoder()

	// decode data
	log.Println("decoding data...")
	location, _, err := d.DecodePosition("0002c420ff005ed85a12b4180719142607", 1, "")
	if err != nil {
		panic(err)
	}

	log.Printf("latitude: %f, longitude: %f, altitude: %f\n", location.GetLatitude(), location.GetLongitude(), *location.GetAltitude())
}
