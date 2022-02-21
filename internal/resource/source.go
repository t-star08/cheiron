package resource

import (
	"github.com/t-star08/cheiron/internal/preparator"
)

type Source struct {
	PathToSource	string
	Contents		[]string
	Strategy		*preparator.SuffixPreference
	Message			*Message
}

func newSource(pathToSource string) *Source {
	return &Source {
		PathToSource: pathToSource,
		Message: newMessage(),
	}
}
