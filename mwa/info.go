package main

import (
	"fmt"
	"runtime"
)

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
	fmt.Println(fmt.Sprintf("GOLANG RUNTIME: %v", runtime.Version()))
	fmt.Println(fmt.Sprintf("VERSION: %v", Version))
	fmt.Println(fmt.Sprintf("GIT: %v", GitCommit))
}
