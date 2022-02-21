package result

import (
	"fmt"
	"sort"

	"github.com/t-star08/cheiron/internal/preparator"
	"github.com/t-star08/cheiron/internal/resource"
	"github.com/t-star08/hand/pkg/io/ioJson"
	"github.com/t-star08/hand/pkg/io/ioStd"
)

func Make(prj *resource.Project, practice bool) *Res {
	res := newRes(practice)
	for _, branch := range prj.Branches {
		if !branch.Ignored && branch.FoundInsertTarget && branch.MetInsertTargetRequirements {
			rb := &RBranch {
				PathToBranch: branch.PathToBranch,
				Template: branch.PathToBestTemplate,
				Cuz: branch.Message.Cuz,
			}
			for _, point := range branch.InsertTarget.Points {
				if point.Message.Result {
					rb.Met = append(rb.Met, &RSource{point.RequirePath, branch.Sources[fmt.Sprintf("%s/%s", branch.PathToBranch, point.RequirePath)].Strategy.Owner, point.Message.Cuz})
				} else {
					rb.UnMet = append(rb.UnMet, &RSource{point.RequirePath, "-", point.Message.Cuz})
				}
			}
			res.Recognized = append(res.Recognized, rb)
		} else {
			res.Ignored = append(
				res.Ignored,
				&RBranch {
					PathToBranch: branch.PathToBranch,
					Cuz: branch.Message.Cuz,
				},
			)
		}
	}

	sort.Slice(
		res.Recognized,
		func(i, j int) bool {
			return res.Recognized[i].PathToBranch < res.Recognized[j].PathToBranch
		},
	)

	return res
}

func Report(pre *preparator.Preparator, res *Res) error {
	succeeded, failed := make([]string, len(res.Recognized)), make([]string, len(res.Ignored))
	for i := 0; i < len(succeeded); i++ {
		succeeded[i] = res.Recognized[i].PathToBranch
	}
	for i := 0; i < len(failed); i++ {
		failed[i] = res.Ignored[i].PathToBranch
	}

	if res.Performance {
		fmt.Println("MODE: PERFORMANCE")
	} else {
		fmt.Println("MODE: PRACTICE")
	}
	fmt.Println("Recognized (succeeded)")
	ioStd.MonospacedPuts("  ", succeeded)
	fmt.Println("Ignored (failed)")
	ioStd.MonospacedPuts("  ", failed)

	return nil
}

func Record(pre *preparator.Preparator, res *Res) error {
	historyTitle := fmt.Sprintf("%s/%s/%s.json", pre.PathToConfigDir, REPORT_DIR_NAME, REPORT_FILE_NAME)
	if err := ioJson.Puts(historyTitle, res); err != nil {
		return err
	}
	fmt.Printf("details are shown in \"%s\"\n", historyTitle)
	return nil
}
