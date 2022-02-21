package statuscmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/internal/config"
	"github.com/t-star08/cheiron/internal/result"
)

const (
	NO_DATA = "-- no data --"
)

var CMD = &cobra.Command{
	Use:	"status",
	Run:	run,
}

var (
	logger = log.New(os.Stderr, "status: ", log.LstdFlags)
)

func checkArgs(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too much args")
	}
	return nil
}

func getWhenLastUsed(pathToConfDir string) string {
	pathToHistory := fmt.Sprintf("%s/%s", pathToConfDir, result.REPORT_DIR_NAME)
	if lsDir, err := os.ReadDir(pathToHistory); err != nil || len(lsDir) == 0 {
		return NO_DATA
	} else {
		lastUse := strings.Replace(lsDir[len(lsDir)-1].Name(), ".json", "", 1)
		if m := strings.LastIndex(lastUse, "-"); m == -1 {
			return NO_DATA
		} else {
			return fmt.Sprintf("%s %s", strings.Replace(lastUse[:m], "-", "/", -1), lastUse[m+1:])
		}
	}
}

func showStatus(pathToConfDir string) {
	fmt.Printf("CONFIG: %s/%s\n", pathToConfDir, config.CONF_FILE_NAME)
	fmt.Printf("HISTORY: %s/%s\n", pathToConfDir, result.REPORT_DIR_NAME)
	fmt.Printf("LAST USED: %s\n", getWhenLastUsed(pathToConfDir))
}

func run(c *cobra.Command, args []string) {
	if err := checkArgs(args); err != nil {
		logger.Fatalln(err)
	}

	if pathToConfDir, err := config.SearchConfigDirPath(); err != nil {
		logger.Fatalln(err)
	} else {
		showStatus(pathToConfDir)
	}
}
