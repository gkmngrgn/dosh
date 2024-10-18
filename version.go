package main

import "fmt"

const Version = "0.1.0"

func getVersion() string {
	return fmt.Sprintf("DOSH version v%s", Version)
}
