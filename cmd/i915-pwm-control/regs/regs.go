package regs

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// executable from package intel-gpu-tools
const execName = "intel_reg"

// Bytelength is a length of gpu regs in bytes
const Bytelength = 4

// ReadReg returns a byteslice of reg's content
func ReadReg(reg int) []byte {
	_hexPrefix := "0x"
	_fmt := fmt.Sprintf("%s%%0%dx", _hexPrefix, Bytelength*2) // Making a fixed-length hex string
	regAddr := fmt.Sprintf(_fmt, reg)
	out, err := exec.Command(execName, "read", regAddr).Output()
	if err != nil {
		panic(err)
	}
	regValue := findHex(string(out))[1] // TODO: remove reg addr from output slice
	bytes, err := hex.DecodeString(strings.ReplaceAll(regValue, _hexPrefix, ""))
	if err != nil {
		panic(err)
	}
	return bytes
}

// findHex parses input string for hexadeciaml numbers
func findHex(input string) []string {
	pattern := fmt.Sprintf("0x[[:xdigit:]]{%d}", Bytelength*2) // Hexadeciaml number of fixed length
	hexRe := regexp.MustCompile(pattern)
	return hexRe.FindAllString(input, -1)
}
