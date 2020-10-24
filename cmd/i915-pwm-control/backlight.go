package main

import (
	"log"

	"github.com/vdromanov/i915-pwm-control/cmd/i915-pwm-control/regs"
)

// var blcRegContents = regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)

func getBacklightPercent() int {
	period, cycle := regs.SplitPayload(blcRegContents)
	percent, err := regs.CycleToPercent(cycle, period)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Got backlight percent: %d\n", percent)
	return percent
}

func setBacklightPercent(percent int) {
	period, _ := regs.SplitPayload(blcRegContents)
	wantedCycle, err := regs.PercentToCycle(percent, period)
	if err != nil {
		log.Fatal(err)
	}
	payload := regs.BuildPayload(period, wantedCycle)
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, payload)
}

func changeBacklightPercent(value int) {
	actualBl := getBacklightPercent()
	newBl := actualBl + value
	setBacklightPercent(newBl)
}
