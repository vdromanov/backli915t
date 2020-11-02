package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	b "github.com/vdromanov/backli915t/internal/pkg/backli915t"
	log "github.com/vdromanov/backli915t/pkg/multilog"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

func main() {
	log.Info.AddOutput(os.Stderr)
	if !checkModuleIsLoaded() {
		log.Info.Fatalf("Driver %s was not found in loaded ones.\nNothing to do...\n", DriverName)
	}

	// Optional args for all modes
	generalArgs := flag.NewFlagSet("", flag.ExitOnError)
	debugPointer := generalArgs.Bool("debug", false, "show debug info in stdout")
	configPointer := generalArgs.String("c", "/usr/share/backli915t/config.json", "config fname")

	generalArgs.Usage = func() {
		fmt.Fprintf(os.Stderr, "General args (optional) are:\n")
		generalArgs.VisitAll(func(f *flag.Flag) { fmt.Fprintf(os.Stderr, "\t-%s - %v\n", f.Name, f.Usage) })
	}

	// Actions for both pwm and bl modes. Only one arg is allowed at the moment
	perModeArgs := flag.NewFlagSet("", flag.ExitOnError)
	incPointer := perModeArgs.Int("inc", 0xFFFFFFFF, "increment value")
	decPointer := perModeArgs.Int("dec", 0xFFFFFFFF, "decrement value")
	setPointer := perModeArgs.Int("set", 0xFFFFFFFF, "explicitly set value")

	perModeArgs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Action could be the one from:\n")
		perModeArgs.VisitAll(func(f *flag.Flag) { fmt.Fprintf(os.Stderr, "\t-%s - %v\n", f.Name, f.Usage) })
	}

	overallUsage := func(naming string) {
		fmt.Fprintf(os.Stderr, "\n\nUsage:\n")
		fmt.Fprintf(os.Stderr, "%s - utility to control pwm frequency and backlight level of displays with pwm-backlight via %s driver\n\n", naming, DriverName)
		fmt.Fprintf(os.Stderr, "Syntax is:\n\t%s [general args] <mode> <action>\n\n", os.Args[0])
		generalArgs.Usage()
		fmt.Fprintf(os.Stderr, "Modes are:\n\tpwm - adjusting frequency in Hz\n\tbl - adjusting brightness in %%\n")
		perModeArgs.Usage()
		fmt.Fprintf(os.Stderr, "Example:\n\t%s --debug bl -inc 10\n\t\n\nLaunching without args will apply last pwm and backlight\n", os.Args[0])
	}

	generalArgs.Parse(os.Args[1:])

	if *debugPointer {
		log.Debug.AddOutput(os.Stdout)
	}

	config, err := LoadConfig(*configPointer)
	if err != nil {
		log.Debug.Println(err)
	}

	if len(generalArgs.Args()) == 0 { // Applying values from a config, if no keys provided
		log.Info.Printf("Applying config: %v\n\n", config)
		b.SetBacklightPercent(config.BacklightPercent)
		b.SetFrequency(config.PwmFrequency)
		overallUsage(os.Args[0])
		return
	}

	dump := func() {
		err := DumpConfig(&config, *configPointer)
		if err != nil {
			log.Info.Println(err)
		}
	}

	defer dump()

	mode := generalArgs.Args()[0]

	perModeArgs.Parse(generalArgs.Args()[1:])

	setFlags := []string{}
	var newVal int
	var changeVal int

	perModeArgs.Visit(func(f *flag.Flag) { setFlags = append(setFlags, f.Name) }) // Iterating over all explicitly set args
	if len(setFlags) != 1 {                                                       // Only one arg is allowed per mode
		fmt.Fprintf(os.Stderr, "Got incorrect actions: %v\n\n", setFlags)
		overallUsage(os.Args[0])
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
			actualFreq = b.GetFrequency() + changeVal
		} else {
			actualFreq = newVal
		}
		b.SetFrequency(actualFreq) // TODO: pull out error
		config.PwmFrequency = actualFreq

	case "bl":
		log.Debug.Println("Operating with backlight")
		var actualBl int
		if newVal == 0xFFFFFFFF {
			actualBl = b.GetBacklightPercent() + changeVal
		} else {
			actualBl = newVal
		}
		b.SetBacklightPercent(actualBl) // TODO: pull out error
		config.BacklightPercent = actualBl
	default:
		fmt.Fprintf(os.Stderr, "%s - incorrect mode!\n\n", mode)
		overallUsage(os.Args[0])
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
