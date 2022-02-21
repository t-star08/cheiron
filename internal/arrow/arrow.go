package arrow

import (
	"fmt"
	"os"
	"strings"

	"github.com/t-star08/cheiron/internal/config"
	"github.com/t-star08/cheiron/internal/document"
	"github.com/t-star08/cheiron/internal/preparator"
	"github.com/t-star08/cheiron/internal/resource"
	"github.com/t-star08/cheiron/internal/result"
	"github.com/t-star08/cheiron/pkg/fileOpe"
	"github.com/t-star08/cheiron/pkg/io/ioFile"
)

const (
	PROTECT_PREFIX = ".>>>["
)

func DetectProjectRoot(pre *preparator.Preparator) (*resource.Project, error) {
	if p, err := resource.GlobProjectRootPath(pre.Cheiron.ProjectRoot); err != nil {
		return p, err
	} else {
		return p, err
	}
}

func DetectBranchPaths(pre *preparator.Preparator, prj *resource.Project) error {
	if err := resource.GlobBranchPaths(prj, pre.Cheiron.Branch); err != nil {
		return err
	} else {
		for pathToBranch, branch := range prj.Branches {
			if info, _ := os.Stat(pathToBranch); !info.IsDir() {
				branch.Ignored = true
				branch.Message.Failed(fmt.Sprintf("%s is not directory", pathToBranch))
			}
		}
		for _, IgnoreBranchName := range pre.Cheiron.Ignore {
			if _, exist := prj.Branches[fmt.Sprintf("%s/%s", prj.PathToProjectRoot, IgnoreBranchName)]; exist {
				prj.Branches[fmt.Sprintf("%s/%s", prj.PathToProjectRoot, IgnoreBranchName)].Ignored = true
				prj.Branches[fmt.Sprintf("%s/%s", prj.PathToProjectRoot, IgnoreBranchName)].Message.Failed(fmt.Sprintf("in %s ignore list", config.CONF_FILE_NAME))
			}
		}
	}
	return nil
}

func DetectInsertTargetFilePaths(pre *preparator.Preparator, prj *resource.Project) error {
	return resource.GlobInsertTargetFilePaths(prj, pre.Cheiron.InsertTarget)
}

func PullInsertTargetSources(prj *resource.Project) error {
	for _, branch := range prj.Branches {
		if branch.Ignored || !branch.FoundInsertTarget {
			continue
		}
		branch.InsertTarget.Source , _ = ioFile.Gets(branch.InsertTarget.PathToInsertTarget)
	}
	return nil
}

func MeetInsertTargetRequirements(pre *preparator.Preparator, prj *resource.Project) error {
	for _, branch := range prj.Branches {
		if branch.Ignored || !branch.FoundInsertTarget {
			continue
		}
		DecideBranchSources(pre, prj, branch.InsertTarget.Source, func(b *resource.Branch) {})
	}
	return nil
}

func DecideBranchSources(pre *preparator.Preparator, prj *resource.Project, insertTarget []string, more func(b *resource.Branch)) error {
	arrowLines := document.FindArrowLines(insertTarget)
	arrowPaths, whetherMust, protectedStrs := document.ParseArrowExps(insertTarget)
	for _, branch := range prj.Branches {
		if branch.Ignored || !branch.FoundInsertTarget || branch.MetInsertTargetRequirements {
			continue
		}

		branch.InsertTarget.Source = insertTarget
		branch.InsertTarget.SetPoints(arrowLines, arrowPaths, whetherMust, protectedStrs)
		if err := resource.GlobBranchSourceFilePaths(branch, pre.SuffixStrategies); err != nil {
			branch.Message.Failed(err.Error())
			continue
		}
		if err := decideBranchSourceStrategies(branch, pre.SuffixStrategies); err != nil {
			branch.Message.Failed(err.Error())			
			continue
		}
		branch.InsertTarget.Source = document.Copy(insertTarget)
		branch.MetInsertTargetRequirements = true
		branch.Message.Reset()
		more(branch)
	}
	return nil
}

func PullSourcesEachBranch(prj *resource.Project) error {
	for _, branch := range prj.Branches {
		if branch.Ignored || !branch.FoundInsertTarget || !branch.MetInsertTargetRequirements {
			continue
		}
		for pathToSource, source := range branch.Sources {
			source.Contents, _ = ioFile.Gets(pathToSource)
		}
	}
	return nil
}

func Arrow(pre *preparator.Preparator, prj *resource.Project, overwrite, practice bool) error {
	if practice {
		return nil
	}

	for _, branch := range prj.Branches {
		if branch.Ignored || !branch.FoundInsertTarget || !branch.MetInsertTargetRequirements {
			continue
		}

		if !overwrite {
			if newPath, err := protectCurrentInsertTargetFile(branch); err != nil {
				branch.Message.Failed(fmt.Sprintf("something why, cannot rename current \"%s\" \"%s\"", branch.InsertTarget.PathToInsertTarget, newPath))
				continue
			}	
		}

		if err := makeUpBranchSources(pre, branch); err != nil {
			continue
		}

		n, dn := 0, 0
		for _, point := range branch.InsertTarget.Points {
			branch.InsertTarget.Source = fileOpe.RemoveContent(branch.InsertTarget.Source, point.RequireLine+n)

			pathToSource := fmt.Sprintf("%s/%s", branch.PathToBranch, point.RequirePath)
			if _, exist := branch.Sources[pathToSource]; !exist {
				n += -1
				continue
			}

			if len(point.ProtectedStr) > 0 {
				branch.InsertTarget.Source = fileOpe.InsertContent(branch.InsertTarget.Source, point.ProtectedStr, point.RequireLine+n)
				continue
			}

			dn, branch.InsertTarget.Source = fileOpe.InsertContents(branch.InsertTarget.Source, branch.Sources[pathToSource].Contents, point.RequireLine+n)
			n += dn - 1
		}

		if err := ioFile.Puts(branch.InsertTarget.PathToInsertTarget, branch.InsertTarget.Source); err != nil {
			branch.Message.Failed(fmt.Sprintf("cannot write result into %s", branch.InsertTarget.PathToInsertTarget))
		}
	}
	return nil
}

func Return(pre *preparator.Preparator, prj *resource.Project, secret, quiet, practice bool) error {
	if quiet {
		return nil
	}

	res := result.Make(prj, practice)
	if err := result.Report(pre, res); err != nil {
		return err
	}
	if secret {
		return nil
	}
	return result.Record(pre, res)
}

func decideBranchSourceStrategies(br *resource.Branch, strategies preparator.SuffixPreferences) error {
	for pathToSource, source := range br.Sources {
		suffix := pathToSource[strings.LastIndex(pathToSource, "."):]
		if strategy, exist := strategies[suffix]; exist {
			source.Strategy = strategy
		} else if strategy, exist := strategies[".*"]; exist {
			source.Strategy = strategy
		} else {
			source.Strategy = strategies["*"]
		}
	}
	return nil
}

func protectCurrentInsertTargetFile(br *resource.Branch) (string, error) {
	if _, err := os.Stat(br.InsertTarget.PathToInsertTarget); os.IsNotExist(err) {
		return "", nil
	}

	splited := strings.Split(br.InsertTarget.PathToInsertTarget, "/")
	splited[len(splited)-1] = fmt.Sprintf("%s%s", PROTECT_PREFIX, splited[len(splited)-1])
	if err := os.Rename(br.InsertTarget.PathToInsertTarget, strings.Join(splited, "/")); err != nil {
		return strings.Join(splited, "/"), err
	}
	return strings.Join(splited, "/"), nil
}

func makeUpBranchSources(pre *preparator.Preparator, br *resource.Branch) error {
	for pathToSource, source := range br.Sources {
		if source.Strategy.UseEscape {
			source.Contents = fileOpe.EscapeLT(source.Contents)
		}
		if source.Strategy.UsePreLang {
			suffix := pathToSource[strings.LastIndex(pathToSource, "."):]
			source.Contents = fileOpe.Sand(source.Contents, fmt.Sprintf(config.PRE_TAG_PREFIX, pre.Cheiron.PreLangSuffixes[suffix]), config.PRE_TAG_SUFFIX)
		}
	}
	return nil
}
