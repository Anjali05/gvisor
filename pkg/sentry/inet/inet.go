// Copyright 2018 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package inet defines semantics for IP stacks.
package inet

// Stack represents a TCP/IP stack.
type Stack interface {
	// Interfaces returns all network interfaces as a mapping from interface
	// indexes to interface properties. Interface indices are strictly positive
	// integers.
	Interfaces() map[int32]Interface

	// InterfaceAddrs returns all network interface addresses as a mapping from
	// interface indexes to a slice of associated interface address properties.
	InterfaceAddrs() map[int32][]InterfaceAddr

	// SupportsIPv6 returns true if the stack supports IPv6 connectivity.
	SupportsIPv6() bool

	// TCPReceiveBufferSize returns TCP receive buffer size settings.
	TCPReceiveBufferSize() (TCPBufferSize, error)

	// SetTCPReceiveBufferSize attempts to change TCP receive buffer size
	// settings.
	SetTCPReceiveBufferSize(size TCPBufferSize) error

	// TCPSendBufferSize returns TCP send buffer size settings.
	TCPSendBufferSize() (TCPBufferSize, error)

	// SetTCPSendBufferSize attempts to change TCP send buffer size settings.
	SetTCPSendBufferSize(size TCPBufferSize) error

	// TCPSACKEnabled returns true if RFC 2018 TCP Selective Acknowledgements
	// are enabled.
	TCPSACKEnabled() (bool, error)

	// SetTCPSACKEnabled attempts to change TCP selective acknowledgement
	// settings.
	SetTCPSACKEnabled(enabled bool) error

	// RouteTable returns the
	RouteTable() []Route
}

// Interface contains information about a network interface.
type Interface struct {
	// Keep these fields sorted in the order they appear in rtnetlink(7).

	// DeviceType is the device type, a Linux ARPHRD_* constant.
	DeviceType uint16

	// Flags is the device flags; see netdevice(7), under "Ioctls",
	// "SIOCGIFFLAGS, SIOCSIFFLAGS".
	Flags uint32

	// Name is the device name.
	Name string

	// Addr is the hardware device address.
	Addr []byte

	// MTU is the maximum transmission unit.
	MTU uint32
}

// InterfaceAddr contains information about a network interface address.
type InterfaceAddr struct {
	// Keep these fields sorted in the order they appear in rtnetlink(7).

	// Family is the address family, a Linux AF_* constant.
	Family uint8

	// PrefixLen is the address prefix length.
	PrefixLen uint8

	// Flags is the address flags.
	Flags uint8

	// Addr is the actual address.
	Addr []byte
}

// TCPBufferSize contains settings controlling TCP buffer sizing.
//
// +stateify savable
type TCPBufferSize struct {
	// Min is the minimum size.
	Min int

	// Default is the default size.
	Default int

	// Max is the maximum size.
	Max int
}

// Route contains information about a network route.
type Route struct {
	// Keep these fields sorted in the order they appear in rtnetlink(7).

	// Family is the address family, a Linux AF_* constant.
	Family uint8

	// DstLen is the length of the destination address.
	DstLen uint8

	// SrcLen is the length of the source address.
	SrcLen uint8

	// Tos is the TOS filter
	Tos uint8

	// Table is the routing table ID.
	Table uint8

	// Protocol is the route origin, a Linux RTPROT_* constant.
	Protocol uint8

	// Scope is the distance to destination, a Linux RT_SCOPE_* constant.
	Scope uint8

	// Type is the route origin, a Linux RTN_* constant.
	Type uint8

	// Flags are route flags. See rtnetlink(7) under "rtm_flags".
	Flags uint32

	// DstAddr is the route destination address (RTA_DST).
	DstAddr []byte

	// SrcAddr is the route source address (RTA_SRC).
	SrcAddr []byte

	// OutputInterface is the output interface index (RTA_OIF)
	OutputInterface int32

	// GatewayAddr is the route gateway address (RTA_GATEWAY).
	GatewayAddr []byte

	// Metrics is the route metric (RTA_METRICS)
	Metrics int32
}
