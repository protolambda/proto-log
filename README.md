# proto-log

Yet another Go logging library.

Optimized for convenience, typing, compatibility.

Features:
- Built on top of `slog` (introduced in Go 1.21, see [Go blog](https://go.dev/blog/slog))
  - All handlers and utils can be wrapped with a `*slog.Logger` or other logger frontends,
    for logging compatibility with projects that do not use `proto-log/log.Logger`.
- Extended `Logger` interface:
  - `Trace`, `TraceContext`:
  - `Crit`, `CritContext`: without `os.Exit`, instead attach a `PostHandlerMod`
    to follow-up crit logs with your preferred crit handling.
  - `Context` to access the default context
  - `WithContext` to make a logger clone and attach a new default context
- Handler `Unwrap` pattern, to find handler-wrappers easily
- A set of `HandlerMod` to adjust log-handlers of (sub-)loggers at runtime:
  - `ContextMod` to adjust the default `context`
  - `LevelMod` to adjust the log-level
  - `CapturingMod` to capture logging
  - `PostProcessMod` to post-process log records (e.g. handle special log levels)
- A set of `slog.Handler` implementations:
  - `DiscardHandler` 
  - `JSONHandler`
  - `LogfmtHandler`: for human-readable but Loki-compatible logging.
  - `TerminalHandler`:
    - Like logfmt, but stylish, with automatic source-filepath and field padding.
    - Looks for `TerminalString() string` on types for custom formatting.
    - `uint64`, `*big.Int` and `*uint256.Int` are logged with `_` thousand-separators.
- `TestLogger`: minimal test log-handling stack on top
  of `T.Output()` (introduced in [Go 1.23](https://github.com/golang/go/issues/59928))
  - Can be customized with additional `HandlerMod`
  - `Logger.Crit` / `Logger.CritContext` are followed up with `T.FailNow()`
- `FormatOption` to configure formatting of handlers:
  - Option to exclude time, for logging in Go `Example` output to be stable
  - Option to resolve file-paths of source-file data to relative paths
  - Option to color output of `TerminalHandler`
- No dependencies


## Usage

See: [Log example tests](./log/log_example_test.go)


## Credits

Some other log libraries influenced this:
- [`log15`](https://github.com/inconshreveable/log15) v1
  ([obligatory xkcd](https://imgs.xkcd.com/comics/standards.png)):
  - Originally [introduced](https://github.com/ethereum/go-ethereum/pull/3696) in go-ethereum.
  - Pushed forward Go structured log fields, sub-logger contexts.
  - Influenced terminal log-formatting in go-ethereum / optimism / others (see notice in `format.go`).
- Logging was adapted in [optimism](https://github.com/ethereum-optimism/optimism/):
  - Optimism introduced initial log-capturer (forked now).
  - First version of handler-wrapping idea for composition.
  - MIT, see `LICENSE` entry.
- [`slog`](https://go.dev/blog/slog) was [introduced in Go 1.21](https://go.dev/doc/go1.21).
  - Go finally got standard modern structured logging.
  - [Introduced](https://github.com/ethereum/go-ethereum/pull/28187) into go-ethereum.
  - Reduced code, adopted standard log-level values and utils.


# License

MIT License, see [LICENSE](./LICENSE) file.
