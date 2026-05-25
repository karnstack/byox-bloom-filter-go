# karnstack/byox-bloom-filter-go

Starter template (Go) for the [Bloom Filter](https://karnstack.com/build/bloom-filter) primitive on karnstack.

Six stages. Paper-backed tests. You implement the interface; karnstack tells you what to read at each stage.

## Prerequisites

[mise](https://mise.jdx.dev/) is the only thing you need installed globally. It pins Go 1.26 for this repo and runs the stage tasks. If you do not want to install mise, the equivalent `go test` commands are documented under [Without mise](#without-mise) below.

Install mise:

```bash
curl https://mise.run | sh
```

## Quick start

```bash
mise trust              # allow this repo's .mise.toml (one time)
mise install            # installs Go 1.26 if you do not have it
mise run stage 1        # runs the tests for stage 1 (they fail until you implement)
```

Open [stage 1 on karnstack](https://karnstack.com/build/bloom-filter/01-bit-array-and-hashing). Implement `bloom/bloom.go` until `mise run stage 1` passes. Then move on:

```bash
mise run stage 2
```

`mise run all` runs every stage in one go. `mise run bench` runs the benchmarks used by stage 4 and stage 5.

## Layout

```
.
├── .mise.toml                          # toolchain + tasks
├── go.mod
└── bloom/
    ├── bloom.go                        # you implement here
    ├── stage01_bit_array_test.go
    ├── stage02_multi_hash_test.go
    ├── stage03_sizing_test.go
    ├── stage04_blocked_test.go
    ├── stage05_concurrent_test.go
    └── stage06_serialize_test.go
```

Tests live in the `bloom` package as `*_test.go` files (Go convention). The test files declare `package bloom_test` so they only see the exported API, which is the same surface a real consumer would use.

## Stages

1. Bit array and single hash
2. Multiple hashes (Kirsch-Mitzenmacher)
3. Optimal sizing math
4. Cache-line-blocked layout
5. Concurrent-safe Add
6. Serialize and saturation

Each stage is described on karnstack. Read first, then implement.

## What you are building

A constant-size data structure that says "definitely not in the set" or "maybe in the set" in `O(k)` time, with a tunable false-positive rate. The structure inside every production LSM-tree (RocksDB, LevelDB, Cassandra) used to skip disk reads on missing keys.

## Papers cited

- Bloom, B. (1970). [Space/Time Trade-offs in Hash Coding with Allowable Errors](https://dl.acm.org/doi/10.1145/362686.362692). CACM 13(7).
- Kirsch, A.; Mitzenmacher, M. (2006). [Less Hashing, Same Performance: Building a Better Bloom Filter](https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf). ESA 2006.
- Putze, F.; Sanders, P.; Singler, J. (2007). [Cache-, Hash- and Space-Efficient Bloom Filters](https://doi.org/10.1007/978-3-540-72845-0_9). WEA 2007.

## Without mise

If you do not want to install mise, ensure you have Go 1.26+ installed and run:

```bash
# Stage 1
go test -v -run '^TestStage01_' ./bloom/...

# Stage N (replace 01 with the zero-padded stage number)
go test -v -run '^TestStageNN_' ./bloom/...

# All stages
go test -v ./bloom/...
```

## License

MIT. See [LICENSE](LICENSE). Your fork is yours.
