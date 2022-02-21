package multicmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/internal/arrow"
	"github.com/t-star08/cheiron/internal/options/arrowOptions"
	"github.com/t-star08/cheiron/internal/preparator"
	"github.com/t-star08/cheiron/internal/resource"
)

var CMD = &cobra.Command{
	Use:	"multi",
	Run:	run,
}

var (
	logger = log.New(os.Stderr, "multi: ", log.LstdFlags)

	opt	= arrowOptions.New()
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

func run(c *cobra.Command, args []string) {
	setPreWorker()
	if err := pre.Execute(args); err != nil {
		logger.Fatalln(err)
	}

	var project *resource.Project
	if prj, err := arrow.DetectProjectRoot(pre); err != nil {
		logger.Fatalln(err)
	} else {
		project = prj
	}

	if err := arrow.DetectBranchPaths(pre, project); err != nil {
		logger.Fatalln(err)
	}

	if err := arrow.DetectInsertTargetFilePaths(pre, project); err != nil {
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
}
