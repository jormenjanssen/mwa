package main

import (
	"flag"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func main() {

	test := flag.Bool("test", false, "Check the configured action")
	debug := flag.Bool("debug", false, "Enable debugging")

	flag.Parse()

	// Try to open our config file from disk
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalf("Failed to open config file %v", err)
	}

	// Try and validate config from reader
	cfg, err := JsonConfigFromReader(file)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Configure logging
	configureLogging(cfg.LogPath, *debug, !cfg.IsDiskLoggingEnabled())

	// Preset out host to host
	host := TryGetHost(cfg.Host, cfg.Ipv4GatewayDetectionInterface)
	scr := NewScriptAction(cfg.Script)
	nhs := NewNetworkHealthService(host, PingHealthCheck{}, cfg.RecoveryDuration(), scr)

	// Console helpers
	if *test {
		if nhs.RecoveryAction != nil {
			nhs.RecoveryAction()
		}
		os.Exit(0)
	}

	// Log to our users
	log.Debugf("Logger debug enabled: %v", *debug)
	log.Printf("Running network health watchdog against target: %v with recovery time: [%v]", nhs.Address, nhs.RecoveryTime)

	// Run Application
	appCtx := NewApplicationContext()
	Watchdog(appCtx, nhs)
}

func TryGetHost(host string, autoDetectInterface string) string {

	targetHost := host

	if host != "" {

		timeout := 3 * time.Minute
		log.Infof("Target host selection is set to autodetect (using default ipv4 gateway) for adapter: %v with a maximum detection timeout of: %v", host, timeout)

		ip, err := GetIpv4TargetForAdapterGatewayWithTimeout(autoDetectInterface, timeout)

		if err != nil {
			log.Fatalf("Could not autodetect gateway error: %v took: %v", err, timeout)
		}

		// Assign our target host we discovered
		targetHost = ip.String()
	}

	return targetHost
}

func configureLogging(logpath string, debug bool, consoleOnly bool) {

	// Enable/Disable debug
	if debug {
		log.SetLevel(log.DebugLevel)

	} else {
		log.SetLevel(log.InfoLevel)
	}

	// For console only situations
	if consoleOnly {
		return
	}

	// Configure file based logger and log rotating functionality
	writer, err := rotatelogs.New(
		logpath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	// Make sure we always output to the console as fallback
	if err != nil {
		log.Errorf("Cannot set file logging error: %v (falling back to console only)", err)
	} else {
		log.Infof("Switching to file-based logging logfile could be found at: %v", path.Dir(logpath))
		log.SetOutput(writer)
	}
}
