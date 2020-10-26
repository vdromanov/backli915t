package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/vdromanov/i915-pwm-control/cmd/i915-pwm-control/regs"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

func main() {
	if !checkModuleIsLoaded() {
		log.Fatalf("Driver %s was not found in loaded ones.\nExiting...\n", DriverName)
	}

	incPointer := flag.Int("inc", 0xFFFFFFFF, "Increment value")
	decPointer := flag.Int("dec", 0xFFFFFFFF, "Decrement value")
	setPointer := flag.Int("set", 0xFFFFFFFF, "Set value")
	flag.Parse()

	setFlags := make(map[string]bool)
	var newVal int
	var actualFlag string

	flag.Visit(func(f *flag.Flag) { setFlags[f.Name] = (f.Name == "set"); actualFlag = f.Name })
	if len(setFlags) > 1 { // Only one flag is allowed
		log.Fatalln("Choose only one option")
	}

	switch actualFlag {
	case "inc":
		newVal = *incPointer
	case "dec":
		newVal = -*decPointer
	case "set":
		newVal = *setPointer
	}

	log.Println(newVal)
	log.Println(setFlags)
	log.Println(actualFlag)

	blcRegVal := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)

	switch flag.Arg(0) { // Working mode select
	case "pwm":
		log.Println("Have chosen PWM mode")
		if setFlags[actualFlag] { // Setting mode
			setFrequency(newVal, &blcRegVal)
		} else { // Changing mode
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
		log.Fatalln("Choose operating mode from bl/pwm")
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
