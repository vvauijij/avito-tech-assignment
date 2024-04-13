package env

import "os"

var env = os.Getenv("ENV")

func IsTestENV() bool {
	return env == "TEST"
}
