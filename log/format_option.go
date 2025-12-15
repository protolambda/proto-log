package log

// FormatConfig configures how the common handlers format their output.
type FormatConfig struct {
	// UseColor formats the output with color-codes.
	// No-op in handlers where color is not supported.
	UseColor bool
	// IncludeSource shows the file path and line number of the log source
	IncludeSource bool
	// ExcludeTime shows the date / time
	ExcludeTime bool
	// SourceRelDir is the dir to resolve sources to as relative files
	SourceRelDir string
}

func (cfg *FormatConfig) Apply(opts ...FormatOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

type FormatOption func(cfg *FormatConfig)

// WithColor sets FormatConfig.UseColor
func WithColor(useColor bool) FormatOption {
	return func(cfg *FormatConfig) {
		cfg.UseColor = useColor
	}
}

// WithIncludeSource sets FormatConfig.IncludeSource
func WithIncludeSource(includeSource bool) FormatOption {
	return func(cfg *FormatConfig) {
		cfg.IncludeSource = includeSource
	}
}

// WithExcludeTime sets FormatConfig.ExcludeTime
func WithExcludeTime(excludeTime bool) FormatOption {
	return func(cfg *FormatConfig) {
		cfg.ExcludeTime = excludeTime
	}
}

// WithSourceRelDir sets FormatConfig.SourceRelDir
func WithSourceRelDir(dir string) FormatOption {
	return func(cfg *FormatConfig) {
		cfg.SourceRelDir = dir
	}
}
