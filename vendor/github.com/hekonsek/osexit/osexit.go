package osexit

import (
	"fmt"
	"os"
	"runtime/debug"
)

const UnixExitCodeOK = 0

const UnixExitCodeGeneralError = 1

func ExitOnError(err error) {
	if err != nil {
		fmt.Printf("Something went wrong: %s\n", err)
		if os.Getenv("OSEXIT_DEBUG") == "true" {
			debug.PrintStack()
		}
		os.Exit(UnixExitCodeGeneralError)
	}
}

func ExitBecauseError(errorMessage string) {
	fmt.Println(errorMessage)
	os.Exit(UnixExitCodeGeneralError)
}
