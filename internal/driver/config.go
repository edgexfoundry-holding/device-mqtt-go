//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package driver

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type configuration struct {
	Incoming SubscribeInfo
	Response SubscribeInfo
}
type SubscribeInfo struct {
	Protocol     string
	Host         string
	Port         int
	Username     string
	Password     string
	Qos          int
	KeepAlive    int
	MqttClientId string
	Topic        string
}

// LoadConfigFromFile use to load toml configuration
func LoadConfigFromFile() (*configuration, error) {
	config := new(configuration)
	filePath := "./res/configuration-driver.toml"

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("could not load configuration file (%s): %v", filePath, err.Error())
	}

	err = toml.Unmarshal(file, config)
	if err != nil {
		return config, fmt.Errorf("unable to parse configuration file (%s): %v", filePath, err.Error())
	}
	return config, err
}
