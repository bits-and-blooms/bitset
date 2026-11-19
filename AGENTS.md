# Instructions for AI agents

This file is for AI coding agents (Claude Code, Codex, Copilot, Cursor and
similar tools) and for the people who use them to work on
`github.com/bits-and-blooms/bitset`.

## Human review policy

Every issue and pull request in this project is read and assessed by a human
being. Please respect their time.

- **Do not generate answers automatically.** Do not post AI-written comments,
  replies, reviews or triage on issues or pull requests. If you use an agent to
  help you, you are responsible for what you submit: read it, understand it, and
  write the discussion yourself.
- **Every issue must be well motivated.** Explain the concrete problem: what you
  did, what you expected, what happened, and the Go version and platform. Include
  a minimal reproducer. Do not open issues about speculative, cosmetic or
  stylistic concerns, and do not file bulk reports produced by a tool.
- **Optimizations and bug fixes are welcome**, but they must be of high quality.
  A human maintainer will assess them, so they must be easy to review:
  - A bug fix comes with a regression test that fails before the fix and passes
    after it.
  - An optimization comes with benchmark results (before and after, ideally with
    `benchstat`) on realistic inputs, and states the hardware and Go version.
    A change that makes the code harder to read needs a measurable gain.
  - Keep each pull request small and focused on one change. Do not mix in
    reformatting, renames or unrelated cleanups.
- Do not open pull requests that only reword comments or documentation, add
  generated boilerplate, or "modernize" code without a clear benefit.

If you are an agent acting on someone's behalf: do not open issues, pull
requests or comments on GitHub unless the person explicitly asks you to, and
show them the exact text first.

## Project overview

Package `bitset` maps non-negative integers to boolean values. It is a small,
widely used library with no dependencies, so stability and backward
compatibility matter more than new features.

- `bitset.go`: the `BitSet` type and almost all of the API.
- `bitset_iter.go`: Go 1.23 iterators (`EachSet`), behind a `go1.23` build tag.
- `popcnt.go`, `select.go`: population count and select helpers.
- `pext.gen.go`: **generated** lookup table. Do not edit it by hand; regenerate
  it with `go run cmd/pextgen/main.go -pkg=bitset`.
- `*_test.go`: unit, example, fuzz (`bitset_fuzz_test.go`) and benchmark
  (`bitset_benchmark_test.go`) tests.

## Constraints

- **Go 1.16 compatibility.** `go.mod` declares `go 1.16`, and CI tests every Go
  release from 1.16 onwards. Do not use generics, `min`/`max` builtins, `clear`,
  range-over-int, `slices`, `maps`, `math/bits` functions added after 1.16, or
  other newer language and library features in files without a build tag.
  Code that needs a newer Go must sit behind a `//go:build` tag (with the
  matching `// +build` line), as `bitset_iter.go` does.
- **32-bit platforms.** `uint` is 32 bits on 386 and ARM. Watch for overflow in
  index and length arithmetic, and test with `GOARCH=386`.
- **No new dependencies.**
- **Public API.** Do not remove or change the signature or semantics of an
  exported identifier. New exported functions need a clear use case and doc
  comments.
- **Serialization format.** `WriteTo`/`ReadFrom`, `MarshalBinary`/
  `UnmarshalBinary` and the JSON encoding must stay compatible with data written
  by earlier releases.
- Keep the code idiomatic and consistent with the surrounding style; run
  `gofmt`.

## Testing

Run these before proposing a change:

```bash
go vet ./...
go test ./...
GOARCH=386 go test ./...            # on Linux or Windows
GOOS=linux GOARCH=386 go vet ./...  # on macOS, which cannot run 386 binaries
```

For changes to core operations, also run the fuzz tests briefly
(`./run_fuzz_tests.sh 10s`), and for performance work, the benchmarks
(`go test -run=^$ -bench=. -count=10`).
