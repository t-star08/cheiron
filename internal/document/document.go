package document

import (
	"strings"

	"github.com/t-star08/cheiron/pkg/fileOpe"
)

const (
	ARROW_SYMBOL = "<<<"
	IF_SYMBOL = "?"
)

func Copy(body []string) []string {
	p := make([]string, len(body))
	_ = copy(p, body)
	return p
}

func FindArrowLines(body []string) []int {
	return fileOpe.FindHeadSymbolLines(body, ARROW_SYMBOL)
}

func ParseArrowExps(body []string) ([]string, []bool, []string) {
	arrowLines := FindArrowLines(body)
	arrowPaths := make([]string, len(arrowLines))
	whetherMust := make([]bool, len(arrowLines))
	protectedStrs := make([]string, len(arrowLines))
	for i, line := range arrowLines {
		afterArrowSymbol := strings.TrimSpace(strings.Replace(body[line], ARROW_SYMBOL, "", -1))
		if ifSymbolPlace := strings.Index(afterArrowSymbol, IF_SYMBOL); ifSymbolPlace == -1 {
			arrowPaths[i] = afterArrowSymbol
			whetherMust[i] = true
			protectedStrs[i] = ""
		} else {
			arrowPaths[i] = afterArrowSymbol[:ifSymbolPlace]
			whetherMust[i] = false
			protectedStrs[i] = strings.TrimSpace(afterArrowSymbol[ifSymbolPlace+1:])
		}
	}

	return arrowPaths, whetherMust, protectedStrs
}
