package basis

import (
	"testing"
)

func TestBasicSuperFuncImpl(t *testing.T) {
	SuperFuncTestCase(BasicSuperFuncImpl, t)
}

func BenchmarkBasicSuperFuncImpl(b *testing.B) {
	SuperFuncBenchmark(BasicSuperFuncImpl, b)
}
