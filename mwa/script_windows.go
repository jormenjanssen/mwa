package main

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ExecuteScript(invokeScript string) error {
	log.Infof("Invoking cmd.exe %v", invokeScript)
	cmd := exec.Command("cmd.exe /c", invokeScript)
	cmdErr := cmd.Run()

	if cmdErr != nil {
		log.Warnf("Failed to execute script %v", cmdErr)
	}

	return cmdErr
}
