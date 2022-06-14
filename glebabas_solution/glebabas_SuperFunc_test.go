package my_solution

import "testing"
import "github.com/comdiv/task_func_optimize_base_go/basis"

//BenchmarkBasicSuperFuncImpl оставлен для наглядного сравнения производительности
func BenchmarkBasicSuperFuncImpl(b *testing.B) {
	basis.SuperFuncBenchmark(basis.BasicSuperFuncImpl, b)
}

func TestMySuperFuncImpl(t *testing.T) {
	basis.SuperFuncTestCase(MySuperFuncImpl, t)
}

func BenchmarkMySuperFuncImpl(b *testing.B) {
	basis.SuperFuncBenchmark(MySuperFuncImpl, b)
}
