package main

import "mattianatali.it/sds011-mqtt/internal/sensor"

func main() {
	c := sensor.Config{
		Topic:          "dust/PM25",
		SensorPortPath: "/dev/ttyUSB0",
		CycleMinutes:   5,
		MqttBroker:     "tcp://192.168.1.60:1883",
	}
	sensor.Start(c)
}
