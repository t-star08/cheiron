package resource

type InsertTarget struct {
	PathToInsertTarget	string
	Source				[]string
	Points				[]*Point
}

type Point struct {
	RequireLine		int
	RequirePath		string
	WhetherMust		bool
	ProtectedStr	string
	Message			*Message
}

func newInsertTarget(pathToInsertTarget string) *InsertTarget {
	return &InsertTarget {
		PathToInsertTarget: pathToInsertTarget,
		Points: make([]*Point, 0),
	}
}

func newPoint(requireLine int, requirePath string, whetherMust bool, protectedStr string) *Point {
	return &Point {
		RequireLine: requireLine,
		RequirePath: requirePath,
		WhetherMust: whetherMust,
		ProtectedStr: protectedStr,
		Message: newMessage(),
	}
}

func CreateInsertTarget(newPathToInsertTarget string) *InsertTarget {
	return newInsertTarget(newPathToInsertTarget)
}

func (it *InsertTarget) SetPoints(requireLines []int, requirePaths []string, whetherMust []bool, protectedStrs []string) {
	it.Points = make([]*Point, len(requireLines))
	for i := 0; i < len(it.Points); i++ {
		it.Points[i] = newPoint(requireLines[i], requirePaths[i], whetherMust[i], protectedStrs[i])
	}
}
