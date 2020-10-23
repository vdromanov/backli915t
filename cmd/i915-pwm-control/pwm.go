package main

import (
	"log"

	"github.com/vdromanov/i915-pwm-control/cmd/i915-pwm-control/regs"
)

var blcRegContents = regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)

func getFrequency() int {
	period, _ := regs.SplitPayload(blcRegContents)
	freq, err := regs.PeriodToFreq(period)
	if err != nil {
		log.Fatal(err)
	}
	return freq
}

func setFrequency(frequency int) {
	_, cycle := regs.SplitPayload(blcRegContents)
	log.Printf("Got cycle: 0x%08x\n", cycle)
	period, err := regs.FreqToPeriod(frequency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got period: 0x%08x\n", period)
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, regs.BuildPayload(period, cycle))
}

func changeFrequency(value int) {
	actualFreq := getFrequency()
	newFreq := actualFreq + value
	setFrequency(newFreq)
}
