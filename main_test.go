package main

import (
	"testing"

	"github.com/truvami/decoder/internal/logger"
)

func TestMain(t *testing.T) {
	logger.NewLogger()
	defer logger.Sync()
	main()
}
