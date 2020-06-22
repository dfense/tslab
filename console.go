package tslab

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dfense/tslab/things"
	log "github.com/sirupsen/logrus"
)

const (
	maxQuantity = 100
)

var (
	errNoCommandEntered   = errors.New("incorrect command, try again")
	errInvalidThingType   = errors.New("invalid thing type, try again")
	errImproperNumberArgs = errors.New("wrong number of arguents to command, try again")
	errConvertingToInt    = errors.New("error converting param to int, try again")
	errMaxQtyExceeded     = errors.New("error maximum qty of things to create is (100)")

	lastType = things.TBatteryPack
)

// Console main loop for console text menu
func Console() {

	welcomeScreen()

	//read commands
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("Command (h for help): ")
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error: %s\n", err)
		}

		// convert CRLF to LF
		command = strings.Replace(command, "\n", "", -1)
		err = processCommand(command)
		if err != nil {
			fmt.Println(err) // output to console
		}

	}
}

// welcomeScreen print welcome text to console
func welcomeScreen() {
	logo := `
===========================================================
||       TESLA Code Challenge                            ||
||            v.1.0.0                                    ||
===========================================================

`
	fmt.Println(logo)
}

// printMenu prints menu options and back to command prompt
func printMenu() {

	menu := `

              TESLA Code Challeng Menu
Command | Arguments       | Description 
----------------------------------------------------------------
   h    |                 | print this help menu
   li   |                 | list all things running/publishing
   nt   | <type> <qty>    | new thing by <type> <qty 1-1000>
   stop |                 | stop & delete all things, exit program 
   st   | <type>          | stop things by type [see types below]
   si   | <id>            | stop thing by id number
   sa   |                 | stop all things, do NOT exit program
   q    |                 | quit, stop all things, exit program
-----------------------------------------------------------------
valid thing <type> -> [b=battery, i=inverter, l=light]

`
	fmt.Println(menu)
}

// processCommand verify and dispatch command from menu
func processCommand(command string) error {

	command = strings.ToLower(command)
	c := strings.Fields(command) // simple space delimited parser
	if len(c) == 0 {
		return errNoCommandEntered
	}
	// simple simple parser. If it gets more complex,  reconsider a lib
	switch c[0] {
	case "h":
		printMenu()
	case "li":
		cids := GetThingsList()
		fmt.Println("\n                      list of things                              ")
		fmt.Println(" CID     | ThingType        | CreatedOn                 | TTLEvts    ")
		fmt.Println("-----------------------------------------------------------------------")
		for _, cid := range cids {
			fmt.Printf(" %-7d| %-18s| %-26s| %-10d\n", cid.CidNumber, cid.Type, cid.CreateTime.Format(time.RFC3339), cid.TTLEvents)
		}
		fmt.Printf("(%d total thing(s) running) \n\n", len(cids))

	// new thing, create
	case "nt":

		// use defaults if no arguments
		switch len(c) {
		case 1:
			fmt.Printf("\nCreated Default 1 %s\n\n", lastType)
			CreateThing(lastType, 1)

			// default used, then rotate to next thing in line, variety :-)
			if lastType == things.TLight {
				lastType = things.TBatteryPack
				break
			}
			lastType++

		// use parameters if provided
		case 3:
			qty, err := strconv.Atoi(c[2])
			if err != nil {
				return errConvertingToInt
			}
			if qty > maxQuantity {
				return errMaxQtyExceeded
			}

			thingType, err1 := verifyThingType(c[1])
			if err1 != nil {
				return errInvalidThingType
			}
			CreateThing(thingType, qty)

		default:
			return errImproperNumberArgs
		}

		fmt.Println("\n--- success: new thing(s) created ---")
		fmt.Println("")

	case "sa":
		Stop(false)

	case "st":
		switch len(c) {
		case 2:
			thingType, err := verifyThingType(c[1])
			if err != nil {
				return err
			}
			StopThingsByType(thingType)
			// do it
		default:
			return errImproperNumberArgs
		}
	case "si":
		switch len(c) {
		case 2:
			id, err := strconv.ParseUint(c[1], 10, 64)
			if err != nil {
				return errConvertingToInt
			}
			err = StopThingsByCID(id)
			if err != nil {
				return err
			}
		default:
			return errImproperNumberArgs
		}
	case "q", "stop":
		Stop(true)
	default:
		fmt.Println("\nunrecognized command, try again!")
		fmt.Println("")
	}

	return nil
}

func verifyThingType(c string) (things.ThingType, error) {
	switch c {
	case "b":
		return things.TBatteryPack, nil
	case "i":
		return things.TInverter, nil
	case "l":
		return things.TLight, nil
	default:
		return 0, errInvalidThingType
	}
}
