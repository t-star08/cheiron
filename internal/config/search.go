package config

import (
	"github.com/t-star08/hand/pkg/evaluator/fileName"
	"github.com/t-star08/hand/pkg/fileUtil/fileSearcher"
)

var (
	evaluator = fileName.NewExactEvaluator()
)

func SearchConfigDirPath() (string, error) {
	if pathToConfDir, err := fileSearcher.BackwardSearchUntilOrderedDepth(evaluator, CONF_DIR_NAME, ".", 10); err != nil {
		return pathToConfDir, err
	} else {
		return pathToConfDir, nil
	}
}
