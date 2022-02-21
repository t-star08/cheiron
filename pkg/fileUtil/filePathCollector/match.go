package filePathCollector

import (
	"fmt"
	"strings"

	"github.com/t-star08/hand/pkg/evaluator"
	"github.com/t-star08/hand/pkg/fileUtil/fileCollector"
)

func ForwardCollectMatchedPath(e *evaluator.Evaluator, targetPath, startPath string) ([]string, error) {
	res := make([]string, 0)
	path := []string {startPath}
	targetPathLevel := len(strings.Split(targetPath, "/"))
	for i, target := range strings.Split(targetPath, "/") {
		nPath := make([]string, 0)
		for len(path) > 0 {
			p := path[len(path)-1]
			path = append(path[:len(path)-1], path[len(path):]...)
			if collection, err := fileCollector.ForwardCollectUntilOrderedDepth(e, target, p, 1); err != nil {
				if err.Error() == "not found" {
					continue
				}
				return collection, err
			} else {
				if i == targetPathLevel - 1 {
					res = append(res, collection...)
				}
				nPath = append(nPath, collection...)
			}
		}
		path = nPath
	}
	if len(res) == 0 {
		return res, fmt.Errorf("not found")
	} else {
		return res, nil
	}
}
