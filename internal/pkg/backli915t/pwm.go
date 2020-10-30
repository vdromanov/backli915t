package backli915t

import (
	"github.com/vdromanov/backli915t/internal/pkg/regs"
	log "github.com/vdromanov/backli915t/pkg/multilog"
)

func GetFrequency(blcRegContents *int) int {
	period, _ := regs.SplitPayload(*blcRegContents)
	freq, err := regs.PeriodToFreq(period)
	if err != nil {
		log.Info.Fatal(err)
	}
	return freq
}

func SetFrequency(frequency int, blcRegContents *int) {
	initialBacklightPercent := GetBacklightPercent(blcRegContents)
	_, cycle := regs.SplitPayload(*blcRegContents)
	log.Debug.Printf("Got cycle: 0x%08x\n", cycle)
	period, err := regs.FreqToPeriod(frequency)
	if err != nil {
		log.Info.Fatal(err)
	}
	log.Debug.Printf("Got period: 0x%08x\n", period)
	regs.WriteReg(regs.BLC_PWM_PCH_CTL2_REG, regs.BuildPayload(period, cycle))
	newBlcRegContents := regs.ReadReg(regs.BLC_PWM_PCH_CTL2_REG)
	SetBacklightPercent(initialBacklightPercent, &newBlcRegContents)
}

func ChangeFrequency(value int, blcRegContents *int) {
	actualFreq := GetFrequency(blcRegContents)
	newFreq := actualFreq + value
	SetFrequency(newFreq, blcRegContents)
}
