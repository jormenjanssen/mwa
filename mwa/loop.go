package main

import (
	"log"
)

func KeepAlive(f func() error, r func(int, error) bool, disableOuterloop bool) error {

	if disableOuterloop {
		return innerloop(f, r)
	}

	return outerloop(f, r)
}

func outerloop(f func() error, r func(int, error) bool) error {

	log.Println("Running outer loop")

	// If we get a vallid response then start the innerloop handling.
	for {
		err := f()
		if err == nil {
			return innerloop(f, r)
		}
	}
}

func innerloop(f func() error, r func(int, error) bool) error {

	log.Println("Running inner loop")

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

	log.Println("Running error loop")
	errorCount := 1

	for {
		if !r(errorCount, err) {
			return err
		}
		errorCount++
	}
}
