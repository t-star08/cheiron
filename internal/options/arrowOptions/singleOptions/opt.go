package singleOptions

import (
	"github.com/t-star08/cheiron/internal/options/arrowOptions"
)

type Options struct {
	arrowOptions.Options
	Target	string
}

func New() *Options {
	return &Options{}
}
