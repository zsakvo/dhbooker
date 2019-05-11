package main

import (
	"github.com/Unknwon/goconfig"
)

func setConfig(section string, key string, value string) bool {
	return config.SetValue(section, key, value)
}

func writeConfig() bool {
	err := goconfig.SaveConfigFile(config, "conf.ini")
	if err != nil {
		return false
	}
	return true
}
