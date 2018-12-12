package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var notReachableError = fmt.Errorf("Target not reachable")
var TimeOutError = fmt.Errorf("Connection is not recovered within timely fashion")

type NetworkHealth struct {
	Address        string
	RecoveryTime   time.Duration
	RecoveryAction func() error
}

func (nh NetworkHealth) Verify() error {
	log.Debugf("Invoking ping to: %v", nh.Address)
	return nh.TryVerifyMultipleAttempts(nh.VerifyOnce, 3, 2*time.Second)
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
			log.Warnf("Invoking network recovery action because [%v] exceeds maximum of [%v]", recoveryDuration, nh.RecoveryTime)
			return LastErrorFunc(nh.RecoveryAction, TimeOutError)
		}

		// We can recover by waiting, let our caller know we succeeded by just having some patience
		if nh.VerifyOnce() == nil {
			return nil
		}

		// Wait a short while
		<-time.After(1 * time.Second)
	}
}

func (nh NetworkHealth) TryVerifyMultipleAttempts(f func() error, attempts int, delay time.Duration) error {

	var err error = nil

	for i := 0; i < attempts; i++ {

		err = f()
		if err == nil {
			return nil
		}

		// Wait a short while
		<-time.After(delay)
	}

	return err
}

func (nh NetworkHealth) VerifyOnce() error {
	// TODO: remove this Debug Windows, refactor to +WIN BUILD only
	//_, err := http.Get("http://www.google.nl")

	up, err := Ping(nh.Address)

	// Extra check if we do not have an error but also not a response
	if !up && err == nil {
		return notReachableError
	}

	return err
}

// LastErrorFunc wraps a function where the result is always an error
// If the function thats being called is returning nil, then we return the alternate (aerr) error
func LastErrorFunc(f func() error, aerr error) error {

	if f == nil {
		return fmt.Errorf("No recovery action is defined we cannot recover")
	}

	err := f()

	if err == nil {
		err = aerr
	}

	return err
}
