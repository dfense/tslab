package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dfense/tslab"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	errProcessingCLI  = 1 // error return code from main
	errSettingLogLvl  = "error setting log level %s:"
	errCreatingFile   = "error creating file %s"
	errCreatingLogDir = "error creating log dir %s"

	logFile   = "teslacc.log" // file location for logged output from program code
	eventFile = "events.txt"  // the events database file basename. (uses rollover logging)
	logPath   = "log"         // directory created for both files above
)

var (

	// CLI args
	app       = kingpin.New("TESLA Code Challenge", "A short async demonstration for concurrency")
	autoStart = kingpin.Flag("autostart", "start (1) of each thing type {t, true, f, false}").Short('a').Default("true").String()
	loglevel  = kingpin.Flag("loglevel", "Set log level {PANIC, FATAL, ERROR, WARN, INFO, DEBUG}").Short('l').Default("INFO").String()

	// TODO build data at compile time
	// version   string
	// builddate string
	// githash   string
)

// main simplest CLI client.
func main() {
	kingpin.Parse()

	// set log level, default sys.out
	setupLogger(*loglevel)

	// allow for configuration and injection of items before starting supervisor
	// this can also be expanded into a richer Confuration object/service/factory
	configData := tslab.ConfigData{
		Autostart: *autoStart,
	}

	// create the io.WriterCloser and inject into listener
	eventWriter, err := tslab.NewEventWriter(logPath + string(os.PathSeparator) + eventFile)
	if err != nil {
		log.Fatalf(errCreatingFile, err)
	}

	// create a listener to inject
	listener := tslab.NewListener()
	listener.SetWriter(eventWriter)

	//inject Listener into supervisor
	tslab.SetListener(listener)
	listener.StartListener()

	// initialize Supervisor
	err = tslab.Initialize(configData)
	if err != nil {
		fmt.Printf("Error on Initialize %s\n", err)
		os.Exit(2)
	}

	// --- catch Ctrl-C on terminal ----
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-ch
		signal.Stop(ch)
		fmt.Println("") // clear the terminal ^C

		// stop all things and close out listener
		tslab.Stop(true)
		os.Exit(0)

	}()

	// create interactive console
	tslab.Console()
}

// setupLogger ensure log dir is created/existing, and configure loglevel, and logfile
func setupLogger(logLevel string) {

	err := os.MkdirAll(logPath, 0744)
	if err != nil {
		log.Fatal(errCreatingLogDir, err)
	}

	level, err := log.ParseLevel(*loglevel)
	if err != nil {
		log.Fatalf(errSettingLogLvl, err)
		os.Exit(errProcessingCLI)
	}
	log.SetLevel(level)

	// file writer
	log.SetOutput(&lumberjack.Logger{
		Filename:   logPath + string(os.PathSeparator) + logFile,
		MaxSize:    2,
		MaxAge:     2,
		MaxBackups: 5,
		Compress:   true,
	})
}
