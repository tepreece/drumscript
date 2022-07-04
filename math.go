package main

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func LeastCommonMultiple(integers []int) int {
	switch len(integers) {
	case 0:
		return 0
	case 1:
		return integers[0]
	case 2:
		return LCM(integers[0], integers[1])
	default:
		return LCM(integers[0], integers[1], integers[2:]...)
	}
}
