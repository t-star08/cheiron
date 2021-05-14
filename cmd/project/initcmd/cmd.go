package initcmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var CMD = &cobra.Command {
	Use: "init",
	Run: run,
}

type Options struct {
	force bool
}

type SummarizeArrow struct {
	Project_root	string					`json:"project_root"`
	Branch			string					`json:"branch"`
	Sources 		map[string]SourceDetail `json:"sources"`
	Insert_target 	string 					`json:"insert_target"`
}

type SourceDetail struct {
	Dire	string `json:"code_dire"`
	File	string `json:"code_file"`
	Key		string `json:"keyword"`
}

type SummarizeQuiver struct {
	Specify	[]string `json:"specify"`
	Ignore	[]string `json:"ignore"`
}

var (
	logger = log.New(os.Stderr, "init: ", log.LstdFlags)

	opt Options

	arrowJson = SummarizeArrow {
		".",
		".*",
		map[string]SourceDetail {
			"python3": {
				"code_python3",
				"main.py",
				"Python3",
			},
			"java": {
				"code_java",
				"Main.java",
				"Java",
			},
			"cpp": {
				"code_c-plus-plus",
				"main.cpp",
				"C++",
			},
			"cc": {
				"code_c-plus-plus",
				"main.cc",
				"C++",
			},
		},
		"DEFAULT.md",
	}

	quiverJson = SummarizeQuiver {
		make([]string, 0),
		make([]string, 0),
	}
)

func createSettingsDirectory() error {
	if _, err := os.Stat("./.cheiron"); !opt.force && os.IsExist(err) {
		return fmt.Errorf("[.cheiron] already exists\nif you wanna overwrite, use f option")
	}

	os.Mkdir("./.cheiron", 0777)

	return nil;
}

func createJson(source interface{}, name string) error {
	json_data, err := json.Marshal(source)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("./.cheiron/%s", name), json_data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func run(c *cobra.Command, args []string) {
	if err := createSettingsDirectory(); err != nil {
		logger.Fatalln(err)
	}

	if err := createJson(arrowJson, "arrow.json"); err != nil {
		logger.Fatalln(err)
	}

	if err := createJson(quiverJson, "quiver.json"); err != nil {
		logger.Fatalln(err)
	}

	fmt.Println("Created directory [./.cheiron]\nBy Editing json, you can add hook")
}

func init() {
	CMD.Flags().BoolVarP(&opt.force, "force", "f", false, "overwrite [arrow.json]")
}
