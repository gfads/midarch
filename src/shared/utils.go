package shared

import (
	"os"
	"strings"
)

func Contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func EnvironmentVariableValue(variable string) (value string) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == variable {
			value = pair[1]
		}
	}
	return value
}