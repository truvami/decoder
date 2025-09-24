package cmd

// Version is the CLI version string. It defaults to "dev" and is overridden at build time via ldflags.
// Example with goreleaser:
//
//	-X github.com/truvami/decoder/cmd.Version={{.Version}}
var Version = "dev"

// Optional metadata (can be set via ldflags as well if desired).
var Commit = ""
var Date = ""
