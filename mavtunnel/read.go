package mavtunnel

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sync"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/standard"
)

func readFrame(frm *gomavlib.EventFrame, capture func([]byte)) {
	switch msg := frm.Message().(type) {
	case *standard.MessageTunnel:
		if msg.TargetSystem == uint8(systemID) && msg.TargetComponent == uint8(standard.MAV_COMP_ID_ALL) && msg.PayloadType == 400 {
			processFragment(msg.Payload, capture)
		}
	}
}

type packetInBuild struct {
	payload       []byte
	fragmentCount int
}

var packetMutex = sync.Mutex{}
var packets = make(map[uint32]*packetInBuild) // Maps reconstruction_identifier to packetInBuild

func processFragment(fragment [128]uint8, capture func([]byte)) {
	packetMutex.Lock()
	defer packetMutex.Unlock()

	reconstructionIdentifier := binary.LittleEndian.Uint32(fragment[0:4])
	sequenceNumber := int(fragment[4])
	totalLength := int(binary.LittleEndian.Uint32(fragment[5:9]))
	lengthOfPayload := minOf(119, totalLength-(sequenceNumber*119))
	payload := fragment[9 : lengthOfPayload+9]

	var packet *packetInBuild

	if p, ok := packets[reconstructionIdentifier]; ok {
		packet = p
	} else {
		packet = &packetInBuild{
			payload:       make([]byte, totalLength),
			fragmentCount: 0,
		}
		packets[reconstructionIdentifier] = packet
	}

	packet.fragmentCount = packet.fragmentCount + 1
	chunkStart := sequenceNumber * 119
	chunkEnd := (sequenceNumber * 119) + len(payload)
	if chunkStart > cap(packet.payload) || chunkEnd > cap(packet.payload) {
		fmt.Fprintf(os.Stderr, "Error when sending packet: Packet 0x%x had a fragmentation error\n", reconstructionIdentifier)
		delete(packets, reconstructionIdentifier)
		return
	}
	copy(packet.payload[chunkStart:chunkEnd], payload)

	if packet.fragmentCount == int(math.Ceil(float64(totalLength)/119.0)) {
		delete(packets, reconstructionIdentifier)
		capture(packet.payload)
	}
}
