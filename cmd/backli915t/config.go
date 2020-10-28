package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config stores PWM frequency and brightness level, defined by user
type Config struct {
	BacklightPercent int `json:"bl"`
	PwmFrequency     int `json:"pwm"`
}

// LoadConfig initiallizes a Config struct with values, stored in json (or a fallback-ones)
func LoadConfig(fpath string) (Config, error) {
	config := Config{ // A fallback values
		BacklightPercent: 30,
		PwmFrequency:     2000,
	}

	fp, err := os.Open(fpath)
	defer fp.Close()
	if err != nil {
		return config, fmt.Errorf("open %s", fpath, err)
	}
	parser := json.NewDecoder(fp)
	if err := parser.Decode(&config); err != nil {
		return config, fmt.Errorf("json parsing from %s", fpath, err)
	} else {
		return config, nil
	}
}

// DumpConfig writes actual config to a json
func DumpConfig(config *Config, fpath string) error {
	fp, err := os.Create(fpath)
	defer fp.Close()
	if err != nil {
		return fmt.Errorf("writing to %s", fpath, err)
	}
	writer := json.NewEncoder(fp)
	if err := writer.Encode(config); err != nil {
		return fmt.Errorf("encoding to json %v", *config, err)
	} else {
		return nil
	}
}
