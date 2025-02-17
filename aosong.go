//--------------------------------------------------------------------------------------------------
//
// Copyright (c) 2018 Denis Dyakov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
// BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
//--------------------------------------------------------------------------------------------------

package aosong

import (
	i2c "github.com/syke99/go-i2c"
)

// SensorType identify which Aosong Electronics
// humidity and temperature sensor is used.
// DHT12, AM2320 are supported.
type SensorType int

// Implement Stringer interface.
func (v SensorType) String() string {
	if v == DHT12 {
		return "DHT12"
	} else if v == AM2320 {
		return "AM2320"
	} else {
		return "!!! unknown !!!"
	}
}

const (
	// Aosong Electronics humidity and temperature sensor model DHT12.
	DHT12 SensorType = iota
	// Aosong Electronics humidity and temperature sensor model AM2320.
	AM2320
)

// Abstract Aosong Electronics sensor interface
// to control and gather data via I2C-bus.
type SensorInterface interface {
	ReadRelativeHumidityAndTemperatureMult10(i2c *i2c.I2C) (humidity int16, temperature int16, err error)
	//ReadTemperatureMult10(i2c *i2c.I2C) (int32, error)
	//ReadRelativeHumidityMult10(i2c *i2c.I2C) (int32, error)
}

type Sensor struct {
	sensorType SensorType
	sensor     SensorInterface
}

func NewSensor(sensorType SensorType) *Sensor {
	v := &Sensor{sensorType: sensorType}
	switch sensorType {
	case AM2320:
		v.sensor = &SensorAM2320{}
	case DHT12:
		v.sensor = &SensorDHT12{}
	}

	return v
}

func (v *Sensor) GetSensorType() SensorType {
	return v.sensorType
}

func (v *Sensor) ReadRelativeHumidityAndTemperature(i2c *i2c.I2C) (humidity float32,
	temperature float32, err error) {
	rh, temp, err := v.sensor.ReadRelativeHumidityAndTemperatureMult10(i2c)
	if err != nil {
		return 0, 0, err
	}
	rhf, tempf := float32(rh)/10, float32(temp)/10
	return rhf, tempf, nil
}
