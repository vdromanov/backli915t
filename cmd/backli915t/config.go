package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	BacklightPercent int `json:"bl"`
	PwmFrequency     int `json:"pwm"`
}

func LoadConfig(fpath string) Config {
	config := Config{ // Default values
		BacklightPercent: 30,
		PwmFrequency:     2000,
	}

	fp, err := os.Open(fpath)
	defer fp.Close()
	if err == nil {
		parser := json.NewDecoder(fp)
		err := parser.Decode(&config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println(err.Error())
	}
	return config
}

func DumpConfig(config *Config, fpath string) {
	fp, err := os.Create(fpath)
	defer fp.Close()
	if err == nil {
		writer := json.NewEncoder(fp)
		err := writer.Encode(config)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln(err)
	}
}
