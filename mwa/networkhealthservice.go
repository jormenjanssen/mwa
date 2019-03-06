package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var NotReachableError = fmt.Errorf("Target not reachable")

const DelayForNextAttemptInRecovery = 1 * time.Second
const DelayBetweenAttemptsInVerify = 2 * time.Second
const VerifyAttempts = 3

type NetworkHealthService struct {
	Address        string
	HealthCheck    NetworkHealthCheck
	RecoveryTime   time.Duration
	RecoveryAction func() error
}

func NewNetworkHealthService(addr string, nhc NetworkHealthCheck, recvtime time.Duration, recva func() error) NetworkHealthService {
	return NetworkHealthService{Address: addr, HealthCheck: nhc, RecoveryTime: recvtime, RecoveryAction: recva}
}

func (nh NetworkHealthService) Verify(attempts int) error {
	log.Debugf("Invoking network healthcheck to: %v", nh.Address)
	return nh.TryVerifyMultipleAttempts(nh.VerifyOnce, VerifyAttempts, DelayBetweenAttemptsInVerify)
}

func (nh NetworkHealthService) Recover() error {
	return nh.RecoverWithinTime(time.Now())
}

func (nh NetworkHealthService) RecoverWithinTime(startTime time.Time) error {

	// We recover by a succesfull call or by throwing an error after we reached our timeout.
	for {

		recoveryDuration := time.Since(startTime)

		// We cannot recover by waiting, run our network recovery action.
		if recoveryDuration > nh.RecoveryTime {
			log.Warnf("Invoking network recovery action because [%v] exceeds maximum of [%v]", recoveryDuration, nh.RecoveryTime)
			return LastErrorFunc(nh.RecoveryAction, TimeOutError("Recovery-Time", recoveryDuration))
		}

		// We can recover by waiting, let our caller know we succeeded by just having some patience
		if nh.VerifyOnce() == nil {
			return nil
		}

		// Wait a short while
		<-time.After(DelayForNextAttemptInRecovery)
	}
}

func (nh NetworkHealthService) TryVerifyMultipleAttempts(f func() error, attempts int, delay time.Duration) error {

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

func (nh NetworkHealthService) VerifyOnce() error {

	if nh.HealthCheck == nil {
		return fmt.Errorf("Healthcheck not implemented")
	}

	return nh.HealthCheck.Perform(nh.Address)
}
