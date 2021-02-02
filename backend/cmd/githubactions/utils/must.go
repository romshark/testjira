package utils

import "log"

// Must not be an error, otherwise we panic
func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
