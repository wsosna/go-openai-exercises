package utils

import "log"

func HandleFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
