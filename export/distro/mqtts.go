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
	"crypto/x509"
	"io/ioutil"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/edgex-go/export"
	"go.uber.org/zap"
)

type mqttsSender struct {
	client MQTT.Client
	topic  string
}

// NewMqttSender - create new mqtt sender
func NewMqttsSender(addr export.Addressable) Sender {
	tlsConfig := new(tls.Config)

	tlsConfig.ClientCAs = x509.NewCertPool()
	tlsConfig.RootCAs = x509.NewCertPool()

	if cert, err := tls.LoadX509KeyPair("/home/sergiom/mosquitto_certs/mosq-serv.crt", "/home/sergiom/mosquitto_certs/mosq-serv.key"); err == nil {
		tlsConfig.Certificates = []tls.Certificate{cert}

		if pemData, err := ioutil.ReadFile("/home/sergiom/mosquitto_certs/mosq-serv.crt"); err == nil {
			if !tlsConfig.ClientCAs.AppendCertsFromPEM(pemData) {
				logger.Fatal("Failed")
			}
		}

		if pemData, err := ioutil.ReadFile("/home/sergiom/mosquitto_certs/mosq-ca.crt"); err == nil {
			if !tlsConfig.RootCAs.AppendCertsFromPEM(pemData) {
				logger.Fatal("Failed appending certs")
			}
		} else {
			logger.Fatal("Failed reading CA")
		}
	} else {
		logger.Fatal("Failed loading x509 data")
		return nil
	}

	opts := MQTT.NewClientOptions()
	broker := addr.Protocol + "://" + addr.Address + ":" + strconv.Itoa(addr.Port)
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
