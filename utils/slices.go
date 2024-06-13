package utils

func MapStructTo[I any, O any](input []I, mapFunc func(input I) O) []O {
	out := make([]O, len(input))
	for i, el := range input {
		out[i] = mapFunc(el)
	}
	return out
}
