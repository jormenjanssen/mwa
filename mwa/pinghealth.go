package main

type PingHealthCheck struct {
}

func (p PingHealthCheck) Perform(addr string) error {

	up, err := Ping(addr)

	// Extra check if we do not have an error but also not a response
	if !up && err == nil {
		return NotReachableError
	}

	return err
}
