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
		AssertIsValidSuperFuncF(t, 4.917359272354959e+65, 1.0001, 1.00002, 30, impl, DEFAULT_PRECESSION)
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
				AssertIsValidSuperFunc(t, x1, x2, n, impl, DEFAULT_PRECESSION)
			}
		})
	}
}

//SuperFuncBenchmark - обобщенный бенчмарк для SuperFuncType
// использованиек в свооем коде:
//  `func BenchmarkMySuperFunc(b *testing.B) { SuperFuncBenchmark(MyFunc, b) }`
func SuperFuncBenchmark(impl SuperFuncType, b *testing.B) {
	b.StopTimer()
	var seria [1024]float64
	for i := 0; i < 512; i++ {
		seria[i*2] = rand.Float64()
		seria[i*2+1] = rand.Float64()
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		basis := i % 512
		impl(seria[basis*2], seria[basis*2+1], DEFAULT_N_FOR_BENCHMARK)
	}

}

//DEFAULT_PRECESSION - максимальная допустимая погрешность по умолчанию - 0.1%
var DEFAULT_PRECESSION float64 = 0.001

//DEFAULT_N_FOR_BENCHMARK - n c которым функция проверяется при выполнении бенчмарка
var DEFAULT_N_FOR_BENCHMARK uint8 = 30

//AssertIsValidSuperFunc - вспомогательное утверждение для сравнения результатов
//float64, использует IsEqualWithPrecession для сравнения с допустимой погрешностью
// для сравнения использует вычисляемый эталон на основе базовой функции
func AssertIsValidSuperFunc(t *testing.T, x1 float64, x2 float64, n uint8, impl SuperFuncType, precession float64) {
	etalon := BasicSuperFuncImpl(x1, x2, n)
	actual := impl(x1, x2, n)
	if !IsEqualWithPrecession(etalon, actual, precession) {
		assert.Failf(t, "Функция не соответствует требованию расчета с погрешностью не более 0.01% от эталона",
			"x1: %f, x2: %f, n:%d, etalon: %f, actual: %f", x1, x2, n, etalon, actual)
	}
}

//AssertIsValidSuperFunc - вспомогательное утверждение для сравнения результатов
//float64, использует IsEqualWithPrecession для сравнения с допустимой погрешностью
// для сравнения использует заранее вычисленное значение [reference]
func AssertIsValidSuperFuncF(t *testing.T, reference float64, x1 float64, x2 float64, n uint8, impl SuperFuncType, precession float64) {
	actual := impl(x1, x2, n)
	if !IsEqualWithPrecession(reference, actual, precession) {
		assert.Failf(t, "Функция не соответствует требованию расчета с погрешностью не более 0.01% от эталона",
			"x1: %f, x2: %f, n:%d, reference: %f, actual: %f", x1, x2, n, reference, actual)
	}
}

//IsEqualWithDelta вспомогательная функция для сравнения float64
//с учетом допустимой точности вычислений относительно эталона,
//например погрешность 1% `IsEqualWithPrecession(x1,x2,0.01)`
func IsEqualWithPrecession(etalon float64, actual float64, precession float64) bool {
	return math.Abs(etalon-actual) <= (math.Abs(etalon) * precession)
}
