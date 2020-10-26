package regs

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// executable from package intel-gpu-tools
const execName = "intel_reg"

const bytelength = 4
const hexPrefix = "0x"
const hexMask = "%s%%0%dx"

// Configuration regs. Look at documentation of i915 driver.
const (
	BLC_PWM_PCH_CTL2_REG = 0xC8254
	PCH_RAWCLK_FREQ_REG  = 0xC6204
)

// ReadReg returns an int32 of reg's content
func ReadReg(reg int) int {
	_fmt := fmt.Sprintf(hexMask, hexPrefix, bytelength*2) // Making a fixed-length hex string
	regAddr := fmt.Sprintf(_fmt, reg)
	log.Println("Will read from: ", regAddr)
	out, err := exec.Command(execName, "read", regAddr).Output()
	if err != nil {
		log.Fatalln(err)
	}
	regValue := findHex(strings.Replace(string(out), regAddr, "", -1))[0] // A reg's addr is in output of intel-gpu-tools => replacing
	ret, err := strconv.ParseInt(strings.Replace(regValue, hexPrefix, "", -1), 16, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return int(ret)
}

// WriteReg writes an int32 to reg with given int number
func WriteReg(reg int, val int) {
	_fmt := fmt.Sprintf(hexMask, hexPrefix, bytelength*2) // Making a fixed-length hex string
	regAddr := fmt.Sprintf(_fmt, reg)
	regVal := fmt.Sprintf(_fmt, val)
	log.Printf("Will write %s to %s\n", regVal, regAddr)
	_, err := exec.Command(execName, "write", regAddr, regVal).Output()
	if err != nil {
		log.Fatalln(err)
	}
}

// findHex parses input string for hexadeciaml numbers and returns a slice of strings
func findHex(input string) []string {
	pattern := fmt.Sprintf("%s[[:xdigit:]]{%d}", hexPrefix, bytelength*2) // Hexadeciaml number of fixed length
	hexRe := regexp.MustCompile(pattern)
	return hexRe.FindAllString(input, -1)
}
