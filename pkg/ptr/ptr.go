package ptr

func Ref[T any](d T) *T {
	return &d
}
