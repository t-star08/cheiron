package config

type Cheiron struct {
	ProjectRoot		string					`json:"projectRoot"`
	Branch			string					`json:"branch"`
	Ignore			[]string				`json:"ignore"`
	InsertTarget 	string 					`json:"insertTarget"`
	Strategies		map[string]*Strategy	`json:"strategies"`
	StrategyAliases	map[string][]string		`json:"strategyAliases"`
	PreLangSuffixes	map[string]string		`json:"preLangSuffixes"`
	Routine			[]*Routine				`json:"routine"`
}

type Strategy struct {
	EscapeOpt		bool		`json:"useEscapeOption"`
	PreLangOpt		bool		`json:"usePreLangOption"`
	TargetSuffixes	[]string	`json:"targetSuffixes"`
}

type Routine struct {
	Template	string	`json:"template"`
	Priority	int		`json:"priority"`
}

const (
	CONF_DIR_NAME = ".cheiron"
	CONF_FILE_NAME = "cheiron.json"
	PRE_TAG_PREFIX = "<pre lang=\"%s\">"
	PRE_TAG_SUFFIX = "</pre>"
)

var (
	TEMPLATE = Cheiron {
		".",
		".*",
		[]string {
			"branch path written here be ignored",
		},
		"DEFAULT.md",
		map[string]*Strategy {
			"strategy1": {
				false,
				false,
				[]string {
					".py",
					".ruby",
				},
			},
			"strategy2": {
				true,
				false,
				[]string {
					".java",
					".cpp",
					".cc",
				},
			},
			"strategy3": {
				true,
				true,
				[]string {
					".*",
				},
			},
		},
		map[string][]string {
			"aliase":  {
				"strategy comb",
			},
			"aliase1": {
				"stratgey1",
				"strategy2",
			},
			"aliase2": {
				"strategy1",
				"strategy3",
			},
		},
		map[string]string {
			"suffix": "language",
			".py": "Python3",
			".java": "Java",
			".ruby": "Ruby",
			".cpp": "C++",
			".cc": "C++",
			".c": "C",
			".go": "GO",
		},
		[]*Routine {
			{
				"Path/to/template",
				0,
			},
			{
				"template.md",
				1,
			},
		},
	}
)
