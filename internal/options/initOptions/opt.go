package initOptions

type Options struct {
	Force	bool
	Any		bool
}

func New() *Options {
	return &Options{}
}
