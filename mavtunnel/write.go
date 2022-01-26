package mavtunnel

import (
	"encoding/binary"
	"errors"

	"github.com/aler9/gomavlib/pkg/dialects/standard"
)

var packetQueue = make(chan []byte, 100)

func SendPacket(packet []byte) {
	packetQueue <- packet
}

//Blocks (always returns a non nil error)
func forwardPackets() error {
	var reconstructionIdentifier uint32
	for packet := range packetQueue {
		// The format of packet data over mavlink is as follows:
		// uint32_t (LE) reconstruction_identifier - The ID of that the packet is part of. Increments with every packet sent.
		// uint8_t sequence_number - The sequence number of the packet. Increments with every packet fragment sent.
		// uint32_t (LE) total_length - The total length of the packet.
		// uint8_t[119] payload - The payload of the packet fragment.
		//
		// The total number of expected fragments can be implicitly calculated from total_length / 119, rounded up.

		var sequenceNumber byte
		for i := 0; i < len(packet); i += 119 {
			var fragment [128]uint8

			ri := make([]byte, 4)
			binary.LittleEndian.PutUint32(ri, reconstructionIdentifier)
			copy(fragment[0:4], ri)

			fragment[4] = sequenceNumber

			tl := make([]byte, 4)
			binary.LittleEndian.PutUint32(tl, uint32(len(packet)))
			copy(fragment[5:9], tl)

			copy(fragment[9:], packet[i:minOf(i+119, len(packet))])

			node.WriteMessageAll(&standard.MessageTunnel{
				TargetSystem:    uint8(systemID),
				TargetComponent: uint8(standard.MAV_COMP_ID_ALL),
				PayloadType:     400,
				PayloadLength:   uint8(minOf(128, len(packet)-i)),
				Payload:         fragment,
			})

			sequenceNumber++
		}

		reconstructionIdentifier++
	}
	return errors.New("packet queue closed")
}

func minOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
