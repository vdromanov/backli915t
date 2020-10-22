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

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage:\n%s <frequency_to_set>\n", os.Args[0])
	}
	if checkModuleIsLoaded() {
		log.Printf("Driver %s has found in loaded ones\n", DriverName)
		blcRegContents, pchRegContents := regs.GetInfo() // Reading config regs once
		log.Printf("Actual pwm frequency is: %d\n", regs.ParsePayload(&blcRegContents, &pchRegContents))
		desiredFreq, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			log.Fatalln(err)
		}
		payload := regs.CalculatePayload(&blcRegContents, &pchRegContents, int(desiredFreq))
		log.Printf("Has calculated payload value: 0x%08x\n", payload)
		regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, payload)
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
