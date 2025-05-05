package slice

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

func Sum[T Number](s []T) T {
	var result T
	for _, v := range s {
		result += v
	}
	return result
}
