package utils

func Contains(slc []string, item string) int {
	for i, v := range slc {
		if v == item {
			return i
		}
	}
	return -1;
}
