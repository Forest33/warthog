package structs

func If[V any](condition bool, then V, els V) V {
	if condition {
		return then
	}
	return els
}

func IfVal[V any](condition bool, then *V, els *V) V {
	if condition {
		return *then
	}
	return *els
}
