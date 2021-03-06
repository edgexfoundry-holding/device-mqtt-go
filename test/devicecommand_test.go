//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"testing"

	"github.com/eclipse/paho.mqtt.golang"
)

// TestCommandReceive use to test receiving commands from the device service and responded back for get/set commands.
//
// Use a REST client to send a command to the service like:
// http://localhost:49982/api/v1/devices/{device id}>/message - use POST on this one with
// {"message":"some text"} in body http://localhost:49982/api/v1/devices/<device id>/ping - use GET
// http://localhost:49982/api/v1/devices/<device id>/randnum - use GET
//
// If command micro service is running, the same can be performed through command to device service
// like this http://localhost:48082/api/v1/device/<device id>/command/<command id>
//
// Requires the Device Service, Command, Core Data, Metadata and Mongo to all be running
func TestCommandReceive(t *testing.T) {
	var brokerUrl = "m12.cloudmqtt.com"
	var brokerPort = 17217
	var username = "tobeprovided"
	var password = "tobeprovided"
	var mqttClientId = "CommandSubscriber"
	var qos = 0
	var topic = "CommandTopic"

	uri := &url.URL{
		Scheme: "tcp",
		Host:   fmt.Sprintf("%s:%d", brokerUrl, brokerPort),
		User:   url.UserPassword(username, password),
	}

	client, err := createMqttClient(mqttClientId, uri)
	defer client.Disconnect(5000)
	if err != nil {
		t.Fatal(err)
	}

	token := client.Subscribe(topic, byte(qos), onCommandReceivedFromBroker)
	if token.Wait() && token.Error() != nil {
		t.Fatal(token.Error())
	}

	select {}
}

func onCommandReceivedFromBroker(client mqtt.Client, message mqtt.Message) {
	var request map[string]interface{}

	json.Unmarshal(message.Payload(), &request)
	uuid, ok := request["uuid"]
	if ok {
		log.Println(fmt.Sprintf("Command response received: topic=%v uuid=%v msg=%v", message.Topic(), uuid, string(message.Payload())))

		if request["method"] == "set" {
			sendTestData(request)
		} else {
			switch request["cmd"] {
			case "ping":
				request["ping"] = "pong"
				sendTestData(request)
			case "randnum":
				request["randnum"] = rand.Float64()
				sendTestData(request)
			case "message":
				request["message"] = "test-message"
				sendTestData(request)
			}
		}
	} else {
		log.Println(fmt.Sprintf("Command response ignored. No UUID found in the message: topic=%v msg=%v", message.Topic(), string(message.Payload())))
	}
}

func sendTestData(response map[string]interface{}) {
	var brokerUrl = "m12.cloudmqtt.com"
	var brokerPort = 17217
	var username = "tobeprovided"
	var password = "tobeprovided"
	var mqttClientId = "ResponsePublisher"
	var qos = byte(0)
	var topic = "ResponseTopic"

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

	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	client.Publish(topic, qos, false, jsonData)

	fmt.Println(fmt.Sprintf("Send response: %v", string(jsonData)))
}

func createMqttClient(clientID string, uri *url.URL) (mqtt.Client, error) {
	fmt.Println(fmt.Sprintf("Create MQTT client and connection: uri=%v clientID=%v ", uri.String(), clientID))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	opts.SetClientID(clientID)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		fmt.Println(fmt.Sprintf("Connection lost : %v", e))
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			fmt.Println(fmt.Sprintf("Reconnection failed : %v", e))
		} else {
			fmt.Println(fmt.Sprintf("Reconnection sucessful : %v", e))
		}
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	return client, nil
}
