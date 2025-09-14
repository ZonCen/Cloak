package helpers

import "fmt"

var (
	Verbose bool
)

func LogVerbose(msg string, args ...interface{}) {
	if Verbose {
		fmt.Printf(msg+"\n", args...)
	}
}
