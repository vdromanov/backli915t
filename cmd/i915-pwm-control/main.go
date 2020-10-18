package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stdout, "Usage:\n%s <frequency_to_set>\n", os.Args[0])
		os.Exit(-1)
	}
	if checkModuleIsLoaded() {
		fmt.Fprintf(os.Stdout, "Driver %s has successfully loaded\n", DriverName) // TODO: set frequency here
	} else {
		fmt.Fprintf(os.Stderr, "Driver %s was not found in loaded. Exiting...\n", DriverName)
		os.Exit(-1)
	}
}

// All loaded modules are listed in procfs
func checkModuleIsLoaded() bool {
	content, err := ioutil.ReadFile("/proc/modules")
	if err != nil {
		panic(err)
	}
	return strings.Contains(string(content), DriverName)
}
