package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/truvami/decoder/pkg/decoder/tagxl/v1"
	"github.com/truvami/decoder/pkg/solver/aws"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	logger := zap.NewExample()
	defer func() {
		_ = logger.Sync() // flushes buffer, if any
	}()

	solver, err := aws.NewAwsPositionEstimateClient(ctx, logger)
	if err != nil {
		panic(err)
	}

	logger.Info("initializing tag XL decoder...")
	d := tagxl.NewTagXLv1Decoder(ctx, solver, logger)

	// decode data
	logger.Info("decoding data...")
	data, err := d.Decode("05ab859590e78d0cc1805a9428b2de73d80cc9c9a3329a01a5e3cba3546b7454395747a1cd6effd2fdeebefe8fac39a60e", 192)
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
