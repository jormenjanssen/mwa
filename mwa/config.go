package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const MinimalRecoveryTime = 10 * time.Second

type Config struct {
	RecoveryTime                  string
	MonitorOnly                   bool
	Host                          string
	Ipv4GatewayDetectionInterface string
	Script                        string
	LogPath                       string
	Protocol                      string
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

	err := c.IsValidRecoveryDuration()
	if err != nil {
		return err
	}

	if c.RecoveryDuration() < MinimalRecoveryTime {
		return fmt.Errorf("Configured recovery time of %v is < %v treshold", c.RecoveryDuration(), MinimalRecoveryTime)
	}

	if !c.MonitorOnly && c.Script == "" {
		return fmt.Errorf("Script not configured and not in monitor only mode")
	}

	if c.Host == "" && c.Ipv4GatewayDetectionInterface == "" {
		return fmt.Errorf("No host or auto detection using default gateway interface configured in config file")
	}

	return nil
}

func (c *Config) IsValidRecoveryDuration() error {
	_, err := time.ParseDuration(c.RecoveryTime)
	return err
}

//RecoveryDuration gets the recovery duration from a config structure
func (c *Config) RecoveryDuration() time.Duration {
	t, _ := time.ParseDuration(c.RecoveryTime)
	return t
}

//IsDiskLoggingEnabled checks if disk logging is enabled by loooking at the configured logpath
func (c *Config) IsDiskLoggingEnabled() bool {
	return c.LogPath != ""
}
