package main

import (
	"log"

	"github.com/vdromanov/backli915t/internal/pkg/regs"
)

func getFrequency(blcRegContents *int) int {
	period, _ := regs.SplitPayload(*blcRegContents)
	freq, err := regs.PeriodToFreq(period)
	if err != nil {
		log.Fatal(err)
	}
	return freq
}

func setFrequency(frequency int, blcRegContents *int) {
	_, cycle := regs.SplitPayload(*blcRegContents)
	log.Printf("Got cycle: 0x%08x\n", cycle)
	period, err := regs.FreqToPeriod(frequency)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got period: 0x%08x\n", period)
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, regs.BuildPayload(period, cycle))
}

func changeFrequency(value int, blcRegContents *int) {
	actualFreq := getFrequency(blcRegContents)
	newFreq := actualFreq + value
	setFrequency(newFreq, blcRegContents)
}
