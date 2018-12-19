package main

import (
	"flag"
	"os"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {

	recoveryTime := flag.Duration("recoverytime", 60*time.Second, "The seconds to wait before executing the recovery action")
	host := flag.String("host", "192.168.10.2", "The default host to check against")

	script := flag.String("script", "", "The command to execute when we are failed")
	test := flag.Bool("test", false, "Check the configured action")
	debug := flag.Bool("debug", false, "Enable debugging")
	isMonitor := flag.Bool("monitor", false, "Only monitor do not execute script")
	gatewayTarget := flag.String("default-ipv4-gateway-for-adapter", "", "The default gateway to use for host selection")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)

	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Preset out host to host
	targetHost := *host

	if *gatewayTarget != "" {

		timeout := 3 * time.Minute
		log.Infof("Target host selection is set to autodetect (using default gateway) for adapter: %v with a maximum detection timeout of: %v", *gatewayTarget, timeout)

		ip, err := GetIpv4TargetForAdapterGatewayWithTimeout(targetHost, timeout)

		if err != nil {
			log.Fatalf("Could not autodetect gateway error: %v took: %v", err, timeout)
		}

		// Assign our target host we discovered
		targetHost = ip.String()
	}

	nh := NetworkHealth{Address: targetHost, RecoveryTime: *recoveryTime}

	// Hookup the recovery action
	if *script != "" && !*isMonitor {
		log.Infof("The error action is configured to /bin/sh %v", *script)
		nh.RecoveryAction = func() error {
			return executeScript(*script)
		}
	} else if !*isMonitor {
		log.Fatalf("The error action is not configured use (-script <file>) or (-monitor)")
	} else if *isMonitor && *test {
		log.Fatalf("Cannot use (-monitor) together with (-test). Use (-test -script <file>)")
	}

	// Test execution for console
	// Todo: Move to cobra sfp
	if *test {
		if nh.RecoveryAction != nil {
			nh.RecoveryAction()
		}
		os.Exit(0)
	}

	// Log to our users
	log.Debugf("Logger debug enabled: %v", *debug)
	log.Printf("Running network health against target: %v with recovery time: [%v]", nh.Address, nh.RecoveryTime)

	// Run Application
	appCtx := NewApplicationContext()
	Watchdog(appCtx, nh, nh)
}

func executeScript(invokeScript string) error {

	log.Infof("Invoking /bin/sh %v", invokeScript)
	cmd := exec.Command("/bin/sh", invokeScript)
	cmdErr := cmd.Run()

	if cmdErr != nil {
		log.Warnf("Failed to execute script %v", cmdErr)
	}

	return cmdErr
}
