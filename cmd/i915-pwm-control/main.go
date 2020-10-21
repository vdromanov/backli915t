package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vdromanov/i915-pwm-control/cmd/i915-pwm-control/regs"
)

// DriverName is a name of compatible Linux kernel module
const DriverName = "i915"

const (
	BLC_PWM_PCH_CTL2_REG = 0xC8254
	PCH_RAWCLK_FREQ_REG  = 0xC6204
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:\n%s <frequency_to_set>\n", os.Args[0])
	}
	if checkModuleIsLoaded() {
		log.Printf("Driver %s has found in loaded ones\n", DriverName) // TODO: set frequency here
		desiredFreq, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			log.Fatalln(err)
		}
		payload := regs.CalculatePayload(regs.ReadReg(BLC_PWM_PCH_CTL2_REG), regs.ReadReg(PCH_RAWCLK_FREQ_REG), int(desiredFreq))
		log.Printf("Has calculated payload value: 0x%08x\n", payload)
		regs.WriteReg(BLC_PWM_PCH_CTL2_REG, payload)
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
