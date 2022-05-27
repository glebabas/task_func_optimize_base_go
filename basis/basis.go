package basis

//SuperFuncType - сигнатура типа для реализуемой функции
type SuperFuncType func(x1 float64, x2 float64, n uint8) float64

//BasicSuperFuncImpl - начальная не оптимизированная версия функции
// описывает основной инвариантный алгоритм:
// 1. `n==0` -> `x1`
// 2. `n==1` -> `x1 * x2`
// 3. `n>1` -> `f(x1, x2, n-2) * f(x1, x2, n-1)`
func BasicSuperFuncImpl(x1 float64, x2 float64, n uint8) float64 {
	switch n {
	case 0:
		return x1
	case 1:
		return x1 * x2
	default:
		return BasicSuperFuncImpl(x1, x2, n-2) * BasicSuperFuncImpl(x1, x2, n-1)
	}
}
