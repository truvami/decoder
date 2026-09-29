//go:build ignore

package gorules

import "github.com/quasilyte/go-ruleguard/dsl"

func errorTextMatch(m dsl.Matcher) {
	m.Match(`strings.Contains($e.Error(), $_)`, `strings.HasPrefix($e.Error(), $_)`,
		`strings.HasSuffix($e.Error(), $_)`, `$e.Error() == $_`, `$e.Error() != $_`).
		Where(m["e"].Type.Implements("error")).
		Report(`compare errors with errors.Is/errors.As, not by message text`)
}

func grpcInternalLeak(m dsl.Matcher) {
	m.Import("google.golang.org/grpc/status")
	m.Import("google.golang.org/grpc/codes")
	m.Match(`status.Errorf(codes.Internal, $_, $_, $*_)`).
		Report(`formatted codes.Internal can leak internals; log the cause, return a fixed message`)
	m.Match(`status.Error(codes.Internal, $e.Error())`).
		Where(m["e"].Type.Implements("error")).
		Report(`do not send err.Error() to gRPC clients; log it, return a fixed message`)
}

func apiMessageLeak(m dsl.Matcher) {
	m.Match(`$t{$*_, Message: $e.Error(), $*_}`).
		Where(m["e"].Type.Implements("error")).
		Report(`do not send err.Error() to API clients; log it, return a stable message`)
}

func promErrorLabel(m dsl.Matcher) {
	m.Match(`$v.WithLabelValues($*args)`).
		Where(m["args"].Text.Matches(`\.Error\(\)`)).
		Report(`err.Error() as a label value is unbounded; use a fixed error class`)
}

func httpClientNoTimeout(m dsl.Matcher) {
	m.Match(`http.Client{$*fields}`).
		Where(!m["fields"].Text.Matches(`\bTimeout\s*:`)).
		Report(`http.Client without Timeout can hang forever; set Timeout`)
}

func unboundedBodyRead(m dsl.Matcher) {
	m.Match(`io.ReadAll($x.Body)`).
		Report(`wrap the body in io.LimitReader or http.MaxBytesReader before io.ReadAll`)
}
