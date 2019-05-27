/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type OmciMsgHandler func(class OmciClass, content OmciContent, key OnuKey) ([]byte, error)

var Handlers = map[OmciMsgType]OmciMsgHandler{
	MibReset:         mibReset,
	MibUpload:        mibUpload,
	MibUploadNext:    mibUploadNext,
	Set:              set,
	Create:           create,
	Get:              get,
	GetAllAlarms:     getAllAlarms,
	GetAllAlarmsNext: getAllAlarmsNext,
	SynchronizeTime:  syncTime,
	Delete:           delete,
	Reboot:           reboot,
}

func mibReset(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	log.Printf("%v - Omci MibReset", key)
	if state, ok := OnuOmciStateMap[key]; ok {
		log.Printf("%v - Reseting OnuOmciState", key)
		state.ResetOnuOmciState()
	}

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return pkt, nil
}

func mibUpload(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	log.Printf("%v - Omci MibUpload", key)

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	pkt[9] = NumMibUploads // Number of subsequent MibUploadNext cmds

	return pkt, nil
}

func mibUploadNext(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	state := OnuOmciStateMap[key]

	// commandNumber is the "Command number" attribute received in "MIB Upload Next" OMCI message
	commandNumber := ( (uint16(content[1])) | (uint16(content[0])<<8) )
	log.Printf("%v - Omci MibUploadNext %d", key, commandNumber)

	switch commandNumber {
	case 0:
		// ONT Data (2)
		log.Println("ONT DATA")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 1:
		// Circuit Pack (6) - #1
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x01, 0xf0, 0x00, 0x2f, 0x04,
			0x49, 0x53, 0x4b, 0x54, 0x71, 0xe8, 0x00, 0x80,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 2:
		// Circuit Pack (6) - #2
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x01, 0x0f, 0x00, 0x42, 0x52,
			0x43, 0x4d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 3:
		// Circuit Pack (6) - #3
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x01, 0x00, 0xf8, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 4:
		// Circuit Pack (6) - #4
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x01, 0x00, 0x04, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 5:
		// Circuit Pack (6) - #5
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x80, 0xf0, 0x00, 0xee, 0x01,
			0x49, 0x53, 0x4b, 0x54, 0x71, 0xe8, 0x00, 0x80,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 6:
		// Circuit Pack (6) - #6
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x80, 0x0f, 0x00, 0x42, 0x52,
			0x43, 0x4d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 7:
		// Circuit Pack (6) - #7
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x80, 0x00, 0xf8, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20,
			0x20, 0x20, 0x00, 0x08, 0x40, 0x10, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 8:
		// Circuit Pack (6) - #8
		log.Println("Circuit Pack")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x06, 0x01, 0x80, 0x00, 0x04, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 9, 10, 11, 12:
		// PPTP (11)
		log.Println("PPTP")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x00, 0x0b, 0x01, 0x01, 0xff, 0xfe, 0x00, 0x2f,
			0x00, 0x00, 0x00, 0x00, 0x03, 0x05, 0xee, 0x00,
			0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		pkt[11] = state.pptpInstance // ME Instance
		state.pptpInstance++
	case 13, 14, 15, 16, 17, 18, 19, 20:
		// T-CONT (262)
		log.Println("T-CONT")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x06, 0x80, 0x00, 0xe0, 0x00, 0xff, 0xff,
			0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		pkt[11] = state.tcontInstance // TCONT ME Instance
		state.tcontInstance++
	case 21:
		// ANI-G (263)
		log.Println("ANI-G")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x07, 0x80, 0x01, 0xff, 0xff, 0x01, 0x00,
			0x08, 0x00, 0x30, 0x00, 0x00, 0x05, 0x09, 0x00,
			0x00, 0xe0, 0x54, 0xff, 0xff, 0x00, 0x00, 0x0c,
			0x63, 0x81, 0x81, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	case 22, 23, 24, 25:
		// UNI-G (264)
		log.Println("UNI-G")
		pkt = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x08, 0x01, 0x01, 0xf8, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		pkt[11] = state.uniGInstance // UNI-G ME Instance
		state.uniGInstance++

	case 26, 30, 34, 38, 42, 46, 50, 54:
		// Prior-Q with mask downstream
		log.Println("Mib-upload for prior-q with mask")
		// For downstream PQ, pkt[10] is 0x00
		// So the instanceId will be like 0x0001, 0x0002,... etc
		pkt = []byte{
			0x00, 0x42, 0x2e, 0x0a, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x15, 0x00, 0x00, 0x00, 0x0f, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		state.priorQInstance++
		pkt[11] = state.priorQInstance

	case 27, 31, 35, 39, 43, 47, 51, 55:
		// Prior-Q with attribute list downstream
		log.Println("Mib-upload for prior-q with attribute list")
		pkt = []byte{
			0x00, 0x43, 0x2e, 0x0a, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x15, 0x00, 0x00, 0xff, 0xf0, 0x00, 0x01,
			0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
			0x20, 0x00, 0x00, 0x01, 0x20, 0x01, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}

		pkt[11] = state.priorQInstance
		state.tcontInstance--
		pkt[24] = state.tcontInstance // related port points to tcont
		pkt[28] = state.tcontInstance

	case 28, 32, 36, 40, 44, 48, 52, 56:
		// Prior-Q with mask upstream
		log.Println("Mib-upload for prior-q with mask")
		pkt = []byte{
			0x00, 0x42, 0x2e, 0x0a, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x15, 0x80, 0x00, 0x00, 0x0f, 0xff, 0xff,
			0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}
		pkt[11] = state.priorQInstance

	case 29, 33, 37, 41, 45, 49, 53, 57:
		// Prior-Q with attribute list upstream
		log.Println("Mib-upload for prior-q with attribute list")
		// For upstream pkt[10] is fixed as 80
		// So for upstream PQ, instanceId will be like 0x8001, 0x8002 ... etc
		pkt = []byte{
			0x00, 0x43, 0x2e, 0x0a, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x15, 0x80, 0x00, 0xff, 0xf0, 0x00, 0x01,
			0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80,
			0x20, 0x00, 0x00, 0x80, 0x20, 0x01, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}

		pkt[11] = state.priorQInstance
		pkt[24] = state.tcontInstance // related port points to tcont
		pkt[28] = state.tcontInstance

	case 58, 59, 60, 61, 62, 63, 64, 65:
		// Traffic Scheduler
		log.Println("Traffic Scheduler")
		pkt = []byte{
			0x02, 0xa4, 0x2e, 0x0a, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x16, 0x80, 0x00, 0xf0, 0x00, 0x80, 0x00,
			0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}

		pkt[15] = state.tcontInstance
		state.tcontInstance++

	case 66:
		// ONT-2G
		log.Println("ONT-2G")
		pkt = []byte{
			0x00, 0x16, 0x2e, 0x0a, 0x00, 0x02, 0x00, 0x00,
			0x01, 0x01, 0x00, 0x00, 0x07, 0xfc, 0x00, 0x40,
			0x08, 0x01, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x7f, 0x00, 0x00, 0x3f, 0x00, 0x01, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}

	default:
		state.extraMibUploadCtr++
		errstr := fmt.Sprintf("%v - Invalid MibUpload request: %d, extras: %d", key, state.mibUploadCtr, state.extraMibUploadCtr)
		return nil, errors.New(errstr)
	}

	state.mibUploadCtr++
	return pkt, nil
}

func set(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci Set", key)

	return pkt, nil
}

func create(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	if class == GEMPortNetworkCTP {
		if onuOmciState, ok := OnuOmciStateMap[key]; !ok {
			log.Printf("%v - ONU Key Error", key)
			return nil, errors.New("ONU Key Error")
		} else {
			onuOmciState.gemPortId = binary.BigEndian.Uint16(content[:2])
			log.Printf("%v - Gem Port Id %d", key, onuOmciState.gemPortId)
			// FIXME
			OnuOmciStateMap[key].state = DONE
		}
	}

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x01, 0x10, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci Create", key)

	return pkt, nil
}

func get(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x2d, 0x02, 0x01,
		0x00, 0x20, 0xc0, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci Get", key)

	return pkt, nil
}

func getAllAlarms(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci GetAllAlarms", key)

	return pkt, nil
}

func syncTime(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci syncTime", key)

	return pkt, nil
}

func getAllAlarmsNext(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x0b, 0x01, 0x02, 0x80, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci GetAllAlarmsNext", key)

	return pkt, nil
}

func delete(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte

	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x0b, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci Delete", key)

	return pkt, nil
}

func reboot(class OmciClass, content OmciContent, key OnuKey) ([]byte, error) {
	var pkt []byte
	pkt = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	log.Printf("%v - Omci Reboot", key)
	return pkt, nil
}
