package regs

import (
	"errors"
)

// SplitPayload makes 2 2-byte ints from a4-byte one
func SplitPayload(val int) (int, int) {
	period := val >> 16   // A first 2 bytes
	cycle := val & 0xFFFF // Last 2 vytes
	return period, cycle
}

// BuildPayload makes one 4-byte int from 2 2-byte ines
func BuildPayload(period, cycle int) int {
	return ((period << 16) + cycle)
}

func FreqToPeriod(freq int) (int, error) {
	period := int(1E6 * ReadReg(PCH_RAWCLK_FREQ_REG) / 128 / freq)
	if (period < 0xFFFF) && (period > 0) {
		return period, nil
	} else {
		return -1, errors.New("Got invalid period value")
	}
}

func PeriodToFreq(period int) (int, error) {
	if (period < 0xFFFF) && (period > 0) {
		return int(1E6 * ReadReg(PCH_RAWCLK_FREQ_REG) / 128 / period), nil
	} else {
		return -1, errors.New("Got invalid period value")
	}
}

func CycleToPercent(cycle, period int) (int, error) {
	percent := int(float32(cycle) / float32(period) * 100.0)
	if (percent < 0) || (percent > 100) {
		return -1, errors.New("Incorrect percent value")
	}
	return percent, nil
}

func PercentToCycle(percent, period int) (int, error) {
	if (percent < 0) || (percent > 100) {
		return -1, errors.New("Incorrect percent value")
	}
	return int(period * percent / 100), nil
}
