package basis

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

//SuperFuncTestCase - стандартный набор тестов для тестирования
// переданного экземпляра SuperFuncType
// должен в обязательном порядке выполняться для всех реализаций
// используется достаточно просто:
// `func TestMySuperFunc(t *testing.T) { SuperFuncTestCase(MyFunc, t) }`
func SuperFuncTestCase(impl SuperFuncType, t *testing.T) {
	t.Run("`n==0` -> x1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			x1 := rand.Float64()
			x2 := rand.Float64()
			assert.Equal(t, x1, impl(x1, x2, 0))
		}
	})
	t.Run("`n==1` -> x1 * x2", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			x1 := rand.Float64()
			x2 := rand.Float64()
			assert.Equal(t, x1*x2, impl(x1, x2, 1))
		}
	})
	t.Run("`x1=1,x2=2,n=3` -> 0->1 1->2 2->2  3->4 ->4", func(t *testing.T) {
		assert.Equal(t, 4.0, impl(1, 2, 3))
	})
	t.Run("`x1=2,x2=3,n=5` -> 0->2 1->6 2->12  3->72 4->864  5->62208 -> 62208", func(t *testing.T) {
		assert.Equal(t, 62208.0, impl(2.0, 3.0, 5))
	})
	t.Run("`big(n) over x1 = 1.0001, x2 = 1.00002, n = 30 -> 4.917359272354959e+65", func(t *testing.T) {
		AssertIsValidSuperFuncF(t, 4.917359272354959e+65, 1.0001, 1.00002, 30, impl, DefaultPrecession)
	})
	basisRef := reflect.ValueOf(BasicSuperFuncImpl).Pointer()
	implRef := reflect.ValueOf(impl).Pointer()
	if implRef != basisRef {
		t.Run("`n>1` -> f(x1,x2,n-2) * f(x1,x2,n-1) (0.1% precession)", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				x1 := rand.Float64()
				x2 := rand.Float64()
				n := uint8(rand.Uint32() % 30)
				// сравнивать будем с эталонной базовой функцией
				AssertIsValidSuperFunc(t, x1, x2, n, impl, DefaultPrecession)
			}
		})
	}
}

type funcCase struct {
	x1 float64
	x2 float64
	n  uint8
}

const CASES_SIZE = 2048

func buildCases() []*funcCase {
	// функция вызывется только один раз на запуск - все получат одинаковый набор кейсов
	var seed int64 = rand.Int63()
	rnd := rand.New(rand.NewSource(seed))
	var result []*funcCase
	for i := 0; i < CASES_SIZE; i++ {
		x1 := rnd.Float64()
		x2 := rnd.Float64()
		var n uint8
		up_or_down := rnd.Float64()
		preN := rand.Uint32()
		// гораздо чаще будут попадаться кейсы 20 - 30 и реже 2 - 19
		if up_or_down > 0.2 {
			n = uint8((preN % 11) + 20) // 20-30
		} else {
			n = uint8((preN % 18) + 2) // 2-19
		}
		result = append(result, &funcCase{x1, x2, n})
	}
	return result
}

var cases []*funcCase = buildCases()

// SuperFuncBenchmark - обобщенный бенчмарк для SuperFuncType
// использование в своем коде:
// `func BenchmarkMySuperFunc(b *testing.B) { SuperFuncBenchmark(MyFunc, b) }`
func SuperFuncBenchmark(impl SuperFuncType, b *testing.B) {
	for i := 0; i < b.N; i++ {
		fCase := cases[i%CASES_SIZE]
		impl(fCase.x1, fCase.x2, fCase.n)
	}
}

//DefaultPrecession - максимальная допустимая погрешность по умолчанию - 0.1%
var DefaultPrecession = 0.001

// AssertIsValidSuperFunc - вспомогательное утверждение для сравнения результатов
// float64, использует IsEqualWithPrecession для сравнения с допустимой погрешностью
// для сравнения использует вычисляемый эталон на основе базовой функции
func AssertIsValidSuperFunc(t *testing.T, x1 float64, x2 float64, n uint8, impl SuperFuncType, precession float64) {
	reference := BasicSuperFuncImpl(x1, x2, n)
	actual := impl(x1, x2, n)
	if !IsEqualWithPrecession(reference, actual, precession) {
		assert.Failf(t, "Функция не соответствует требованию расчета с погрешностью не более 0.01% от эталона",
			"x1: %f, x2: %f, n:%d, reference: %f, actual: %f", x1, x2, n, reference, actual)
	}
}

// AssertIsValidSuperFuncF - вспомогательное утверждение для сравнения результатов
// float64, использует IsEqualWithPrecession для сравнения с допустимой погрешностью
// для сравнения использует заранее вычисленное значение [reference]
func AssertIsValidSuperFuncF(t *testing.T, reference float64, x1 float64, x2 float64, n uint8, impl SuperFuncType, precession float64) {
	actual := impl(x1, x2, n)
	if !IsEqualWithPrecession(reference, actual, precession) {
		assert.Failf(t, "Функция не соответствует требованию расчета с погрешностью не более 0.01% от эталона",
			"x1: %f, x2: %f, n:%d, reference: %f, actual: %f", x1, x2, n, reference, actual)
	}
}

// IsEqualWithPrecession IsEqualWithDelta вспомогательная функция для сравнения float64
// с учетом допустимой точности вычислений относительно эталона,
// например погрешность 1% `IsEqualWithPrecession(x1,x2,0.01)`
func IsEqualWithPrecession(reference float64, actual float64, precession float64) bool {
	return math.Abs(reference-actual) <= (math.Abs(reference) * precession)
}

const DefaultNForBenchmark = 30

// SuperFuncBenchmark - обобщенный бенчмарк для SuperFuncType
// использование в своем коде:
// `func BenchmarkMySuperFunc(b *testing.B) { SuperFuncBenchmark(MyFunc, b) }`
func OldHackedSuperFuncBenchmark(impl SuperFuncType, b *testing.B) {
	b.StopTimer()
	var xs [1024]float64
	for i := 0; i < 512; i++ {
		xs[i*2] = rand.Float64()
		xs[i*2+1] = rand.Float64()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		basis := i % 512
		impl(xs[basis*2], xs[basis*2+1], DefaultNForBenchmark)
	}
}
