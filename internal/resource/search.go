package resource

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/t-star08/cheiron/internal/config"
	"github.com/t-star08/cheiron/internal/preparator"
	"github.com/t-star08/cheiron/pkg/fileUtil/filePathCollector"
	"github.com/t-star08/hand/pkg/evaluator/fileName"
	"github.com/t-star08/hand/pkg/fileUtil/fileSearcher"
)

var (
	exEvaluator = fileName.NewExactEvaluator()
	parEvaluator = fileName.NewPartialEvaluator()
)

func GlobProjectRootPath(targetPath string) (*Project, error) {
	res := newProject()
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return res, err
	} else {
		res.PathToProjectRoot = targetPath
	}
	return res, nil
}

func FindBranchPaths(prj *Project, targetPath string) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return err
	}
	prj.Branches[targetPath] = newBranch(targetPath)
	return nil
}

func GlobBranchPaths(prj *Project, targetPath string) error {
	if branchPaths, err := filePathCollector.ForwardCollectMatchedPath(parEvaluator, targetPath, prj.PathToProjectRoot); err != nil {
		return err
	} else {
		for _, branchPath := range branchPaths {
			if _, exist := prj.Branches[branchPath]; exist {
				continue
			}
			prj.Branches[branchPath] = newBranch(branchPath)
		}
	}
	return nil
}

func GlobInsertTargetFilePaths(prj *Project, targetPath string) error {
	for pathToBranch, branch := range prj.Branches {
		if branch.Ignored {
			continue
		}

		if pathToInsertTarget, err := fileSearcher.ForwardSearchUntilOrderedDepth(exEvaluator, targetPath, pathToBranch, 1); err != nil {
			branch.Message.Failed("cannot glob insert target file")
		} else {
			branch.FoundInsertTarget = true
			branch.InsertTarget = newInsertTarget(pathToInsertTarget)
		}
	}
	return nil
}

func GlobTemplateFilePaths(prj *Project, pathToConfig string, routine []*config.Routine) error {
	sort.SliceStable(
		routine,
		func(i, j int) bool {
			return routine[i].Priority < routine[j].Priority
		},
	)

	for _, arrow := range routine {
		pathToTemplate := fmt.Sprintf("%s/%s", pathToConfig, arrow.Template)
		if _, err := os.Stat(pathToTemplate); err != nil {
			continue
		}
		prj.PathToTemplates = append(prj.PathToTemplates, pathToTemplate)
	}
	
	return nil
}

func GlobBranchSourceFilePaths(br *Branch, strategies preparator.SuffixPreferences) error {
	for _, point := range br.InsertTarget.Points {
		point.Message.Reset()

		// check suffix valid
		suffix := point.RequirePath[strings.LastIndex(point.RequirePath, "."):]
		if !(strategies.Has(".*") || strategies.Has("*")) {
			if !strategies.Has(suffix) {
				point.Message.Failed("out of strategy suffixes")
				continue
			}
		}
		// check path valid
		pathToSource := fmt.Sprintf("%s/%s", br.PathToBranch, point.RequirePath)
		if _, err := os.Stat(pathToSource); os.IsNotExist(err) {
			if !point.WhetherMust {
				point.Message.Failed(fmt.Sprintf("cannot find %s", pathToSource))
				continue
			}
			br.Sources = map[string]*Source{}
			return fmt.Errorf("lack %s", pathToSource)
		}
		br.Sources[pathToSource] = newSource(pathToSource)
	}
	return nil
}
