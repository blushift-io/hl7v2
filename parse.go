package hl7v2

// ParserOptions holds configuration options for message parsing.
type ParserOptions struct {
	preParse   []RawTransform
	postparse  []RawTransform
	onlyHeader bool
}

// ParserOption is a function that configures ParserOptions.
type ParserOption func(*ParserOptions)

// NewParserOptions creates a new ParserOptions struct configured by the given options.
func NewParserOptions(opts ...ParserOption) *ParserOptions {
	options := &ParserOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

// FixLineEndings returns a ParserOption that normalizes message line endings prior to parsing.
func FixLineEndings() ParserOption {
	return func(options *ParserOptions) {
		options.preParse = append(options.preParse, fixLineEndings)
	}
}

// PreParse returns a ParserOption that applies a custom RawTransform before parsing.
func PreParse(t RawTransform) ParserOption {
	return func(options *ParserOptions) {
		options.preParse = append(options.preParse, t)
	}
}

// OnlyHeader returns a ParserOption that stops parsing after the message header.
func OnlyHeader() ParserOption {
	return func(options *ParserOptions) {
		options.onlyHeader = true
	}
}
