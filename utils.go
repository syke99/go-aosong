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
	"bytes"
	"encoding/binary"

	i2c "github.com/syke99/go-i2c"
)

// Utility functions

// getS16BE extract 2-byte integer as signed big-endian.
func getS16BE(buf []byte) int16 {
	v := int16(buf[0])<<8 + int16(buf[1])
	return v
}

// getS16LE extract 2-byte integer as signed little-endian.
func getS16LE(buf []byte) int16 {
	w := getS16BE(buf)
	// exchange bytes
	v := (w&0xFF)<<8 + w>>8
	return v
}

// getU16BE extract 2-byte integer as unsigned big-endian.
func getU16BE(buf []byte) uint16 {
	v := uint16(buf[0])<<8 + uint16(buf[1])
	return v
}

// getU16LE extract 2-byte integer as unsigned little-endian.
func getU16LE(buf []byte) uint16 {
	w := getU16BE(buf)
	// exchange bytes
	v := (w&0xFF)<<8 + w>>8
	return v
}

// Calc CRC according to AM2320 specification.
func calcCRC_AM2320(buf []byte) uint16 {
	var seed uint16 = 0xFFFF
	for i := 0; i < len(buf); i++ {
		seed ^= uint16(buf[i])
		for j := 0; j < 8; j++ {
			if seed&0x01 != 0 {
				seed >>= 1
				seed ^= 0xA001
			} else {
				seed >>= 1
			}
		}
	}
	return seed
}

func calcCRC1(seed byte, buf []byte) byte {
	for i := 0; i < len(buf); i++ {
		b := buf[ /*len(buf)-1-*/ i]
		for j := 0; j < 8; j++ {
			if (seed^b)&0x01 != 0 {
				seed ^= 0x18
				seed >>= 1
				seed |= 0x80
				// crc = crc ^ 0x8c
			} else {
				seed >>= 1
			}
			b >>= 1
		}
	}
	return seed
}

// Read byte block from i2c device to struct object.
func readDataToStruct(i2c *i2c.I2C, byteCount int,
	byteOrder binary.ByteOrder, obj interface{}) error {
	buf1 := make([]byte, byteCount)
	_, err := i2c.ReadBytes(buf1)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(buf1)
	err = binary.Read(buf, byteOrder, obj)
	if err != nil {
		return err
	}
	return nil
}
