//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"testing"
	"time"
)

// TestDeviceSender use to to generate random numbers and send them into the device service as if a sensor
// was sending the data. Requires the Device Service along with Mongo, Core Data, and Metadata to be running
func TestDeviceSender(t *testing.T) {
	var brokerUrl = "m12.cloudmqtt.com"
	var brokerPort = 17217
	var username = "tobeprovided"
	var password = "tobeprovided"
	var mqttClientId = "IncomingDataPublisher"
	var qos = byte(0)
	var topic = "DataTopic"

	uri := &url.URL{
		Scheme: "tcp",
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createMqttClient(mqttClientId, uri)
	defer client.Disconnect(5000)
	if err != nil {
		fmt.Println(err)
	}

	var data = make(map[string]interface{})
	data["name"] = "MQTT test device"
	data["cmd"] = "randnum"
	data["method"] = "get"

	for {
		data["randnum"] = rand.Float64()
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		client.Publish(topic, qos, false, jsonData)

		fmt.Println(fmt.Sprintf("Send response: %v", string(jsonData)))

		time.Sleep(time.Second * time.Duration(30))
	}

}
