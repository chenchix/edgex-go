//
// Copyright (c) 2018
// Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	models "github.com/edgexfoundry/edgex-go/core/domain/models"
)

//"min":"-40",,"max":"140" {"cosa":"peo"}
const valueDescriptor string = `{"name":"temperature","type":"J","min":"-40","max":"140",
	"uomLabel":"degree cel","defaultValue":"0","formatting":"%s","labels":["temp","hvac"]}`

const event string = `{"origin":1471806386919,"device":"livingroomthermostat",
		"readings":[{"origin":1471806386919,"name":"temperature","value":"{\"cosa\":\"peo\"}"}, {"origin":1471806386919,"name":"humidity","value":"-12"}]}`

func isValidValueDescriptor(reading model.Readings, ev models.Event) bool {
	vd, err := dbc.ValueDescriptorByName(reading.name)

	// Not an error if not found
	if err == clients.ErrNotFound {
		continue
	}

	fmt.Println(vd)

	/*ev := models.Event{}

	if err := json.Unmarshal([]byte(event), &ev); err != nil {
		panic(err)
	}

	vd := models.ValueDescriptor{}

	if err := json.Unmarshal([]byte(valueDescriptor), &vd); err != nil {
		panic(err)
	}*/

	switch vd.Type {
	case "B": // boolean
		return validBoolean(ev)
	case "F": // floating point
		return validFloat(ev, vd)
	case "I": // integer
		return validInteger(ev, vd)
	case "S": // string or character data
		return validString(ev)
	case "J": // JSON data
		return validJSON(ev)
	default:
		return false
	}

	return false
}

func validBoolean(reading models.Event) bool {

	for i := range reading.Readings {
		data, err := strconv.ParseBool(reading.Readings[i].Value)
		fmt.Println(data)

		if err != nil {
			return false
		}
	}
	return true
}

func validFloat(reading models.Event, vd models.ValueDescriptor) bool {

	if vd.Min == nil {
		return true
	}

	min, err := strconv.ParseFloat(vd.Min.(string), 64)
	if err != nil {
		return true
	}

	if vd.Max == nil {
		return true
	}

	max, err := strconv.ParseFloat(vd.Max.(string), 64)
	if err != nil {
		return true
	}

	for i := range reading.Readings {

		value, err := strconv.ParseFloat(reading.Readings[i].Value, 64)
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}

		fmt.Println(reflect.TypeOf(value))
		fmt.Println(min, value, max)

		if value > max || value < min {
			return false
		}

	}
	return true
}

func validInteger(reading models.Event, vd models.ValueDescriptor) bool {
	if vd.Min == nil {
		return true
	}

	min, err := strconv.ParseInt(vd.Min.(string), 10, 64)
	if err != nil {
		return true
	}

	if vd.Max == nil {
		return true
	}

	max, err := strconv.ParseInt(vd.Max.(string), 10, 64)
	if err != nil {
		return true
	}

	for i := range reading.Readings {

		value, err := strconv.ParseInt(reading.Readings[i].Value, 10, 64)
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}

		fmt.Println(reflect.TypeOf(value))
		fmt.Println(min, value, max)

		if value > max || value < min {
			return false
		}

	}
	return true
}

func validString(reading models.Event) bool {
	for i := range reading.Readings {
		if reading.Readings[i].Value == "" {
			return false
		}
	}

	return true
}
func validJSON(reading models.Event) bool {
	var js interface{}
	for i := range reading.Readings {
		err := json.Unmarshal([]byte(reading.Readings[i].Value), &js)
		fmt.Println(reading.Readings[i].Value, err)
		if err != nil {
			return false
		}
	}
	return true
}

/*
func main() {

	t := isValidValueDescriptor()

	fmt.Println("Is valid: ", t)
}*/
