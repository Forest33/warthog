package structs

func Ref[V comparable](v V) *V {
	return &v
}

func Val[V comparable](v *V) V {
	return *v
}
