//
// Copyright (c) 2017
// Cavium
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package distro

import (
	"crypto/tls"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/edgex-go/export"
	"go.uber.org/zap"
)

type mqttsSender struct {
	client MQTT.Client
	topic  string
}

const (
	MQTT_CERT = "certs/dummy.crt"
	MQTT_KEY  = "certs/dummy.key"
)

// NewMqttsSender - create new mqtts sender
func NewMqttsSender(addr export.Addressable) Sender {

	cert, err := tls.LoadX509KeyPair(MQTT_CERT, MQTT_KEY)

	if err != nil {
		logger.Fatal("Failed loading x509 data")
		return nil
	}

	tlsConfig := &tls.Config{
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	opts := MQTT.NewClientOptions()
	broker := addr.Protocol + "://" + addr.Address + ":" + strconv.Itoa(addr.Port) + "/" + addr.Path
	opts.AddBroker(broker)
	opts.SetClientID(addr.Publisher)
	opts.SetUsername(addr.User)
	opts.SetPassword(addr.Password)
	opts.SetAutoReconnect(false)
	opts.SetTLSConfig(tlsConfig)

	sender := &mqttsSender{
		client: MQTT.NewClient(opts),
		topic:  addr.Topic,
	}

	return sender
}

func (sender *mqttsSender) Send(data []byte) {
	if !sender.client.IsConnected() {
		logger.Info("Connecting to mqtts server")
		if token := sender.client.Connect(); token.Wait() && token.Error() != nil {
			logger.Warn("Could not connect to mqtts server, drop event")
			return
		}
	}

	token := sender.client.Publish(sender.topic, 0, false, data)
	// FIXME: could be removed? set of tokens?
	token.Wait()
	if token.Error() != nil {
		logger.Warn("mqtts error: ", zap.Error(token.Error()))
	} else {
		logger.Debug("Sent data: ", zap.ByteString("data", data))
	}
}
