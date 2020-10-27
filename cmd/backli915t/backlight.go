package main

import (
	"log"

	"github.com/vdromanov/backli915t/internal/pkg/regs"
)

func getBacklightPercent(blcRegContents *int) int {
	period, cycle := regs.SplitPayload(*blcRegContents)
	percent, err := regs.CycleToPercent(cycle, period)
	if err != nil {
		log.Fatal(err)
	}
	return percent
}

func setBacklightPercent(percent int, blcRegContents *int) {
	period, _ := regs.SplitPayload(*blcRegContents)
	wantedCycle, err := regs.PercentToCycle(percent, period)
	if err != nil {
		log.Fatal(err)
	}
	payload := regs.BuildPayload(period, wantedCycle)
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, payload)
}

func changeBacklightPercent(value int, blcRegContents *int) {
	actualBl := getBacklightPercent(blcRegContents)
	newBl := actualBl + value
	setBacklightPercent(newBl, blcRegContents)
}
