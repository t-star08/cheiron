package singlecmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/internal/arrow"
	"github.com/t-star08/cheiron/internal/options/arrowOptions/singleOptions"
	"github.com/t-star08/cheiron/internal/preparator"
	"github.com/t-star08/cheiron/internal/resource"
)

var CMD = &cobra.Command{
	Use:	"single",
	Run:	run,
}

var (
	logger = log.New(os.Stderr, "single: ", log.LstdFlags)

	opt	= singleOptions.New()
	pre = preparator.NewPreparator()
)

func setPreWorker() {
	pre.SetArgsChecker(
		func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("too few args")
			}
			return nil
		},
	)
}

func detectProjectRoot(pre *preparator.Preparator) (*resource.Project, error) {
	pre.Cheiron.ProjectRoot = "."
	return arrow.DetectProjectRoot(pre)
}

func detectBranchPathAndInsertTargetPath(prj *resource.Project) error {
	pre.Cheiron.Branch, pre.Cheiron.InsertTarget = splitIntoBranchAndInsertTarget()
	if err := resource.FindBranchPaths(prj, pre.Cheiron.Branch); err != nil {
		return err
	}
	if err := arrow.DetectInsertTargetFilePaths(pre, prj); err != nil {
		return err
	}

	return nil
}

func splitIntoBranchAndInsertTarget() (string, string) {
	border := strings.LastIndex(opt.Target, "/")
	var pathToBranch, pathToInsertTarget string
	if border == -1 {
		pathToBranch = "."
		pathToInsertTarget = opt.Target
	} else {
		pathToBranch = opt.Target[:border]
		pathToInsertTarget = opt.Target[border+1:]
	}

	return pathToBranch, pathToInsertTarget
}

func run(c *cobra.Command, args []string) {
	setPreWorker()
	if err := pre.Execute(args); err != nil {
		logger.Fatalln(err)
	}

	var project *resource.Project
	if prj, err := detectProjectRoot(pre); err != nil {
		logger.Fatalln(err)
	} else {
		project = prj
	}

	if err := detectBranchPathAndInsertTargetPath(project); err != nil {
		logger.Fatalln(err)
	}

	if err := arrow.PullInsertTargetSources(project); err != nil {
		logger.Fatalln(err)
	}

	if err := arrow.MeetInsertTargetRequirements(pre, project); err != nil {
		logger.Fatalln(err)
	}

	if err := arrow.PullSourcesEachBranch(project); err != nil {
		logger.Fatalln(err)
	}

	if err := arrow.Arrow(pre, project, opt.Overwrite, opt.Practice); err != nil {
		logger.Fatalln(err)
	}

	if err := arrow.Return(pre, project, opt.Secret, opt.Quiet, opt.Practice); err != nil {
		logger.Fatalln(err)
	}
}

func init() {
	CMD.Flags().BoolVarP(&opt.Practice, "practice", "p", false, "not insert, just check which is recognized or ignored")
	CMD.Flags().BoolVarP(&opt.Overwrite, "overwrite", "o", false, "not protect original insert target")
	CMD.Flags().BoolVarP(&opt.Secret, "secret", "s", false, "not create history file")
	CMD.Flags().BoolVarP(&opt.Quiet, "quiet", "q", false, "not leave any log")
	CMD.Flags().StringVarP(&opt.Target, "target", "t", "", "specify insert target")
}
