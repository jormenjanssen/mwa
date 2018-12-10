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

	log.Println("Entering [OUTER]")

	// If we get a vallid response then start the innerloop handling.
	for {
		err := f()
		if err == nil {
			return innerloop(f, r)
		}
	}
}

func innerloop(f func() error, r func(int, error) bool) error {

	log.Println("Entering [INNER]")

	// Innerloop is simple on first error start handling errors.
	for {
		err := f()
		if err != nil {
			// This is a helper arround
			err = scopedErrorLoop(f, r, err)
			if err != nil {
				return err
			} else {
				log.Println("Entering [INNER]")
			}
		}
	}
}

func scopedErrorLoop(f func() error, r func(int, error) bool, err error) error {

	log.Println("Entering [ERROR]")
	errorCount := 1

	for {

		if !r(errorCount, err) {
			log.Println("Unrecoverable leaving [ERROR]")
			return err
		}

		err = f()

		if err == nil {
			log.Println("Gracefully Leaving [ERROR]")
			return nil
		}

		errorCount++
	}
}
