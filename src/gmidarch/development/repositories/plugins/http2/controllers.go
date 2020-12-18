package http2

import "apps/http2server/impl"

func GetFunction(functionName string) interface{} {
	switch functionName {
	case "GetHealth": return impl.GetHealth
	case "GetItems": return impl.GetItems
	case "GetFibonacci": return impl.GetFibonacci
	default: return nil
	}
}