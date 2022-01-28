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

Simply run `maveth bridge -m <endpoint>`. \<endpoint\> can be one of:

    udps:listen_ip:port (udp, server mode)
    udpc:dest_ip:port (udp, client mode)
    udpb:broadcast_ip:port (udp, broadcast mode)
    tcps:listen_ip:port (tcp, server mode)
    tcpc:dest_ip:port (tcp, client mode)
    serial:port:baudrate (serial)

In general, only the serial type should be used. As all the others are network based, they have the possibility of creating an infinite loop if you don't know what you're doing.

**Devices on the network will need to either have static IPs, or use link-local addresses.** On windows, link local is automatic. On linux, it can be enabled using `sudo avahi-autoipd --force-bind -D maveth` (Tested on Ubunu)