package main

import log "github.com/sirupsen/logrus"

func NewScriptAction(script string) func() error {

	if script == "" {
		return NoScriptErrorFunc()
	}

	return func() error {
		return ExecuteScript(script)
	}
}

func NoScriptErrorFunc() func() error {
	return func() error {
		log.Warning("No script was configured to run")
		return nil
	}
}
