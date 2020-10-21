package regs

import "log"

// splitPayload makes 2 2-byte ints from a4-byte one
func splitPayload(val int) (int, int) {
	period := val >> 16   // A first 2 bytes
	cycle := val & 0xFFFF // Last 2 vytes
	return period, cycle
}

// buildPayload makes one 4-byte int from 2 2-byte ines
func buildPayload(period, cycle int) int {
	return ((period << 16) + cycle)
}

// CalculatePayload does the magic of embedding desired frequency into a correct payload for intel-gpu reg
func CalculatePayload(blcPwmRegVal, pchRawclkRegVal, freq int) int {
	_, cycle := splitPayload(blcPwmRegVal)
	log.Printf("Got cycle: 0x%08x\n", cycle)
	period := int(1000000 * pchRawclkRegVal / 128 / freq)
	log.Printf("Got period: 0x%08x\n", period)
	regVal := buildPayload(period, cycle)
	return regVal
}

// ParsePayload calculates PWM frequency in Hz from BLC and PCH regs values
func ParsePayload(blcPwmRegVal, pchRawclkRegVal int) int {
	period, _ := splitPayload(blcPwmRegVal)
	freq := int(1E6 * pchRawclkRegVal / 128 / period)
	return freq
}
