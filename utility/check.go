package utility

func InStringSlice(value string, strSlice []string) bool {
	for _, s := range strSlice {
		if value == s {
			return true
		}
	}
	return false
}
