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

type Lang struct {
	Ans_dire	string `json:"ans_dire"`
	Ans_file	string `json:"ans_file"`
	In_com		string `json:"in_com"`
}

func escape_lt(ans []string) []string {
	// commentary.md で「<」が特殊文字として扱われてしまうので、エスケープする

	rst := make([]string, len(ans))

	for n, line := range ans {
		for i := 0; i < len(line); {
			if char := string(line[i]); char == "<" {
				i += 1
				// どの文字が「<」直後に現れると「<」が特殊文字として働くか網羅できてない気がする
				if next_char := string(line[i]); unicode.IsLetter(rune(line[i])) || next_char == "!" || next_char == "?"  {
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

func find_pre_tag(source []string, keyword string) ([]int, []int, bool, bool) {
	pre_open := make([]int, 0)
	pre_close := make([]int, 0)
	pre_opening := false

	for n, line := range source {
		if line == fmt.Sprintf(`<pre lang="%s"></pre>`, keyword) {
			pre_open = append(pre_open, n)
			pre_close = append(pre_close, n)
		} else if line == fmt.Sprintf(`<pre lang="%s">`, keyword) {
			pre_open = append(pre_open, n)
			pre_opening = true
		} else if pre_opening && line == "</pre>" {
			pre_close = append(pre_close, n)
			pre_opening = false
		}
	}

	return pre_open, pre_close, len(pre_open) > 0, len(pre_open) > 1
}

func find_keyword(source []string, keyword string) ([]int, bool, bool) {
	key_point := make([]int, 0)

	for n, line := range source {
		if line == fmt.Sprintf("\\%s", keyword) {
			key_point = append(key_point, n)
		}
	}

	return key_point, len(key_point) > 0, len(key_point) > 1
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
	com_file, err := os.Create(path)
	if err != nil {
		return nil
	}

	defer com_file.Close()

	writer := bufio.NewWriter(com_file)
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

	code_path := args[0]
	target_path := args[1]
	fmt.Printf("%6s: %s\n", "code", code_path)
	fmt.Printf("%6s: %s\n", "target", target_path)

	target_source, err := readLine(target_path)
	if err != nil {
		logger.Fatalln(err)
	}

	code, err := readLine(code_path)
	if err != nil {
		logger.Fatalln(err)
	}
	code = escape_lt(code)

	insert_point := 0
	pre_open := make([]int, 0)
	pre_close := make([]int, 0)
	key_point := make([]int, 0)
	found := false
	multi_found := false
	if !opt.key {
		pre_open, pre_close, found, multi_found = find_pre_tag(target_source, opt.keyword)
		if !found {
			key_point, found, multi_found = find_keyword(target_source, opt.keyword)
		}
	} else {
		key_point, found, multi_found = find_keyword(target_source, opt.keyword)
		if !found {
			pre_open, pre_close, found, multi_found = find_pre_tag(target_source, opt.keyword)
		}
	}
	
	if !found {
		logger.Fatalln("insert point not found")
	}
	if !opt.force && multi_found {
		logger.Fatalln("insert point found more than 2")
	}
	if len(pre_open) > 0 {
		if len(pre_open) != len(pre_close) {
			logger.Fatalln("maybe pre tag not closed")
		}
		insert_point = pre_open[0]
		target_source = append(target_source[:pre_open[0]], target_source[pre_close[0]+1:]...)
	} else {
		insert_point = key_point[0]
		target_source = append(target_source[:key_point[0]], target_source[key_point[0]+1:]...)
	}

	code = append([]string{fmt.Sprintf(`<pre lang="%s">`, opt.keyword)}, code...)
	code = append(code, "</pre>")
	target_source = append(target_source[:insert_point], append(code[:], target_source[insert_point:]...)...)

	if err := writeLine(target_path, target_source); err != nil {
		logger.Fatalln(err)
	}

	fmt.Println("done")
}

func init() {
	CMD.Flags().StringVarP(&opt.keyword, "keyword", "", "", "insertion keyword")
	CMD.Flags().BoolVarP(&opt.force, "force", "f", false, "insert forcibly")
	CMD.Flags().BoolVarP(&opt.key, "key", "k", false, "emphasis keyword when insert")
}
