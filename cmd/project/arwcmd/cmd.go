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
	Run:	Run,
}

var (
	logger = log.New(os.Stderr, "arrow: ", log.LstdFlags)

	arrow Arrow

	branchPath		[]string
	codeDirPath		[]string
	codeFilePath	[]string
	targetFilePath	string

	opt		Options
	hooks 	Hook

	regexpShortCut map[string]string = map[string]string {
		"$A": `^[\w\s]+$`,
		"$B": `^[\w]+$`,
		"$C": ".*",
		"$D": "^[0-9]",
		"$E": ".",
	}
)

type Options struct {
	force	bool
	key		bool
	simple	bool
}

type Hook struct {
	hooks	[]string
	maxLen	int
}

type Arrow struct {
	ProjectRoot		string					`json:"project_root"`
	Branch			string					`json:"branch"`
	Sources 		map[string]SourceDetail `json:"sources"`
	InsertTarget 	string 					`json:"insert_target"`
}

type SourceDetail struct {
	Dire	string `json:"code_dire"`
	File	string `json:"code_file"`
	Key		string `json:"keyword"`
}

func checkArgs(sources map[string]SourceDetail, args []string) error {
	hooks.hooks = make([]string, 0)
	
	for _, hook := range args {
		if _, exist := sources[hook]; exist {
			hooks.hooks = append(hooks.hooks, hook)
			if hooks.maxLen < len(hook) {
				hooks.maxLen = len(hook)
			}
		}
	}

	if len(hooks.hooks) == 0 {
		return fmt.Errorf("args not registered in setting file [arrow.json]")
	}
	return nil
}

func visit(path string, depth int) ([]string, error) {
	spot := make([]string, 0)

	sep := strings.Split(path, "/")
	if depth == len(sep) {
		spot = append(spot, path)
		return spot, nil
	}

	head := sep[depth]
	if head == "" {
		spot = append(spot, strings.Join(sep[:depth], "/"))
		return spot, nil
	}


	if info, err := os.Stat(head); !os.IsNotExist(err) {
		if info.IsDir() {
			wd, _ := os.Getwd()
			defer os.Chdir(wd)

			os.Chdir(head)
			if p, err := visit(strings.Join(sep, "/"), depth+1); err != nil {
				return spot, err
			} else {
				return append(spot, p...), nil
			}
		} else if depth == len(sep)-1 {
			return append(spot, path), nil
		}
	}


	// 完全一致
	re, err := regexp.Compile(head)
	if err != nil {
		return spot, err
	}
	if exp, exist := regexpShortCut[head]; exist {
		re, err = regexp.Compile(exp)
		if err != nil {
			return spot, err
		}
	}

	lsDir, err := ioutil.ReadDir(".")
	if err != nil {
		return spot, err
	}

	for _, info := range lsDir {
		name := info.Name()
		if re.MatchString(name) && info.IsDir() {
			os.Chdir(name)

			sep[depth] = name
			if p, err := visit(strings.Join(sep, "/"), depth+1); err != nil {
				fmt.Println(err)
			} else {
				spot = append(spot, p...)
			}

			os.Chdir("..")
		}
	}

	return spot, nil
}

func detectBranchPath(path string) error {
	if p, err := visit(path, 0); err != nil {
		return err
	} else {
		if len(p) == 0 {
			wd, _ := os.Getwd()
			return fmt.Errorf("no file matched [%s/%s]", wd, path)
		}
		branchPath = p
		return nil
	}
}

func listCodeDirPath(path string) error {
	if p, err := visit(path, 0); err != nil {
		return err
	} else {
		if len(p) == 0 {
			wd, _ := os.Getwd()
			return fmt.Errorf("no file matched [%s/%s]", wd, path)
		}
		codeDirPath = p
		return nil
	}
}

func AdditionalListCodePath() {
	if opt.simple {
		return
	}
	rst := make([]string, 0)
	for _, path := range codeDirPath {
		rst = append(rst, path)
		i := 2
		for {
			if info, err := os.Stat(fmt.Sprintf("%s_%d", path, i)); !os.IsNotExist(err) && info.IsDir() {
				rst = append(rst, fmt.Sprintf("%s_%d", path, i))
			} else {
				break
			}
			i += 1
		}
	}
	codeDirPath = rst
}

func listCodeFilePath(codeFile string) error {
	rst := make([]string, 0)
	re, err := regexp.Compile(codeFile)
	if err != nil {
		return err
	}

	for _, path := range codeDirPath {
		lsDir, _ := ioutil.ReadDir(path)
		if _, err := os.Stat(fmt.Sprintf("%s/%s", path, codeFile)); !os.IsNotExist(err) {
			rst = append(rst, fmt.Sprintf("%s/%s", path, codeFile))
			continue
		}
		for _, info := range lsDir {
			name := info.Name()
			if !info.IsDir() && re.MatchString(name) {
				rst = append(rst, fmt.Sprintf("%s/%s", path, name))
			}
		}
	}
	
	codeFilePath = rst

	return nil
}

func detectTargetFile(path string) error {
	if p, err := visit(path, 0); err != nil {
		return err
	} else {
		if len(p) == 0 {
			wd, _ := os.Getwd()
			return fmt.Errorf("%s/%s does not exist", wd, path)
		}
		targetFilePath = p[0]
		return nil
	}
}

func escapeLt(ans []string) []string {
	// Markdown で「<」が特殊文字として扱われてしまうので、エスケープする

	rst := make([]string, len(ans))

	for n, line := range ans {
		for i := 0; i < len(line); {
			if char := string(line[i]); char == "<" {
				i += 1
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

func findPreTag(source []string, keyword string) ([]int, []int, bool) {
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

	return preOpen, preClose, len(preOpen) > 0
}

func findKeyword(source []string, keyword string) ([]int, bool) {
	keyPoint := make([]int, 0)

	for n, line := range source {
		if line == fmt.Sprintf("\\%s", keyword) {
			keyPoint = append(keyPoint, n)
		}
	}

	return keyPoint, len(keyPoint) > 0
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
	targetFile, err := os.Create(path)
	if err != nil {
		return nil
	}

	defer targetFile.Close()

	writer := bufio.NewWriter(targetFile)
	for _, line := range(source) {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	
	writer.Flush()

	return nil
}

func showSuccessBranch(successBranch map[string][]string) {
	fmt.Println("Success Branch")
	if len(successBranch) == 0 {
		fmt.Println("  None")
		fmt.Println()
		return
	}

	for bPath, hooks := range successBranch {
		fmt.Printf("  [%s]\n    args: %s\n", bPath, strings.Join(hooks, ", "))
	}
	fmt.Println()
}

func showFailureBranch(failureBranch map[string][]string) {
	fmt.Println("Failure Branch")
	if len(failureBranch) == 0 {
		fmt.Println("  None")
		fmt.Println()
		return
	}

	for bPath, hooks := range failureBranch {
		fmt.Printf("  [%s]\n    args: %s\n", bPath, strings.Join(hooks, ", "))
	}
	fmt.Println()
}

func Run(c *cobra.Command, args []string) {
	jsonFile, err := ioutil.ReadFile("./.cheiron/arrow.json")
	if err != nil {
		logger.Fatalln("setting file [arrow.json] does not exist\n\nTry: cheiron project init")
	}

	if err := json.Unmarshal(jsonFile, &arrow); err != nil {
		logger.Fatalln(err)
	}
	
	if err := checkArgs(arrow.Sources, args); err != nil {
		logger.Fatalln(err)
	} else {
		fmt.Printf("valiable args: %s\n", strings.Join(hooks.hooks, ", "))
	}

	if initialPath, err := os.Getwd(); err != nil {
		logger.Fatalln(err)
	} else {
		defer os.Chdir(initialPath)
	}

	if err := os.Chdir(arrow.ProjectRoot); err != nil {
		logger.Fatalln(err)
	}

	if err := detectBranchPath(arrow.Branch); err != nil {
		logger.Fatalln(err)
	}

	successBranch := make(map[string][]string)
	failureBranch := make(map[string][]string)
	for _, bPath := range branchPath {
		if err := detectTargetFile(fmt.Sprintf("%s/%s", bPath, arrow.InsertTarget)); err != nil {
			failureBranch[bPath] = hooks.hooks
			continue
		}
		targetSource, err := readLine(targetFilePath)
		if err != nil {
			logger.Println(err)
			continue
		}

		for _, hook := range hooks.hooks {
			unExecute := false
			sourceDetail := arrow.Sources[hook]

			if err := listCodeDirPath(fmt.Sprintf("%s/%s", bPath, sourceDetail.Dire)); err != nil {
				failureBranch[bPath] = append(failureBranch[bPath], hook)
				continue
			}

			AdditionalListCodePath()

			if err := listCodeFilePath(sourceDetail.File); err != nil {
				failureBranch[bPath] = append(failureBranch[bPath], hook)
				continue
			}

			codes := make([][]string, 0)
			for _, path := range codeFilePath {
				code, err := readLine(path)
				if err != nil {
					logger.Println(err)
					continue
				}
				code = escapeLt(code)
				codes = append(codes, code)
			}

			insertPoint := make([]int, 0)
			preOpen := make([]int, 0)
			preClose := make([]int, 0)
			keyPoint := make([]int, 0)
			found := false
			if !opt.key {
				preOpen, preClose, found = findPreTag(targetSource, sourceDetail.Key)
				if !found {
					keyPoint, found = findKeyword(targetSource, sourceDetail.Key)
				}
			} else {
				keyPoint, found = findKeyword(targetSource, sourceDetail.Key)
				if !found {
					preOpen, preClose, found = findPreTag(targetSource, sourceDetail.Key)
				}
			}
				
			if !found {
				unExecute = true
			} else {
				if len(preOpen) > 0 {
					if opt.force || len(preOpen) == len(codes) {
						// 対象の pre タグの数と解答コードの数が一致しなくても解答コードの数だけ pre タグエリアの一旦削除を行う
						x := 0
						for i := range codes {
							if i > len(preOpen) - 1 {
								break
							}
							insertPoint = append(insertPoint, preOpen[i]-x)
							targetSource = append(targetSource[:preOpen[i]-x], targetSource[preClose[i]+1-x:]...)
							x += preClose[i] - preOpen[i] + 1
						}
					} else {
						// 対象の pre タグの数と解答コードの数が一致しなければ pre タグエリアの一旦削除とその後の挿入は実行しない
						unExecute = true
					}
				} else  {
					// found(==True) なら必ず key_point は要素を持っている
					x := 0
					if opt.force || len(keyPoint) == len(codes) {
						for i := range codes {
							if i > len(keyPoint) - 1 {
								break
							}
							insertPoint = append(insertPoint, keyPoint[i]-x)
							targetSource = append(targetSource[:keyPoint[i]], targetSource[keyPoint[i]+1:]...)
							x += 1
						}
					} else {
						unExecute = true
					}
				}
			}

			if unExecute {
				failureBranch[bPath] = append(failureBranch[bPath], hook)
				// logger.Printf("%*s: did not execute %s/%s\n", hooks.maxLen, hook, arrow.ProjectRoot, bPath)
				continue
			}
			x := 0
			for i := range codes {
				if i > len(insertPoint) - 1 {
					break
				}
				point := insertPoint[i]
				code := codes[i]
				
				code = append([]string{fmt.Sprintf(`<pre lang="%s">`, sourceDetail.Key)}, code...)
				code = append(code, "</pre>")
				targetSource = append(targetSource[:point+x], append(code[:], targetSource[point+x:]...)...)

				x += len(code)
			}
			successBranch[bPath] = append(successBranch[bPath], hook)
		}

		if err := writeLine(targetFilePath, targetSource); err != nil {
			logger.Println(err)
		}
	}

	showSuccessBranch(successBranch)
	showFailureBranch(failureBranch)

	fmt.Println("done")
}

func init() {
	CMD.Flags().BoolVarP(&opt.force, "force", "f", false, "insert forcibly")
	CMD.Flags().BoolVarP(&opt.key, "key", "k", false, "emphasis keyword when insert")
	CMD.Flags().BoolVarP(&opt.simple, "simple", "s", false, "simple directory search")
}
