package result

import (
	"fmt"
	"time"
)

const (
	REPORT_DIR_NAME = "history"
)

type Res struct {
	Performance	bool		`json:"performance"`
	Recognized	[]*RBranch	`json:"recognized"`
	Ignored		[]*RBranch	`json:"ignored"`
}

type RBranch struct {
	PathToBranch	string		`json:"pathToBranch"`
	Met				[]*RSource	`json:"met"`
	UnMet			[]*RSource	`json:"unMet"`
	Template		string		`json:"usedTemplate"`
	Cuz				string		`json:"cuz"`
}

type RSource struct {
	PathToSource	string	`json:"pathToSource"`
	Strategy		string	`json:"usedStrategy"`
	Cuz				string	`json:"cuz"`
}

var (
	t = time.Now()
	REPORT_FILE_NAME = fmt.Sprintf("%d-%02d-%02d-%02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
)

func newRes(practice bool) *Res {
	return &Res {
		Performance: !practice,
	}
}
