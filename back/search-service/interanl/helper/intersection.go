package helper

func Intersect(a []string, b []string) []string {
	set := make([]string, 0)
	hash := make(map[string]bool)
	for _, val := range a {
		hash[val] = true
	}
	for _, val := range b {
		if hash[val] {
			set = append(set, val)
		}
	}
	return set
}
