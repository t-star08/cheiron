package resource

type Branch struct {
	PathToBranch				string
	Ignored						bool
	FoundInsertTarget			bool
	MetInsertTargetRequirements	bool
	InsertTarget				*InsertTarget
	Sources						map[string]*Source
	PathToBestTemplate			string
	Message						*Message
}

func newBranch(pathToBranch string) *Branch {
	return &Branch {
		PathToBranch: pathToBranch,
		Sources: make(map[string]*Source),
		Message: newMessage(),
	}
}
