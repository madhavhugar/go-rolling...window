package main

import "log"

func handleOnError(err error, message string) {
	if err != nil {
		log.Println("$$ ERROR:", message, err)
	}
}
