package helpers

func SliceContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SliceContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
