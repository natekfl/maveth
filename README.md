# MAVETH - Ethernet over MAVLink

Allows an ethernet connection to a MAVLink network. Useful if the drone stack communicates over IP, only telemetry radios connect the drone to the ground.

## Setup

### Windows

The TAP Interface from OpenVPN is required. Grab the installer from https://openvpn.net/community-downloads/. Only the TAP Interface is needed. Everything else, including OpenVPN itself, can be deselected.

### Linux

No setup required

### MacOSX

Not supported

## Usage

Simply run `maveth bridge -m <endpoint>`. \<endpoint\> takes the form `port:baudrate`.

Devices on the network will need to either have static IPs, or use link-local addresses.