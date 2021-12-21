package sensor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ryszard/sds011/go/sds011"
)

type Config struct {
	Topic          string
	SensorPortPath string
	CycleMinutes   uint8
	MqttBroker     string
}

func Start(c Config) {
	opts := mqtt.NewClientOptions().AddBroker(c.MqttBroker)
	opts.AutoReconnect = true
	opts.SetKeepAlive(30 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sensor, err := sds011.New(c.SensorPortPath)

	if err != nil {
		log.Fatalf("ERROR: sds011.New, %v", err)
	}

	defer sensor.Close()

	err = sensor.Awake()
	if err != nil {
		log.Fatalf("Error: sensor.Awake: %v", err)
	}

	if err = sensor.SetCycle(c.CycleMinutes); err != nil {
		log.Printf("ERROR: sensor.SetCycle: %v", err)
	}

	err = sensor.MakeActive()

	if err != nil {
		log.Fatalf("ERROR: sensor.MakeActive: %v", err)
	}

	for {
		point, err := sensor.Get()

		// var noError error
		// point, err := sds011.Point{
		// 	PM10:      float64(rand.Intn(20)),
		// 	PM25:      float64(rand.Intn(20)),
		// 	Timestamp: time.Now(),
		// }, noError
		// time.Sleep(5 * time.Second)

		if err != nil {
			log.Printf("ERROR: sensor.Get: %v", err)
			continue
		}

		point25, err := fmt.Fprintf(os.Stdout, "%v\n", point.PM25)

		pointJSON25, err := json.Marshal(point25)

		if err != nil {
			log.Printf("ERROR: Marshal: %v", err)
			continue
		}

		if token := client.Publish(c.Topic, 0, false, pointJSON25); token.Wait() && token.Error() != nil {
			fmt.Print(token.Error())
		}
	}
}
