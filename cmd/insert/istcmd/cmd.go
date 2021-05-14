package istcmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/spf13/cobra"
)

var CMD = &cobra.Command{
	Use:	"insert",
	Run:	run,
}

var (
	logger = log.New(os.Stderr, "insert: ", log.LstdFlags)

	opt Options
)

type Options struct {
	keyword	string
	force	bool
	key		bool
}

func escapeLt(ans []string) []string {
	// commentary.md で「<」が特殊文字として扱われてしまうので、エスケープする

	rst := make([]string, len(ans))

	for n, line := range ans {
		for i := 0; i < len(line); {
			if char := string(line[i]); char == "<" {
				i += 1
				// どの文字が「<」直後に現れると「<」が特殊文字として働くか網羅できてない気がする
				if nextChar := string(line[i]); unicode.IsLetter(rune(line[i])) || nextChar == "!" || nextChar == "?"  {
					line = line[:i-1] + "&lt;" + line[i:]
					i += 3
				}
			} else {
				i += 1
			}
		}
		rst[n] = line
	}

	return rst
}

func findPreTag(source []string, keyword string) ([]int, []int, bool, bool) {
	preOpen := make([]int, 0)
	preClose := make([]int, 0)
	preOpening := false

	for n, line := range source {
		if line == fmt.Sprintf(`<pre lang="%s"></pre>`, keyword) {
			preOpen = append(preOpen, n)
			preClose = append(preClose, n)
		} else if line == fmt.Sprintf(`<pre lang="%s">`, keyword) {
			preOpen = append(preOpen, n)
			preOpening = true
		} else if preOpening && line == "</pre>" {
			preClose = append(preClose, n)
			preOpening = false
		}
	}

	return preOpen, preClose, len(preOpen) > 0, len(preOpen) > 1
}

func findKeyword(source []string, keyword string) ([]int, bool, bool) {
	keyPoint := make([]int, 0)

	for n, line := range source {
		if line == fmt.Sprintf("\\%s", keyword) {
			keyPoint = append(keyPoint, n)
		}
	}

	return keyPoint, len(keyPoint) > 0, len(keyPoint) > 1
}

func readLine(path string) ([]string, error) {
	doc := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		return doc, err
	}

	defer file.Close()

	scan := bufio.NewScanner(file)
	for scan.Scan() {
		line := scan.Text()
		doc = append(doc, line)
	}

	return doc, nil
}

func writeLine(path string, source []string) error {
	comFile, err := os.Create(path)
	if err != nil {
		return nil
	}

	defer comFile.Close()

	writer := bufio.NewWriter(comFile)
	for _, line := range(source) {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	
	writer.Flush()

	return nil
}

func run(c *cobra.Command, args []string) {
	if opt.keyword == "" {
		logger.Fatalln("must set keyword to use insert")
	}

	codePath := args[0]
	targetPath := args[1]
	fmt.Printf("%6s: %s\n", "code", codePath)
	fmt.Printf("%6s: %s\n", "target", targetPath)

	targetSource, err := readLine(targetPath)
	if err != nil {
		logger.Fatalln(err)
	}

	code, err := readLine(codePath)
	if err != nil {
		logger.Fatalln(err)
	}
	code = escapeLt(code)

	insertPoint := 0
	preOpen := make([]int, 0)
	preClose := make([]int, 0)
	keyPoint := make([]int, 0)
	found := false
	multiFound := false
	if !opt.key {
		preOpen, preClose, found, multiFound = findPreTag(targetSource, opt.keyword)
		if !found {
			keyPoint, found, multiFound = findKeyword(targetSource, opt.keyword)
		}
	} else {
		keyPoint, found, multiFound = findKeyword(targetSource, opt.keyword)
		if !found {
			preOpen, preClose, found, multiFound = findPreTag(targetSource, opt.keyword)
		}
	}
	
	if !found {
		logger.Fatalln("insert point not found")
	}
	if !opt.force && multiFound {
		logger.Fatalln("insert point found more than 2")
	}
	if len(preOpen) > 0 {
		if len(preOpen) != len(preClose) {
			logger.Fatalln("maybe pre tag not closed")
		}
		insertPoint = preOpen[0]
		targetSource = append(targetSource[:preOpen[0]], targetSource[preClose[0]+1:]...)
	} else {
		insertPoint = keyPoint[0]
		targetSource = append(targetSource[:keyPoint[0]], targetSource[keyPoint[0]+1:]...)
	}

	code = append([]string{fmt.Sprintf(`<pre lang="%s">`, opt.keyword)}, code...)
	code = append(code, "</pre>")
	targetSource = append(targetSource[:insertPoint], append(code[:], targetSource[insertPoint:]...)...)

	if err := writeLine(targetPath, targetSource); err != nil {
		logger.Fatalln(err)
	}

	fmt.Println("done")
}

func init() {
	CMD.Flags().StringVarP(&opt.keyword, "keyword", "", "", "insertion keyword")
	CMD.Flags().BoolVarP(&opt.force, "force", "f", false, "insert forcibly")
	CMD.Flags().BoolVarP(&opt.key, "key", "k", false, "emphasis keyword when insert")
}
