package helpers

import (
	"runtime"
	"strings"
)

func GetCurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	strs := strings.Split((runtime.FuncForPC(pc).Name()), "/")
	return strs[len(strs)-1]
}

func PointerTo[T any](value T) *T { return &value }

func UnwrapOrDefault[T any](value *T, base T) T {
	if value == nil {
		return base
	}
	return *value
}
