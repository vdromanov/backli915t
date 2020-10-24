package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vdromanov/i915-pwm-control/cmd/i915-pwm-control/regs"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

func main() {
	if !checkModuleIsLoaded() {
		log.Fatalf("Driver %s was not found in loaded ones.\nExiting...\n", DriverName)
	}

	if len(os.Args) < 2 {
		log.Fatalln("Choose operating mode!")
	}
	var pwmChange string
	var blChange string
	var blcRegVal int

	flag.StringVar(&pwmChange, "pwm", "+0", "Specify PWM frequency in Hz") // A value is +100 like
	flag.StringVar(&blChange, "bl", "+0", "Specify Backlight level in %")  // A value is +100 like
	flag.Parse()

	blcRegVal = regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)

	switch string(pwmChange[0]) {
	case "+", "-":
		if change, err := strconv.ParseInt(pwmChange, 10, 64); err == nil {
			changeFrequency(int(change), &blcRegVal)
		} else {
			log.Fatal(err)
		}
	case "=":
		if set, err := strconv.ParseInt(pwmChange[1:], 10, 64); err == nil {
			setFrequency(int(set), &blcRegVal)
		} else {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Value should be like +<freq>/-<freq>/=<freq>")
	}

	blcRegVal = regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)

	switch string(blChange[0]) {
	case "+", "-":
		if change, err := strconv.ParseInt(blChange, 10, 64); err == nil {
			changeBacklightPercent(int(change), &blcRegVal)
		} else {
			log.Fatal(err)
		}
	case "=":
		if set, err := strconv.ParseInt(blChange[1:], 10, 64); err == nil {
			setBacklightPercent(int(set), &blcRegVal)
		} else {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Value should be like +<freq>/-<freq>/=<freq>")
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
