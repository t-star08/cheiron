package arwcmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"regexp"

	"github.com/spf13/cobra"
)

var CMD = &cobra.Command{
	Use:	"arrow",
	Run:	run,
}

var (
	logger = log.New(os.Stderr, "arrow: ", log.LstdFlags)

	from_json Summarize
	// candidate_directories
	cndt_dires []string
	re *regexp.Regexp

	hooks 	Hook
	opt		Options
)

type Options struct {
	force	bool
	key		bool
	simple	bool
}

type Hook struct {
	hooks	[]string
	max_len	int
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

func check_args(json map[string]SourceDetail, args []string) error {
	for _, hook := range args {
		if _, exist := json[hook]; exist {
			hooks.hooks = append(hooks.hooks, hook)
			if hooks.max_len < len(hook) {
				hooks.max_len = len(hook)
			}
		}
	}

	if len(hooks.hooks) == 0 {
		return fmt.Errorf("args are not registered in setting file [arrow.json]")
	}
	return nil
}

func list_candidate_dires() error {
	ls_dir, err := ioutil.ReadDir("./")
	if err != nil {
		return err
	}

	for _, info := range ls_dir {
		if info.IsDir() && re.MatchString(info.Name()) {
			cndt_dires = append(cndt_dires, info.Name())
		}
	}

	if len(cndt_dires) == 0 {
		return fmt.Errorf("no candidate directories")
	}

	return nil
}

func list_all_code_path(base, code_dir, code_file, hook string) ([]string, error) {
	rst := make([]string, 0)

	
	base = fmt.Sprintf("%s/%s", base, code_dir)
	code_path := fmt.Sprintf("%s/%s", base, code_file)
	_, err := os.Stat(code_path)
	if err != nil {
		return rst, fmt.Errorf("%*s: %s/%s: no code file", hooks.max_len, hook, from_json.Project_root, code_path)
	}
	rst = append(rst, code_path)

	if opt.simple {
		return rst, nil
	}

	i := 2
	for {
		code_path := fmt.Sprintf("%s_%d/%s", base, i, code_file)
		_, err := os.Stat(code_path)
		if err != nil {
			return rst, nil
		}
		rst = append(rst, code_path)

		i += 1
	}
}

func escape_lt(ans []string) []string {
	// Markdown で「<」が特殊文字として扱われてしまうので、エスケープする

	rst := make([]string, len(ans))

	for n, line := range ans {
		for i := 0; i < len(line); {
			if char := string(line[i]); char == "<" {
				i += 1
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

func find_pre_tag(source []string, keyword string) ([]int, []int, bool) {
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

	return pre_open, pre_close, len(pre_open) > 0
}

func find_keyword(source []string, keyword string) ([]int, bool) {
	key_point := make([]int, 0)

	for n, line := range source {
		if line == fmt.Sprintf("\\%s", keyword) {
			key_point = append(key_point, n)
		}
	}

	return key_point, len(key_point) > 0
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

func check_candidate_directories(args []string) error {
	n_cndt_dires := make([]string, 0)
	for _, candidate := range cndt_dires {
		skip := true
		for _, hook := range args {
			code_path := fmt.Sprintf("%s/%s/%s", candidate, from_json.Sources[hook].Dire, from_json.Sources[hook].File)
			if _, err := os.Stat(code_path); err == nil {
				skip = false
				break
			}
		}
		if !skip {
			target_path := fmt.Sprintf("%s/%s", candidate, from_json.Insert_target)
			if _, err := os.Stat(target_path); err != nil {
				skip = true
			}
		}
		if skip {
			logger.Printf("skip: %s", candidate)
		} else {
			n_cndt_dires = append(n_cndt_dires, candidate)
		}
	}

	cndt_dires = n_cndt_dires
	if len(cndt_dires) == 0 {
		return fmt.Errorf("no candidate directory")
	}

	return nil
}

func run(c *cobra.Command, args []string) {
	json_file, err := ioutil.ReadFile("../cheiron_settings/arrow.json")
	if err != nil {
		logger.Fatalln("setting file [arrow.json] does not exist\n\nTry: cheiron init")
	}

	err = json.Unmarshal(json_file, &from_json)
	if err != nil {
		logger.Fatalln(err)
	}

	if err := check_args(from_json.Sources, args); err != nil {
		logger.Fatalln(err)
	} else {
		fmt.Printf("valiable args: %s\n", strings.Join(hooks.hooks, ", "))
	}

	if err := os.Chdir(from_json.Project_root); err != nil {
		logger.Fatalln(err)
	}

	re, _ = regexp.Compile(`[\w\s]+`)
	if err := list_candidate_dires(); err != nil {
		logger.Fatalln(err)
	}

	if err := check_candidate_directories(args); err != nil {
		logger.Fatalln(err)
	}
	
	for _, candidate := range cndt_dires {
		target_source, err := readLine(fmt.Sprintf("%s/%s", candidate, from_json.Insert_target))
		if err != nil {
			logger.Fatalln(err)
		}

		for _, hook := range hooks.hooks {
			un_execute := false
			source_detail := from_json.Sources[hook]

			code_path, err := list_all_code_path(candidate, source_detail.Dire, source_detail.File, hook)
			if err != nil {
				logger.Println(err)
				continue
			}
			codes := make([][]string, 0)
			for _, path := range code_path {
				code, err := readLine(path)
				if err != nil {
					logger.Println(err)
					continue
				}
				code = escape_lt(code)
				codes = append(codes, code)
			}

			insert_point := make([]int, 0)
			pre_open := make([]int, 0)
			pre_close := make([]int, 0)
			key_point := make([]int, 0)
			found := false
			if !opt.key {
				pre_open, pre_close, found = find_pre_tag(target_source, source_detail.Key)
				if !found {
					key_point, found = find_keyword(target_source, source_detail.Key)
				}
			} else {
				key_point, found = find_keyword(target_source, source_detail.Key)
				if !found {
					pre_open, pre_close, found = find_pre_tag(target_source, source_detail.Key)
				}
			}
			
			if !found {
				un_execute = true
			} else {
				if len(pre_open) > 0 {
					if len(pre_open) != len(pre_close) {
						un_execute = true
					} else if opt.force || len(pre_open) == len(codes) {
						// 対象の pre タグの数と解答コードの数が一致しなくても解答コードの数だけ pre タグエリアの一旦削除を行う
						x := 0
						for i := range codes {
							if i > len(pre_open) - 1 {
								break
							}
							insert_point = append(insert_point, pre_open[i]-x)
							target_source = append(target_source[:pre_open[i]-x], target_source[pre_close[i]+1-x:]...)
							x += pre_close[i] - pre_open[i] + 1
						}
					} else {
						// 対象の pre タグの数と解答コードの数が一致しなければ pre タグエリアの一旦削除とその後の挿入は実行しない
						un_execute = true
					}
				} else  {
					// found(==True) なら必ず key_point は要素を持っている
					x := 0
					if opt.force || len(key_point) == len(codes) {
						for i := range codes {
							if i > len(key_point) - 1 {
								break
							}
							insert_point = append(insert_point, key_point[i]-x)
							target_source = append(target_source[:key_point[i]], target_source[key_point[i]+1:]...)
							x += 1
						}
					} else {
						un_execute = true
					}
				}
			}

			if un_execute {
				logger.Printf("%*s: did not execute %s/%s\n", hooks.max_len, hook, from_json.Project_root, candidate)
				continue
			}
			x := 0
			for i := range codes {
				if i > len(insert_point) - 1 {
					break
				}
				point := insert_point[i]
				code := codes[i]
				
				code = append([]string{fmt.Sprintf(`<pre lang="%s">`, source_detail.Key)}, code...)
				code = append(code, "</pre>")
				target_source = append(target_source[:point+x], append(code[:], target_source[point+x:]...)...)

				x += len(code)
			}
		}

		if err := writeLine(fmt.Sprintf("%s/%s", candidate, from_json.Insert_target), target_source); err != nil {
			logger.Println(err)
		}
	}

	fmt.Println("done")
}

func init() {
	CMD.Flags().BoolVarP(&opt.force, "force", "f", false, "insert forcibly")
	CMD.Flags().BoolVarP(&opt.key, "key", "k", false, "emphasis keyword when insert")
	CMD.Flags().BoolVarP(&opt.simple, "simple", "s", false, "simple directory search")
}
