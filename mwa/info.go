package main

import "fmt"

var GitCommit string
var Version string
var BuildDate string

func ShowInfo() {

	if Version == "" {
		Version = "0.0.0"
	}

	if GitCommit == "" {
		GitCommit = "[Development build]"
	}

	fmt.Println("MWA Multi-Wireless-Agent")
	fmt.Println(fmt.Sprintf("VERSION: %v", Version))
	fmt.Println(fmt.Sprintf("GIT: %v", GitCommit))
}
