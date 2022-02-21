package resource

type Project struct {
	PathToProjectRoot	string
	Branches			map[string]*Branch
	PathToTemplates		[]string
	Message				string
}

func newProject() *Project {
	return &Project {
		Branches: make(map[string]*Branch),
	}
}
