// Package ipfix decodes IPFIX packets
//: ----------------------------------------------------------------------------
//: Copyright (C) 2017 Verizon.  All Rights Reserved.
//: All Rights Reserved
//:
//: file:    marshal.go
//: details: TODO
//: author:  Mehrdad Arshad Rad
//: date:    02/01/2017
//:
//: Licensed under the Apache License, Version 2.0 (the "License");
//: you may not use this file except in compliance with the License.
//: You may obtain a copy of the License at
//:
//:     http://www.apache.org/licenses/LICENSE-2.0
//:
//: Unless required by applicable law or agreed to in writing, software
//: distributed under the License is distributed on an "AS IS" BASIS,
//: WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//: See the License for the specific language governing permissions and
//: limitations under the License.
//: ----------------------------------------------------------------------------
package ipfix

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

var errUknownMarshalDataType = errors.New("unknown data type to marshal")

// JSONMarshal encodes IPFIX message
func (m *Message) JSONMarshal(b *bytes.Buffer) ([]byte, error) {
	b.WriteString("{")

	// encode agent id
	m.encodeAgent(b)

	// encode header
	m.encodeHeader(b)

	// encode data sets
	if err := m.encodeDataSet(b); err != nil {
		return nil, err
	}

	b.WriteString("}")

	return b.Bytes(), nil
}

func (m *Message) encodeDataSet(b *bytes.Buffer) error {
	var length, dsLength int

	b.WriteString("\"DataSets\":")
	dsLength = len(m.DataSets)

	b.WriteString("[")
	for i := range m.DataSets {
		length = len(m.DataSets[i])

		b.WriteString("[")
		for j := range m.DataSets[i] {
			b.WriteString("{\"ID\":")
			b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].ID), 10))
			b.WriteString(",\"Value\":")

			switch m.DataSets[i][j].Value.(type) {
			case uint:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint)), 10))
			case uint8:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint8)), 10))
			case uint16:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint16)), 10))
			case uint32:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint32)), 10))
			case uint64:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(uint64)), 10))
			case int:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int)), 10))
			case int8:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int8)), 10))
			case int16:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int16)), 10))
			case int32:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int32)), 10))
			case int64:
				b.WriteString(strconv.FormatInt(int64(m.DataSets[i][j].Value.(int64)), 10))
			case float32:
				b.WriteString(strconv.FormatFloat(float64(m.DataSets[i][j].Value.(float32)), 'E', -1, 32))
			case float64:
				b.WriteString(strconv.FormatFloat(m.DataSets[i][j].Value.(float64), 'E', -1, 64))
			case string:
				b.WriteString("\"")
				b.WriteString(m.DataSets[i][j].Value.(string))
				b.WriteString("\"")
			case net.IP:
				b.WriteString("\"")
				b.WriteString(m.DataSets[i][j].Value.(net.IP).String())
				b.WriteString("\"")
			case net.HardwareAddr:
				b.WriteString("\"")
				b.WriteString(m.DataSets[i][j].Value.(net.HardwareAddr).String())
				b.WriteString("\"")
			case []uint8:
				b.WriteString("\"")
				b.WriteString(fmt.Sprintf("0x%x", m.DataSets[i][j].Value.([]uint8)))
				b.WriteString("\"")
			default:
				return errUknownMarshalDataType
			}
			if j < length-1 {
				b.WriteString("},")
			} else {
				b.WriteString("}")
			}
		}
		if i < dsLength-1 {
			b.WriteString("],")
		} else {
			b.WriteString("]")
		}
	}
	b.WriteString("]")

	return nil
}

func (m *Message) encodeHeader(b *bytes.Buffer) {
	b.WriteString("\"Header\":{\"Version\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.Version), 10))
	b.WriteString(",\"Length\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.Length), 10))
	b.WriteString(",\"ExportTime\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.ExportTime), 10))
	b.WriteString(",\"SequenceNo\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.SequenceNo), 10))
	b.WriteString(",\"DomainID\":")
	b.WriteString(strconv.FormatInt(int64(m.Header.DomainID), 10))
	b.WriteString("},")
}

func (m *Message) encodeAgent(b *bytes.Buffer) {
	b.WriteString("\"AgentID\":\"")
	b.WriteString(m.AgentID)
	b.WriteString("\",")
}