package arrowOptions

type Options struct {
	Practice	bool
	Overwrite	bool
	Secret		bool
	Quiet		bool
}

func New() *Options {
	return &Options{}
}
