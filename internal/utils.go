package internal

import "log"

func ExitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
