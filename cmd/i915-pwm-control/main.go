package main

import (
	"fmt"
	"io/ioutil"
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
		fmt.Fprintf(os.Stdout, "Usage:\n%s <frequency_to_set>\n", os.Args[0])
		os.Exit(-1)
	}
	if checkModuleIsLoaded() {
		fmt.Fprintf(os.Stdout, "Driver %s has successfully loaded\n", DriverName) // TODO: set frequency here
		fmt.Fprintf(os.Stdout, "Reg contents: 0x%08x\n", regs.ReadReg(BLC_PWM_PCH_CTL2_REG))
		freq, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			panic(err)
		}
		newVal := regs.CalculatePayload(BLC_PWM_PCH_CTL2_REG, PCH_RAWCLK_FREQ_REG, int(freq))
		fmt.Printf("Has calculated new value: 0x%08x\n", newVal)
		regs.WriteReg(BLC_PWM_PCH_CTL2_REG, newVal)
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
