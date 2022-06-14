package my_solution

// MySuperFuncImpl MySuperFunc - реализация
// 1. `n==0` -> `x1`
// 2. `n==1` -> `x1 * x2`
// 3. `n>1` -> `f(x1, x2, n-2) * f(x1, x2, n-1)`
func MySuperFuncImpl(x1 float64, x2 float64, n uint8) float64 {
	// изначально текст полностью повторяет BasicSuperFuncImpl
	switch n {
	case 0:
		return x1
	case 1:
		return x1 * x2
	default:
		return MySuperFuncImpl(x1, x2, n-2) * MySuperFuncImpl(x1, x2, n-1)
	}
}
