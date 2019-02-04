package main

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ExecuteScript(invokeScript string) error {
	log.Infof("Invoking /bin/sh %v", invokeScript)
	cmd := exec.Command("/bin/sh", invokeScript)
	cmdErr := cmd.Run()

	if cmdErr != nil {
		log.Warnf("Failed to execute script %v", cmdErr)
	} else {
		out, err := cmd.Output()
		if err != nil {
			log.Infof("Command output: %v", out)
		}
	}

	return cmdErr
}
