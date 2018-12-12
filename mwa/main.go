package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	recoveryTime := flag.Duration("recoverytime", 60*time.Second, "The seconds to wait before executing the recovery action")
	host := flag.String("host", "192.168.10.2", "The default host to check against")

	invokeScript := flag.String("script", "", "The command to execute when we are failed")
	testExec := flag.Bool("exectest", false, "Check the configured action")
	debug := flag.Bool("debug", false, "Enable debugging")

	flag.Parse()

	IsDebugEnabled = *debug

	if *testExec {

		if *invokeScript == "" {
			log.Fatalf("script is not configured")
		}

		executeScript(*invokeScript)
		os.Exit(0)
	}

	nh := NetworkHealth{Address: *host, RecoveryTime: *recoveryTime}

	log.Println(fmt.Printf("Running network health against target: %v with recovery time: %v", nh.Address, nh.RecoveryTime))
	Debug(func() { log.Println("#Debug is enabled") })

	if *invokeScript != "" {
		fmt.Println(fmt.Printf("The error action is configured to /bin/sh %v", *invokeScript))
		nh.RecoveryAction = func() error {
			return executeScript(*invokeScript)
		}
	}

	appCtx := NewApplicationContext()
	Watchdog(appCtx, nh, nh)
}

func executeScript(invokeScript string) error {

	log.Println(fmt.Printf("Invoking /bin/sh %v", invokeScript))
	cmd := exec.Command("/bin/sh", invokeScript)
	cmdErr := cmd.Run()

	if cmdErr != nil {
		log.Println(fmt.Printf("Failed to execute script %v", cmdErr))
	}

	return cmdErr
}
