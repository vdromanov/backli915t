package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:\n%s <frequency_to_set>\n", os.Args[0])
	}
	if checkModuleIsLoaded() {
		log.Printf("Driver %s has found in loaded ones\n", DriverName)
		log.Printf("Actual pwm frequency is: %d\n", getFrequency())
		desiredFreq, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Will set to %d\n", desiredFreq)
		setFrequency(int(desiredFreq))
	} else {
		log.Fatalf("Driver %s was not found in loaded ones.\nExiting...\n", DriverName)
	}
}

// All loaded modules are listed in procfs
func checkModuleIsLoaded() bool {
	content, err := ioutil.ReadFile("/proc/modules")
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Contains(string(content), DriverName)
}
