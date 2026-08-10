package hl7v2

type ParserOptions struct {
	preParse   []RawTransform
	postparse  []RawTransform
	onlyHeader bool
}

type ParserOption func(*ParserOptions)

func NewParserOptions(opts ...ParserOption) *ParserOptions {
	options := &ParserOptions{}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

func FixLineEndings() ParserOption {
	return func(options *ParserOptions) {
		options.preParse = append(options.preParse, fixLineEndings)
	}
}

func PreParse(t RawTransform) ParserOption {
	return func(options *ParserOptions) {
		options.preParse = append(options.preParse, t)
	}
}

func OnlyHeader() ParserOption {
	return func(options *ParserOptions) {
		options.onlyHeader = true
	}
}
