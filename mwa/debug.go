package main

var IsDebugEnabled = false

func Debug(f func()) {
	if IsDebugEnabled {
		f()
	}
}
