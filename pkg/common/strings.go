package common

// PointerToString convert *string to string
func PointerToString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
