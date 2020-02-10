package components

import (
	"reflect"
)

var ProxyLibrary = map[string] reflect.Type{
	reflect.TypeOf(CalculatorProxy{}).String(): reflect.TypeOf(CalculatorProxy{}),
	reflect.TypeOf(Fibonacciproxy{}).String():  reflect.TypeOf(Fibonacciproxy{})}
