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

	//delay := flag.Duration("delay", 30*time.Second, "The duration to delay after an succesfull attempt")
	recoveryTime := flag.Duration("recoverytime", 60*time.Second, "The seconds to wait before executing the recovery action")
	//retries := flag.Int("retries", 3, "The number of retries before giving up")
	host := flag.String("host", "192.168.10.2", "The default host to check against")
	//disableOuterLoop := flag.Bool("without-check", false, "Disables pre checking")
	invokeScript := flag.String("script", "", "The command to execute when we are failed")
	testExec := flag.Bool("exectest", false, "Check the configured action")

	flag.Parse()

	if *testExec {

		if *invokeScript == "" {
			log.Fatalf("script is not configured")
		}

		executeScript(*invokeScript)
		os.Exit(0)
	}

	nh := NetworkHealth{Address: *host, RecoveryTime: *recoveryTime}
	appCtx := NewApplicationContext()
	Watchdog(appCtx, nh, nh)

	log.Println(fmt.Printf("Running network health against target: %v with recovery time: %v", nh.Address, nh.RecoveryTime))

	if *invokeScript != "" {
		fmt.Println(fmt.Printf("The error action is configured to /bin/sh %v", invokeScript))
	}

}

func executeScript(invokeScript string) {

	log.Println(fmt.Printf("Invoking /bin/sh %v", invokeScript))
	cmd := exec.Command("/bin/sh", invokeScript)
	cmdErr := cmd.Run()

	if cmdErr != nil {
		log.Println(fmt.Printf("Failed to execute script %v", cmdErr))
	}
}
