package main

import (
	"log"

	"github.com/vdromanov/i915-pwm-control/cmd/i915-pwm-control/regs"
)

var blcRegContents, pchRegContents = regs.GetInfo() // Reading config regs once

func getFrequency() int {
	freq, _ := regs.ParsePayload(&blcRegContents, &pchRegContents)
	return freq
}

func setFrequency(frequency int) {
	payload, err := regs.CalculatePayload(&blcRegContents, &pchRegContents, frequency)
	if err != nil {
		log.Fatal(err)
	}
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, payload)
}

func changeFrequency(value int) {
	actualFreq := getFrequency()
	newFreq := actualFreq + value
	setFrequency(newFreq)
}
