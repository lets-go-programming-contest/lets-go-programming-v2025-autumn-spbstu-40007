package functionals

func Iter[X any](f func(X), xs []X) {
	for _, x := range xs {
		f(x)
	}
}

func Unique[X comparable](xs []X) []X {
	uniqueXs := make([]X, 0, len(xs))
	set := make(map[X]bool)

	for _, x := range xs {
		if !set[x] {
			set[x] = true
			uniqueXs = append(uniqueXs, x)
		}
	}

	return uniqueXs
}
