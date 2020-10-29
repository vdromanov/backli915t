package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/vdromanov/backli915t/internal/pkg/regs"
	log "github.com/vdromanov/backli915t/pkg/multilog"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

const configFname = "/usr/share/backli915t/config.json"

func main() {
	log.Info.AddOutput(os.Stderr)
	if !checkModuleIsLoaded() {
		log.Info.Fatalf("Driver %s was not found in loaded ones.\nNothing to do...\n", DriverName)
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

	blcRegVal := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)
	config, err := LoadConfig(configFname)
	if err != nil {
		log.Debug.Println(err)
	}

	if len(os.Args) < 2 { // Applying values from a config, if no keys provided
		log.Info.Printf("Applying config: %v\n", config)
		setFrequency(config.PwmFrequency, &blcRegVal)
		newBlRegContents := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)
		setBacklightPercent(config.BacklightPercent, &newBlRegContents)
		return
	}

	dump := func() {
		err := DumpConfig(&config, configFname)
		if err != nil {
			log.Info.Println(err)
		}
	}

	defer dump()

	mode := os.Args[1]

	perModeArgs.Parse(os.Args[2:])

	setFlags := []string{}
	var newVal int
	var changeVal int

	perModeArgs.Visit(func(f *flag.Flag) { setFlags = append(setFlags, f.Name) }) // Iterating over all explicitly set flags
	if len(setFlags) != 1 {                                                       // Only one flag is allowed per mode
		fmt.Fprintf(os.Stderr, "Got keys: %v\nOnly one key is allowed!\n\n", setFlags)
		perModeArgs.Usage()
		os.Exit(-1)
	}

	switch setFlags[0] { // Could increment/decrement/explicitly set value
	case "inc":
		changeVal = *incPointer
		newVal = 0xFFFFFFFF
	case "dec":
		changeVal = -*decPointer
		newVal = 0xFFFFFFFF
	case "set":
		newVal = *setPointer
		changeVal = 0xFFFFFFFF
	}

	switch mode { // Working mode select
	case "pwm":
		log.Debug.Println("Operating with PWM")
		var actualFreq int
		if newVal == 0xFFFFFFFF {
			actualFreq = getFrequency(&blcRegVal) + changeVal
		} else {
			actualFreq = newVal
		}
		setFrequency(actualFreq, &blcRegVal) // TODO: pull out error
		config.PwmFrequency = actualFreq

	case "bl":
		log.Debug.Println("Operating with backlight")
		var actualBl int
		if newVal == 0xFFFFFFFF {
			actualBl = getBacklightPercent(&blcRegVal) + changeVal
		} else {
			actualBl = newVal
		}
		setBacklightPercent(actualBl, &blcRegVal) // TODO: pull out error
		config.BacklightPercent = actualBl
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
		log.Info.Fatalln(err)
	}
	return strings.Contains(string(content), DriverName)
}
