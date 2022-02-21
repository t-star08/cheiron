package fileOpe

import (
	"fmt"
	"unicode"

	"github.com/t-star08/hand/pkg/evaluator/stringValue"
)

var (
	evaluator = stringValue.NewPartialEvaluator()
)

func FindHeadSymbolLines(contents []string, symbol string) []int {
	res := make([]int, 0)
	symbol = fmt.Sprintf("^%s", symbol)
	for line, content := range contents {
		if found, err := evaluator.Evaluate(symbol, content); err != nil {
			continue
		} else if found {
			res = append(res, line)
		}
	}
	return res
}

func EscapeLT(contents []string) []string {
	res := make([]string, len(contents))
	for line, content := range contents {
		for i := 0; i < len(content); {
			if char := string(content[i]); char == "<" {
				if nextChar := string(content[i+1]); unicode.IsLetter(rune(content[i+1])) || nextChar == "!" || nextChar == "?"  {
					content = content[:i] + "&lt;" + content[i+1:]
					i += len("&lt;")
				} else {
					i += 1
				}
			} else {
				i += 1
			}
		}
		res[line] = content
	}

	return res
}

func Sand(contents []string, prefix, suffix string) []string {
	if prefix != "" {
		contents = InsertContent(contents, prefix, 0)
	}
	if suffix != "" {
		contents = append(contents, suffix)
	}
	return contents
}

func RemoveContent(origin []string, line int) []string {
	return append(origin[:line], origin[line+1:]...)
}

func InsertContent(origin []string, content string, line int) []string {
	return append(origin[:line], append([]string{content}, origin[line:]...)...)
}

func InsertContents(origin, contents []string, line int) (int, []string) {
	return len(contents), append(origin[:line], append(contents[:], origin[line:]...)...)
}
