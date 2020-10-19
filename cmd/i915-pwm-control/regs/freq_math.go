package regs

import "fmt"

func splitPayload(val int32) (int32, int32) {
	period := val >> 16
	cycle := val & 0xFFFF
	return period, cycle
}

func buildPayload(period, cycle int32) int32 {
	return ((period << 16) + cycle)
}

func CalculatePayload(blcPwmReg, pchRawclkReg, freq int) int32 {
	_, cycle := splitPayload(ReadReg(blcPwmReg))
	fmt.Printf("Got cycle: 0x%08x\n", cycle)
	ret := int32(1000000 * ReadReg(pchRawclkReg) / 128 / int32(freq))
	fmt.Printf("Got ret: 0x%08x\n", ret)
	// period, _ := splitPayload(ret)
	period := ret
	fmt.Printf("Got period: 0x%08x\n", period)
	regVal := buildPayload(period, cycle)
	fmt.Printf("Got new reg value: 0x%08x\n", regVal)
	return regVal
}
