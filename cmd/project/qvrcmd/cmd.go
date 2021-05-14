package qvrcmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"regexp"

	"github.com/spf13/cobra"
	"github.com/t-star08/cheiron/cmd/project/arwcmd"
)

var CMD = &cobra.Command{
	Use:	"quiver",
	Run:	run,
}

var (
	logger = log.New(os.Stderr, "quiver: ", log.LstdFlags)

	quiver Quiver
	re *regexp.Regexp
)

type Quiver struct {
	Specify	[]string `json:"specify"`
	Ignore	[]string `json:"ignore"`
}

func contains(key string, target []string) bool {
	for _, ele := range target {
		if ele == key {
			return true
		}
	}
	return false
}

func allSubDir() ([]string, error) {
	lsDir, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	re, _ = regexp.Compile(`^[\w\s]+$`)
	targetDir := make([]string, 0)
	for _, info := range lsDir {
		name := info.Name()
		if info.IsDir() && re.MatchString(name) {
			if contains(name, quiver.Ignore) {
				continue
			}
			if _, err := os.Stat(fmt.Sprintf("./%s/.cheiron", name)); os.IsNotExist(err) {
				logger.Printf("skip: %s", name)
				continue
			}
			targetDir = append(targetDir, name)
		}
	}
	return targetDir, nil
}

func specifiedSubDir() ([]string, error) {
	lsDir, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	targetDir := make([]string, 0)
	for _, info := range lsDir {
		name := info.Name()
		if info.IsDir() && contains(name, quiver.Specify) {
			if contains(name, quiver.Ignore) {
				continue
			}
			if _, err := os.Stat(fmt.Sprintf("./%s/.cheiron", name)); os.IsNotExist(err) {
				logger.Printf("skip: %s", name)
				continue
			}
			targetDir = append(targetDir, name)
		}
	}

	return targetDir, nil
}

func run(c *cobra.Command, args []string) {
	jsonFile, err := ioutil.ReadFile("./.cheiron/quiver.json")
	if err != nil {
		logger.Fatalln("[.cheiron] does not exist\n\nTry: cheiron project init")
	}
	
	if err := json.Unmarshal(jsonFile, &quiver); err != nil {
		logger.Fatalln(err)
	}

	var targetDir []string
	if len(quiver.Specify) == 0 {
		if targetDir, err = allSubDir(); err != nil {
			logger.Fatalln(err)
		}
	} else {
		if targetDir, err = specifiedSubDir(); err != nil {
			logger.Fatalln(err)
		}
	}

	initialPath, err := os.Getwd()
	if err != nil {
		logger.Fatalln(err)
	}

	for _, name := range targetDir {
		fmt.Println(name)
		os.Chdir(name)
		arwcmd.Run(c, args)
		os.Chdir(initialPath)
	}
}