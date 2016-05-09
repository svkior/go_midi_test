package main

import (
	"github.com/rakyll/portmidi"
	"log"
	"strings"
	"time"
)

func workMidi(portId portmidi.DeviceId) {
	out, _ := portmidi.NewOutputStream(portId, 1024, 0)
	ticker := time.NewTicker(time.Millisecond * 20)

	out.WriteShort(0xb0, 0x00, 0x00)

	var idx int64
	var row, col, vel int64
	var row1, col1 int64
	for _ = range ticker.C {
		switch (idx >> 6) & 0x3 {
		case 0x0:
			vel = 0
		case 0x1:
			vel = 1 + 2
		case 0x2:
			vel = 16 + 32
		case 0x3:
			vel = 1 + 2 + 16 + 32
		}
		col = (idx) & 0x07
		//row = ((idx >> 5) & 0x01) | ((idx >> 3) & 0x02) | ((idx >> 1) & 0x04)
		row = (idx >> 3) & 0x07
		col1 = (idx) & 0x07
		row1 = (idx >> 3) & 0x07
		out.WriteShort(0x90, col+row*16, vel)
		out.WriteShort(0x90, row1+col1*16, vel)
		idx++
	}
	out.Close()
}

func main() {
	portmidi.Initialize()
	midiCounts := portmidi.CountDevices()
	log.Println("Counts:", midiCounts)
	var idx portmidi.DeviceId
	for idx = 0; idx < portmidi.DeviceId(midiCounts); idx++ {
		di := portmidi.GetDeviceInfo(idx)
		if strings.Contains(di.Name, "Launchpad") && (di.IsOutputAvailable) {
			log.Printf("We Have Launchpad Output %d", idx)

			workMidi(idx)
		}
	}

}
