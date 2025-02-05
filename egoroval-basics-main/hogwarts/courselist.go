//go:build !solution

package hogwarts

type Colour int

const (
	white = iota
	grey
	black
)

func DFS(prereqs map[string][]string, k string, way *[]string, col map[string]Colour) bool {
	col[k] = grey
	for _, v := range prereqs[k] {
		if col[v] == grey {
			return true
		}
		if col[v] == black {
			continue
		}
		if DFS(prereqs, v, way, col) {
			return true
		}
	}
	*way = append(*way, k)
	col[k] = black
	return false
}

func GetCourseList(prereqs map[string][]string) []string {
	col := make(map[string]Colour)
	way := make([]string, 0)
	p := &way
	for k := range prereqs {
		if col[k] == white {
			if DFS(prereqs, k, p, col) {
				panic("Cicle depending!")
			}
		}
	}
	return way
}
