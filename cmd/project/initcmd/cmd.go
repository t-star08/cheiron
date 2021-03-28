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

type Summarize struct {
	Project_root	string					`json:"project_root"`
	Sources 		map[string]SourceDetail `json:"sources"`
	Insert_target 	string 					`json:"insert_target"`
}

type SourceDetail struct {
	Dire	string `json:"code_dire"`
	File	string `json:"code_file"`
	Key		string `json:"keyword"`
}

var (
	logger = log.New(os.Stderr, "init: ", log.LstdFlags)

	opt Options

	lang_json = Summarize {
		".",
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
)

func create_json(data Summarize) error {
	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	if _, err := os.Stat("../cheiron_settings"); os.IsNotExist(err) {
		os.Mkdir("../cheiron_settings", 0777)
	}

	if _, err := os.Stat("../cheiron_settings/arrow.json"); !opt.force && !os.IsNotExist(err) {
		return fmt.Errorf("[arrow.json] already exists\nif you wanna overwrite, use f option")
	}
	
	err = ioutil.WriteFile("../cheiron_settings/arrow.json", json_data, 0777)
	if err != nil {
		return err
	}

	return nil
}

func run(c *cobra.Command, args []string) {
	err := create_json(lang_json)
	if err != nil {
		logger.Fatalln(err)
	}
	fmt.Println("Created directory [../cheiron_settings/arrow.json]\nYou can edit it to add target language")
}

func init() {
	CMD.Flags().BoolVarP(&opt.force, "force", "f", false, "overwrite [arrow.json]")
}
