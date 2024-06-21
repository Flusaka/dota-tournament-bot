package utils

func MapStructTo[I any, O any](input []I, mapFunc func(input I) O) []O {
	out := make([]O, len(input))
	for i, el := range input {
		out[i] = mapFunc(el)
	}
	return out
}

func FilterWhere[I any](input []I, predicate func(element I) bool) []I {
	out := make([]I, 0, len(input))
	for _, el := range input {
		if predicate(el) {
			out = append(out, el)
		}
	}
	return out
}
