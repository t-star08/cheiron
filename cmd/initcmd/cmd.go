package initcmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/internal/config"
	"github.com/t-star08/cheiron/internal/options/initOptions"
	"github.com/t-star08/cheiron/internal/result"
	"github.com/t-star08/hand/pkg/io/ioJson"
)

var CMD = &cobra.Command {
	Use: "init",
	Run: run,
}

var (
	logger = log.New(os.Stderr, "init: ", log.LstdFlags)

	opt = initOptions.New()
	cheironJson = config.TEMPLATE
)

func createConfDirectory() error {
	if _, err := os.Stat(config.CONF_DIR_NAME); !opt.Force && !os.IsNotExist(err) {
		if opt.Any {
			return nil
		}
		return fmt.Errorf("\"%s\" already exists\nif you wanna overwrite, use \"f\" option", config.CONF_DIR_NAME)
	}

	os.Mkdir(config.CONF_DIR_NAME, 0777)

	return nil;
}

func createHistoryDirectory() error {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", config.CONF_DIR_NAME, result.REPORT_DIR_NAME)); !opt.Force && !os.IsNotExist(err) {
		if opt.Any {
			return nil
		}
		return fmt.Errorf("\"%s/%s\" already exists\nif you wanna overwrite, use \"f\" option", config.CONF_DIR_NAME, result.REPORT_DIR_NAME)
	}
	
	os.Mkdir(fmt.Sprintf("%s/%s", config.CONF_DIR_NAME, result.REPORT_DIR_NAME), 0777)
	return nil
}

func createConfFile() error {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", config.CONF_DIR_NAME, config.CONF_FILE_NAME)); !opt.Force && !os.IsNotExist(err) {
		if opt.Any {
			return nil
		}
		return fmt.Errorf("\"%s/%s\" already exists\nif you wanna overwrite, use \"f\" option", config.CONF_DIR_NAME, config.CONF_FILE_NAME)
	}
	
	if err := ioJson.Puts(fmt.Sprintf("%s/%s", config.CONF_DIR_NAME, config.CONF_FILE_NAME), cheironJson); err != nil {
		return err
	}

	fmt.Printf("Created \"%s/%s\"\nBy Editing json, you can set config\n", config.CONF_DIR_NAME, config.CONF_FILE_NAME)
	return nil
}

func run(c *cobra.Command, args []string) {
	if err := createConfDirectory(); err != nil {
		logger.Fatalln(err)
	}

	if err := createConfFile(); err != nil {
		logger.Fatalln(err)
	}

	if err := createHistoryDirectory(); err != nil {
		logger.Fatalln(err)
	}
}

func init() {
	CMD.Flags().BoolVarP(&opt.Force, "force", "f", false, fmt.Sprintf("overwrite \"%s\"", config.CONF_FILE_NAME))
	CMD.Flags().BoolVarP(&opt.Any, "any", "a", false, "compensate lack of")
}
