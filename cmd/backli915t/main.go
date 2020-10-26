package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/vdromanov/backli915t/cmd/backli915t/regs"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

func main() {
	if !checkModuleIsLoaded() {
		log.Fatalf("Driver %s was not found in loaded ones.\nNothing to do...\n", DriverName)
	}

	// Flags for both pwm and bl modes. Only one flag is allowed at the moment
	perModeArgs := flag.NewFlagSet("", flag.ExitOnError)
	incPointer := perModeArgs.Int("inc", 0xFFFFFFFF, "increment value")
	decPointer := perModeArgs.Int("dec", 0xFFFFFFFF, "decrement value")
	setPointer := perModeArgs.Int("set", 0xFFFFFFFF, "explicitly set value")

	perModeArgs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s <mode> -<key> <value>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Allowed modes are:\n\tpwm - adjusting frequency in Hz\n\tbl - adjusting brightness in %%\n")
		fmt.Fprintf(os.Stderr, "Allowed keys are:\n")
		perModeArgs.VisitAll(func(f *flag.Flag) { fmt.Fprintf(os.Stderr, "\t-%s - %v\n", f.Name, f.Usage) })
		fmt.Fprintf(os.Stderr, "\nExample:\n\t%s pwm -inc 500\n", os.Args[0])
	}

	if len(os.Args) < 2 { // TODO: apply config in this case
		perModeArgs.Usage()
		os.Exit(-1)
	}
	mode := os.Args[1]

	perModeArgs.Parse(os.Args[2:])

	setFlags := make(map[string]bool) // Inc/Dec and Set modes are using different funcs => storing info, was the flag a Set mode or not
	var newVal int
	var actualFlag string

	perModeArgs.Visit(func(f *flag.Flag) { setFlags[f.Name] = (f.Name == "set"); actualFlag = f.Name }) // Iterating over all explicitly set flags
	if len(setFlags) != 1 {                                                                             // Only one flag is allowed per mode
		fmt.Fprintf(os.Stderr, "Got keys: %v\nOnly one key is allowed!\n\n", perModeArgs.Args())
		perModeArgs.Usage()
		os.Exit(-1)
	}

	switch actualFlag { // Could increment/decrement/explicitly set value
	case "inc":
		newVal = *incPointer
	case "dec":
		newVal = -*decPointer
	case "set":
		newVal = *setPointer
	}

	blcRegVal := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)

	switch mode { // Working mode select
	case "pwm":
		log.Println("Have chosen PWM mode")
		if setFlags[actualFlag] { // Explicitly setting a new value
			setFrequency(newVal, &blcRegVal)
		} else { // Changing existing value
			changeFrequency(newVal, &blcRegVal)
		}
	case "bl":
		log.Println("Have chosen Backlight mode")
		if setFlags[actualFlag] {
			setBacklightPercent(newVal, &blcRegVal)
		} else {
			changeBacklightPercent(newVal, &blcRegVal)
		}
	default:
		fmt.Fprintf(os.Stderr, "Incorrect mode %s has provided!\n\n", mode)
		perModeArgs.Usage()
		os.Exit(-1)
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
