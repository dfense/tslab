package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dfense/tslab"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	errProcessingCLI = 1
	errSettingLogLvl = "error setting log level %s:"
	errCreatingFile  = "error creating file %s"
)

var (
	app = kingpin.New("TESLA Code Challenge", "A short async demonstration for concurrency")

	loglevel = kingpin.Flag("loglevel", "Set log level {PANIC, FATAL, ERROR, WARN, INFO, DEBUG}").Short('l').Default("INFO").String()

	// frequency range of agent velocity
	// swgNic   = kingpin.Flag("interface", "network interface to blast").Short('i').Required().String()
	swgPulse = kingpin.Flag("pulse", "time between mcast blast").Short('p').Default("10s").Duration()
	// number of default things to fire off
	swgPort = kingpin.Flag("swgport", "TCP port to connect session").Short('s').Default("3333").Int()

	// TODO build data at compile time
	// version   string
	// builddate string
	// githash   string
)

// main simplest CLI client. Expandable to many options
func main() {
	kingpin.Parse()

	// set log level, default sys.out
	// level, err := log.ParseLevel(*loglevel)
	// if err != nil {
	// 	log.Fatalf(errSettingLogLvl, err)
	// 	os.Exit(errProcessingCLI)
	// }
	fmt.Printf("Setting log level: %s", *loglevel)
	//log.SetLevel(level)
	log.SetLevel(log.DebugLevel)

	// catch Ctrl-C on terminal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		signalType := <-ch
		signal.Stop(ch)
		fmt.Println("") // clear the terminal ^C

		// stop all things and close out listener
		tslab.Stop(true)
		// this is a good place to flush everything to disk
		// before terminating.
		log.Println("Signal type : ", signalType)
		os.Exit(0)

	}()

	// probably should put inside supervisor
	eventWriter, err := tslab.NewEventWriter("events.txt")
	if err != nil {
		log.Fatalf(errCreatingFile, err)
	}
	listener := tslab.NewListener()
	listener.SetWriter(eventWriter)
	tslab.SetListener(listener)
	listener.StartListener()

	tslab.Console()
}
