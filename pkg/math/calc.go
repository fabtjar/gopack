package calc

// Max returns the bigger of the 2 given values.
func Max[T int](a, b T) T {
	if a >= b {
		return a
	}
	return b
}
