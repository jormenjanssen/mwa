package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const MinimalRecoveryTime = 10 * time.Second

type Config struct {
	RecoveryTime        string
	MonitorOnly         bool
	Host                string
	AutoIpv4GatewayHost bool
}

func JsonConfigFromReader(r io.Reader) (Config, error) {
	var cfg Config
	err := json.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, cfg.Validate()
}

//Validate validates the config for basic validation errors
func (c *Config) Validate() error {

	d, err := c.RecoveryDuration()
	if err != nil {
		return err
	}

	if d < MinimalRecoveryTime {
		return fmt.Errorf("Configured recovery time of %v is < %v treshold", d, MinimalRecoveryTime)
	}

	if c.Host == "" && !c.AutoIpv4GatewayHost {
		return fmt.Errorf("No host configured in config file")
	}

	return nil
}

//RecoveryDuration gets the recovery duration from a config structure
func (c *Config) RecoveryDuration() (time.Duration, error) {
	return time.ParseDuration(c.RecoveryTime)
}
