package mavtunnel

import (
	"fmt"
	"regexp"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/standard"
)

var node *gomavlib.Node
var nodeWait = make(chan struct{})

const systemID = 210

type endpointType struct {
	args string
	desc string
	make func(args string) gomavlib.EndpointConf
}

var reArgs = regexp.MustCompile("^([a-z]+):(.+)$")
var endpointTypes = map[string]endpointType{
	"serial": {
		"port:baudrate",
		"serial",
		func(args string) gomavlib.EndpointConf {
			return gomavlib.EndpointSerial{Address: args}
		},
	},
	"udps": {
		"listen_ip:port",
		"udp, server mode",
		func(args string) gomavlib.EndpointConf {
			return gomavlib.EndpointUDPServer{Address: args}
		},
	},
	"udpc": {
		"dest_ip:port",
		"udp, client mode",
		func(args string) gomavlib.EndpointConf {
			return gomavlib.EndpointUDPClient{Address: args}
		},
	},
	"udpb": {
		"broadcast_ip:port",
		"udp, broadcast mode",
		func(args string) gomavlib.EndpointConf {
			return gomavlib.EndpointUDPBroadcast{BroadcastAddress: args}
		},
	},
	"tcps": {
		"listen_ip:port",
		"tcp, server mode",
		func(args string) gomavlib.EndpointConf {
			return gomavlib.EndpointTCPServer{Address: args}
		},
	},
	"tcpc": {
		"dest_ip:port",
		"tcp, client mode",
		func(args string) gomavlib.EndpointConf {
			return gomavlib.EndpointTCPClient{Address: args}
		},
	},
}

func Endpoint(endpoint string) (gomavlib.EndpointConf, error) {
	matches := reArgs.FindStringSubmatch(endpoint)
	if matches == nil {
		return nil, fmt.Errorf("invalid endpoint: %s", endpoint)
	}
	key, args := matches[1], matches[2]

	etype, ok := endpointTypes[key]
	if !ok {
		return nil, fmt.Errorf("invalid endpoint: %s", endpoint)
	}

	return etype.make(args), nil
}

//Blocks (always returns a non nil error)
func Connect(endpoint gomavlib.EndpointConf, capture func([]byte)) (err error) {
	node, err = gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints:        []gomavlib.EndpointConf{endpoint},
		Dialect:          standard.Dialect,
		OutVersion:       gomavlib.V2,
		HeartbeatDisable: true,
		OutSystemID:      systemID,
	})
	if err != nil {
		return fmt.Errorf("could not connect to MAVLink endpoint: %s", err)
	}
	close(nodeWait)
	var errChan = make(chan error)

	go (func() {
		errChan <- forwardPackets()
	})()

	for {
		select {
		case evt := <-node.Events():
			if frm, ok := evt.(*gomavlib.EventFrame); ok {
				readFrame(frm, capture)
			}
		case err := <-errChan:
			return err
		}
	}
}

func WaitForConnection() {
	<-nodeWait
}
