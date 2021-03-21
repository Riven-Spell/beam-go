package to

func StringPtr(val string) *string {
	return &val
}

func BoolPtr(val bool) *bool {
	return &val
}

func Uint64Ptr(val uint64) *uint64 {
	return &val
}
