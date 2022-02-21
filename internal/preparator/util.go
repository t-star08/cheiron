package preparator

type SuffixPreferences map[string]*SuffixPreference

type SuffixPreference struct {
	Owner		string
	PreLang		string
	UseEscape	bool
	UsePreLang	bool
}

func newSuffixPreference(owener, preLang string, useEscape, usePreLang bool) *SuffixPreference {
	return &SuffixPreference {
		Owner: owener,
		PreLang: preLang,
		UseEscape: useEscape,
		UsePreLang: usePreLang,
	}
}

func (ss *SuffixPreferences) Has(s string) bool {
	_, exist := (*ss)[s]
	return exist
}
