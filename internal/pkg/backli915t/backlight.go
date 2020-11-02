package backli915t

import (
	"github.com/vdromanov/backli915t/internal/pkg/regs"
	log "github.com/vdromanov/backli915t/pkg/multilog"
)

func GetBacklightPercent() int {
	blcRegContents := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)
	period, cycle := regs.SplitPayload(blcRegContents)
	percent, err := regs.CycleToPercent(cycle, period)
	if err != nil {
		log.Info.Fatal(err) // TODO: pull to main
	}
	return percent
}

func SetBacklightPercent(percent int) {
	blcRegContents := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)
	period, _ := regs.SplitPayload(blcRegContents)
	wantedCycle, err := regs.PercentToCycle(percent, period)
	if err != nil {
		log.Info.Fatal(err) // TODO: pull to main
	}
	payload := regs.BuildPayload(period, wantedCycle)
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, payload)
}

func ChangeBacklightPercent(value int) {
	actualBl := GetBacklightPercent()
	newBl := actualBl + value
	SetBacklightPercent(newBl)
}
