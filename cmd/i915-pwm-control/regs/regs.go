package regs

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// executable from package intel-gpu-tools
const execName = "intel_reg"

// Bytelength is a length of gpu regs in bytes
const Bytelength = 4

const HexPrefix = "0x"
const hexMask = "%s%%0%dx"

// ReadReg returns a byteslice of reg's content
func ReadReg(reg int) int32 {
	_fmt := fmt.Sprintf(hexMask, HexPrefix, Bytelength*2) // Making a fixed-length hex string
	regAddr := fmt.Sprintf(_fmt, reg)
	fmt.Println("Will read from: ", regAddr)
	out, err := exec.Command(execName, "read", regAddr).Output()
	if err != nil {
		panic(err)
	}
	regValue := findHex(string(out))[1] // TODO: remove reg addr from output slice
	// bytes, err := hex.DecodeString(strings.ReplaceAll(regValue, hexPrefix, ""))
	ret, err := strconv.ParseInt(strings.Replace(regValue, HexPrefix, "", -1), 16, 32)
	if err != nil {
		panic(err)
	}
	return int32(ret)
}

// WriteReg writes a byteslice to reg with given int number
func WriteReg(reg int, val int32) {
	_fmt := fmt.Sprintf(hexMask, HexPrefix, Bytelength*2) // Making a fixed-length hex string
	regAddr := fmt.Sprintf(_fmt, reg)
	regVal := fmt.Sprintf(_fmt, val)
	// regVal := hex.EncodeToString(val)
	fmt.Printf("Will write %s to %s\n", regVal, regAddr)
	_, err := exec.Command(execName, "write", regAddr, regVal).Output()
	if err != nil {
		panic(err)
	}
}

// findHex parses input string for hexadeciaml numbers
func findHex(input string) []string {
	pattern := fmt.Sprintf("%s[[:xdigit:]]{%d}", HexPrefix, Bytelength*2) // Hexadeciaml number of fixed length
	hexRe := regexp.MustCompile(pattern)
	return hexRe.FindAllString(input, -1)
}
