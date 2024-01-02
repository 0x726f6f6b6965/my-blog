package helper

func Intersect[T string | int](a []T, b []T) []T {
	set := make([]T, 0)
	hash := make(map[T]bool)
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
