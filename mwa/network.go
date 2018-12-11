package main

import (
	"fmt"
	"net/http"
	"time"
)

var TimeOutError = fmt.Errorf("Connection is not recovered within timely fashion")

type NetworkHealth struct {
	Address        string
	RecoveryTime   time.Duration
	RecoveryAction func() error
}

func (nh NetworkHealth) Verify() error {

	_, err := http.Get("http://www.google.nl/")
	fmt.Println(fmt.Printf("Invoking ping to: %v", nh.Address))
	//_, err := Ping(nh.Address)

	return err
}

func (nh NetworkHealth) Recover() error {
	return nh.RecoverWithinTime(time.Now())
}

func (nh NetworkHealth) RecoverWithinTime(startTime time.Time) error {

	// We recover by a succesfull call or by throwing an error after we reached our timeout.
	for {

		recoveryDuration := time.Since(startTime)

		// We cannot recover by waiting, run our network recovery action.
		if recoveryDuration > nh.RecoveryTime {
			return FuncError(nh.RecoveryAction, TimeOutError)
		}

		// We can recover by waiting, let our caller know we succeeded by just having some patience
		if nh.Verify() == nil {
			return nil
		}

		// Wait a short while
		<-time.After(1 * time.Second)
	}
}

func FuncError(f func() error, aerr error) error {

	if f == nil {
		return fmt.Errorf("No recovery action is defined we cannot recover")
	}

	err := f()

	if err == nil {
		err = aerr
	}

	return err
}
