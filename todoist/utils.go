package todoist

// p<Type> functions are used to convert basic types to their pointer equivalents
func pString(s string) *string {
	return &s
}

func pInt64(i int64) *int64 {
	return &i
}
