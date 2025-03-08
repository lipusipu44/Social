package env

/*
Developed this env go which is to extract the string value
if not present fall back to default string.
this to be called in main.go class for value setup

Imp Note:- env.go is in internal package, because anything in internal
package can't be imported apart from files in this package or its parent
package like Social, cmd n all.

Ex: If we try to import env file in a package external/main.go it will
throw error. but it can be imported in cmd/api/main.go. by design go is like that
*/
import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		val, err := strconv.Atoi(value)
		if err != nil {
			return val
		}
	}
	return fallback
}
