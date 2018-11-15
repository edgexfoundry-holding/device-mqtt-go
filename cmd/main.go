// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-mqtt-go/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	version     string = "1.0"
	serviceName string = "edgex-device-mqtt"
)

func main() {
	sd := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, version, sd)
}
