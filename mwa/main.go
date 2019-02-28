package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Flag parsing (-debug && -test)
	test := flag.Bool("test", false, "Check the configured action")
	debug := flag.Bool("debug", false, "Enable debugging")
	console := flag.Bool("console", false, "Force logging output to console only (stdout)")
	version := flag.Bool("version", false, "Show version info and exit")

	flag.Parse()

	// For version info only
	if *version {
		ShowInfo()
		os.Exit(0)
	}

	// Try to open our config file from disk
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalf("Failed to open config file at: %v error: %v", file, err)
	}

	// Validate config from reader or fail if we cannot parse this
	cfg, err := JsonConfigFromReader(file)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v error: %v", file, err)
	}

	// Configure logging
	configureLogging(cfg.LogPath, *debug, !cfg.IsDiskLoggingEnabled() || *console)

	// Get our host or fail
	host, err := GetTargetHost(cfg.Host, cfg.Ipv4GatewayDetectionInterface)
	if err != nil {
		log.Fatalf("Failed to get target host exitting error: %v", err)
	}

	scr := NewScriptAction(cfg.Script)
	nhs := NewNetworkHealthService(host, GetNeworkHealthCheck(cfg.Host), cfg.RecoveryDuration(), scr)

	// Console helpers
	if *test {
		if nhs.RecoveryAction != nil {
			nhs.RecoveryAction()
		}
		os.Exit(0)
	}

	log.Println("MWA Application starting")

	// Log os uptime
	uptime, err := GetUptime()
	if err == nil {
		log.Println(fmt.Sprintf("Operating System uptime: %s", uptime))
	} else {
		log.Println(fmt.Sprintf("Failed getting os uptime reason: %s", err))
	}

	// Log to our users we are in debug
	if *debug {
		log.Debug("Debug logging [enabled]")
	}

	// Run Application
	appCtx := NewApplicationContext()
	Watchdog(appCtx, nhs)
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
		log.Errorf("Cannot set file logging error: %v (falling back to console log only)", err)
	} else {

		logdir := path.Dir(logpath)
		log.Infof("Switching to file-based logging logfiles could be found at: %v", logdir)

		if _, err := os.Stat(logdir); os.IsNotExist(err) {
			log.Infof("Log directory does not exsist creating: %v", logdir)
			// 0664 (RW+RW+R)
			err := os.MkdirAll(logdir, 0664)
			if err != nil {
				log.Errorf("Cannot create logging directory: %v error: %v (falling back to console log only)", logdir, err)
				return
			}
		}

		log.SetOutput(writer)
	}
}
