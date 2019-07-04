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

package linux

// Netlink message types for NETLINK_ROUTE sockets, from uapi/linux/rtnetlink.h.
const (
	RTM_NEWLINK = 16
	RTM_DELLINK = 17
	RTM_GETLINK = 18
	RTM_SETLINK = 19

	RTM_NEWADDR = 20
	RTM_DELADDR = 21
	RTM_GETADDR = 22

	RTM_NEWROUTE = 24
	RTM_DELROUTE = 25
	RTM_GETROUTE = 26

	RTM_NEWNEIGH = 28
	RTM_DELNEIGH = 29
	RTM_GETNEIGH = 30

	RTM_NEWRULE = 32
	RTM_DELRULE = 33
	RTM_GETRULE = 34

	RTM_NEWQDISC = 36
	RTM_DELQDISC = 37
	RTM_GETQDISC = 38

	RTM_NEWTCLASS = 40
	RTM_DELTCLASS = 41
	RTM_GETTCLASS = 42

	RTM_NEWTFILTER = 44
	RTM_DELTFILTER = 45
	RTM_GETTFILTER = 46

	RTM_NEWACTION = 48
	RTM_DELACTION = 49
	RTM_GETACTION = 50

	RTM_NEWPREFIX = 52

	RTM_GETMULTICAST = 58

	RTM_GETANYCAST = 62

	RTM_NEWNEIGHTBL = 64
	RTM_GETNEIGHTBL = 66
	RTM_SETNEIGHTBL = 67

	RTM_NEWNDUSEROPT = 68

	RTM_NEWADDRLABEL = 72
	RTM_DELADDRLABEL = 73
	RTM_GETADDRLABEL = 74

	RTM_GETDCB = 78
	RTM_SETDCB = 79

	RTM_NEWNETCONF = 80
	RTM_GETNETCONF = 82

	RTM_NEWMDB = 84
	RTM_DELMDB = 85
	RTM_GETMDB = 86

	RTM_NEWNSID = 88
	RTM_DELNSID = 89
	RTM_GETNSID = 90
)

// InterfaceInfoMessage is struct ifinfomsg, from uapi/linux/rtnetlink.h.
type InterfaceInfoMessage struct {
	Family  uint8
	Padding uint8
	Type    uint16
	Index   int32
	Flags   uint32
	Change  uint32
}

// Interface flags, from uapi/linux/if.h.
const (
	IFF_UP          = 1 << 0
	IFF_BROADCAST   = 1 << 1
	IFF_DEBUG       = 1 << 2
	IFF_LOOPBACK    = 1 << 3
	IFF_POINTOPOINT = 1 << 4
	IFF_NOTRAILERS  = 1 << 5
	IFF_RUNNING     = 1 << 6
	IFF_NOARP       = 1 << 7
	IFF_PROMISC     = 1 << 8
	IFF_ALLMULTI    = 1 << 9
	IFF_MASTER      = 1 << 10
	IFF_SLAVE       = 1 << 11
	IFF_MULTICAST   = 1 << 12
	IFF_PORTSEL     = 1 << 13
	IFF_AUTOMEDIA   = 1 << 14
	IFF_DYNAMIC     = 1 << 15
	IFF_LOWER_UP    = 1 << 16
	IFF_DORMANT     = 1 << 17
	IFF_ECHO        = 1 << 18
)

// Interface link attributes, from uapi/linux/if_link.h.
const (
	IFLA_UNSPEC          = 0
	IFLA_ADDRESS         = 1
	IFLA_BROADCAST       = 2
	IFLA_IFNAME          = 3
	IFLA_MTU             = 4
	IFLA_LINK            = 5
	IFLA_QDISC           = 6
	IFLA_STATS           = 7
	IFLA_COST            = 8
	IFLA_PRIORITY        = 9
	IFLA_MASTER          = 10
	IFLA_WIRELESS        = 11
	IFLA_PROTINFO        = 12
	IFLA_TXQLEN          = 13
	IFLA_MAP             = 14
	IFLA_WEIGHT          = 15
	IFLA_OPERSTATE       = 16
	IFLA_LINKMODE        = 17
	IFLA_LINKINFO        = 18
	IFLA_NET_NS_PID      = 19
	IFLA_IFALIAS         = 20
	IFLA_NUM_VF          = 21
	IFLA_VFINFO_LIST     = 22
	IFLA_STATS64         = 23
	IFLA_VF_PORTS        = 24
	IFLA_PORT_SELF       = 25
	IFLA_AF_SPEC         = 26
	IFLA_GROUP           = 27
	IFLA_NET_NS_FD       = 28
	IFLA_EXT_MASK        = 29
	IFLA_PROMISCUITY     = 30
	IFLA_NUM_TX_QUEUES   = 31
	IFLA_NUM_RX_QUEUES   = 32
	IFLA_CARRIER         = 33
	IFLA_PHYS_PORT_ID    = 34
	IFLA_CARRIER_CHANGES = 35
	IFLA_PHYS_SWITCH_ID  = 36
	IFLA_LINK_NETNSID    = 37
	IFLA_PHYS_PORT_NAME  = 38
	IFLA_PROTO_DOWN      = 39
	IFLA_GSO_MAX_SEGS    = 40
	IFLA_GSO_MAX_SIZE    = 41
)

// InterfaceAddrMessage is struct ifaddrmsg, from uapi/linux/if_addr.h.
type InterfaceAddrMessage struct {
	Family    uint8
	PrefixLen uint8
	Flags     uint8
	Scope     uint8
	Index     uint32
}

// Interface attributes, from uapi/linux/if_addr.h.
const (
	IFA_UNSPEC    = 0
	IFA_ADDRESS   = 1
	IFA_LOCAL     = 2
	IFA_LABEL     = 3
	IFA_BROADCAST = 4
	IFA_ANYCAST   = 5
	IFA_CACHEINFO = 6
	IFA_MULTICAST = 7
	IFA_FLAGS     = 8
)

// Device types, from uapi/linux/if_arp.h.
const (
	ARPHRD_LOOPBACK = 772
)

// RouteMessage struct rtmsg, from uapi/linux/rtnetlink.h
type RouteMessage struct {
	Family uint8
	DstLen uint8
	SrcLen uint8
	Tos    uint8

	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8

	Flags uint32
}

// Route types, from uapi/linux/rtnetlink.h
const (
	RTN_UNSPEC      = 0
	RTN_UNICAST     = 1  // Gateway or direct route
	RTN_LOCAL       = 2  // Accept locally
	RTN_BROADCAST   = 3  // Accept locally as broadcast, send as broadcast
	RTN_ANYCAST     = 6  // Accept locally as broadcast, but send as unicast
	RTN_MULTICAST   = 5  // Multicast route
	RTN_BLACKHOLE   = 6  // Drop
	RTN_UNREACHABLE = 7  // Destination is unreachable
	RTN_PROHIBIT    = 8  // Administratively prohibited
	RTN_THROW       = 9  // Not in this table
	RTN_NAT         = 10 // Translate this address
	RTN_XRESOLVE    = 11 // Use external resolver
)

// Route protocols/origins, from uapi/linux/rtnetlink.h
const (
	RTPROT_UNSPEC   = 0
	RTPROT_REDIRECT = 1   // Route installed by ICMP redirects
	RTPROT_KERNEL   = 2   // Route installed by kernel
	RTPROT_BOOT     = 3   // Route installed during boot
	RTPROT_STATIC   = 4   // Route installed by administrator
	RTPROT_GATED    = 8   // Apparently, GateD
	RTPROT_RA       = 9   // RDISC/ND router advertisements
	RTPROT_MRT      = 10  // Merit MRT
	RTPROT_ZEBRA    = 11  // Zebra
	RTPROT_BIRD     = 12  // BIRD
	RTPROT_DNROUTED = 13  // DECnet routing daemon
	RTPROT_XORP     = 14  // XORP
	RTPROT_NTK      = 15  // Netsukuku
	RTPROT_DHCP     = 16  // DHCP client
	RTPROT_MROUTED  = 17  // Multicast daemon
	RTPROT_BABEL    = 42  // Babel daemon
	RTPROT_BGP      = 186 // BGP Routes
	RTPROT_ISIS     = 187 // ISIS Routes
	RTPROT_OSPF     = 188 // OSPF Routes
	RTPROT_RIP      = 189 // RIP Routes
	RTPROT_EIGRP    = 192 // EIGRP Routes
)

// Route scopes, from uapi/linux/rtnetlink.h
const (
	RT_SCOPE_UNIVERSE = 0   // global route
	RT_SCOPE_SITE     = 200 // interior route in the local autonomous system
	RT_SCOPE_LINK     = 253 // route on this link
	RT_SCOPE_HOST     = 254 // route on the local host
	RT_SCOPE_NOWHERE  = 255 // destination doesn't exist
)

// Route flags, from uapi/linux/rtnetlink.h
const (
	RTM_F_NOTIFY       = 0x100
	RTM_F_CLONED       = 0x200
	RTM_F_EQUALIZE     = 0x400
	RTM_F_PREFIX       = 0x800
	RTM_F_LOOKUP_TABLE = 0x1000
	RTM_F_FIB_MATCH    = 0x2000
)

// Route tables, from uapi/linux/rtnetlink.h
const (
	RT_TABLE_UNSPEC  = 0
	RT_TABLE_COMPAT  = 252
	RT_TABLE_DEFAULT = 253
	RT_TABLE_MAIN    = 254
	RT_TABLE_LOCAL   = 255
)

// Route attributes, from uapi/linux/rtnetlink.h
const (
	RTA_UNSPEC        = 0
	RTA_DST           = 1
	RTA_SRC           = 2
	RTA_IIF           = 3
	RTA_OIF           = 4
	RTA_GATEWAY       = 5
	RTA_PRIORITY      = 6
	RTA_PREFSRC       = 7
	RTA_METRICS       = 8
	RTA_MULTIPATH     = 9
	RTA_PROTOINFO     = 10 // no longer used
	RTA_FLOW          = 11
	RTA_CACHEINFO     = 12
	RTA_SESSION       = 13 // no longer used
	RTA_MP_ALGO       = 14 // no longer used
	RTA_TABLE         = 15
	RTA_MARK          = 16
	RTA_MFC_STATS     = 17
	RTA_VIA           = 18
	RTA_NEWDST        = 19
	RTA_PREF          = 20
	RTA_ENCAP_TYPE    = 21
	RTA_ENCAP         = 22
	RTA_EXPIRES       = 23
	RTA_PAD           = 24
	RTA_UID           = 25
	RTA_TTL_PROPAGATE = 26
	RTA_IP_PROTO      = 27
	RTA_SPORT         = 28
	RTA_DPORT         = 29
)
