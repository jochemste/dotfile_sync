package utils

import "log"

func CheckIfError(err error) {
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}
