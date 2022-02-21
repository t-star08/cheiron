package preparator

import (
	"fmt"

	"github.com/t-star08/cheiron/internal/config"
)

type Preparator struct {
	Cheiron					*config.Cheiron
	PathToConfigDir			string
	SuffixStrategies		SuffixPreferences
	CheckArgs				func(args []string) error
}

func NewPreparator() *Preparator {
	return &Preparator {
		SuffixStrategies: make(SuffixPreferences),
	} 
}

func (pre *Preparator) SetArgsChecker(checker func(args []string) error) {
	pre.CheckArgs = checker
}

func (pre *Preparator) Execute(args []string) error {
	if path, conf, err := readConf(); err != nil {
		return err
	} else {
		pre.PathToConfigDir = path
		pre.Cheiron = conf
	}

	if err := pre.CheckArgs(args); err != nil {
		return err
	}

	if err := pre.confirmStrategy(args); err != nil {
		return err
	}

	return nil
}

func (pre *Preparator) integrateSuffixes(strategyName string) error {
	if strategy, exist := pre.Cheiron.Strategies[strategyName]; !exist {
		return nil
	} else {
		for _, suffix := range strategy.TargetSuffixes {
			if _, exist := pre.SuffixStrategies[suffix]; exist {
				continue
			}
			pre.SuffixStrategies[suffix] = newSuffixPreference(strategyName, pre.Cheiron.PreLangSuffixes[suffix], strategy.EscapeOpt, strategy.PreLangOpt)
			if suffix == ".*" || suffix == "*" {
				return fmt.Errorf("appear wild card")
			}
		}
	}
	return nil
}

func (pre *Preparator) confirmStrategy(args []string) error {
	done := false
	for _, requiredStrategyAliase := range args {
		if strategyNames, exist := pre.Cheiron.StrategyAliases[requiredStrategyAliase]; exist {
			for _, strategyName := range strategyNames {
				if err := pre.integrateSuffixes(strategyName); err != nil {
					done = true
					break
				}
			}
		} else {
			if err := pre.integrateSuffixes(requiredStrategyAliase); err != nil {
				done = true
				break
			}
		}
		if done {
			break
		}
	}

	if len(pre.SuffixStrategies) == 0 {
		return fmt.Errorf("all args are not registered in \"%s\"", config.CONF_FILE_NAME)
	}
	return nil
}
