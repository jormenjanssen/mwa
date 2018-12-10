package main

import (
	"log"
)

func KeepAlive(f func() error, r func(int, error) bool) error {
	return outerloop(f, r)
}

func outerloop(f func() error, r func(int, error) bool) error {

	// If we get a vallid response then start the innerloop handling.
	for {
		err := f()
		if err == nil {
			log.Println("Starting inner loop")
			return innerloop(f, r)
		}
	}
}

func innerloop(f func() error, r func(int, error) bool) error {

	// Innerloop is simple on first error start handling errors.
	for {
		err := f()
		if err != nil {
			// This is a helper arround
			return scopedErrorLoop(f, r, err)
		}
	}
}

func scopedErrorLoop(f func() error, r func(int, error) bool, err error) error {

	log.Println("Starting scoped error loop")
	errorCount := 1

	for {
		if !r(errorCount, err) {
			return err
		}
		errorCount++
	}
}
