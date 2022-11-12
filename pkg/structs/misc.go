package structs

// Ref returns reference to value
func Ref[V comparable](v V) *V {
	return &v
}

// Val returns value of reference
func Val[V comparable](v *V) V {
	return *v
}
