package common

func Must[T any](val T, err any, messageArgs ...any) T {
	if err == nil {
		return val
	}

	panic(err)
}

func ToPtr[T any](val T) *T {
	return &val
}
