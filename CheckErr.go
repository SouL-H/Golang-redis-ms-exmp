package main

import "log"

func CheckErr(err error) {
	if err != nil {
		log.Fatalf("Failed: %v", err)

	}
}
