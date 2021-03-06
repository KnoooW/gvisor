// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package header

import (
	"encoding/binary"

	"gvisor.googlesource.com/gvisor/pkg/tcpip"
)

const (
	udpSrcPort  = 0
	udpDstPort  = 2
	udpLength   = 4
	udpChecksum = 6
)

// UDPFields contains the fields of a UDP packet. It is used to describe the
// fields of a packet that needs to be encoded.
type UDPFields struct {
	// SrcPort is the "source port" field of a UDP packet.
	SrcPort uint16

	// DstPort is the "destination port" field of a UDP packet.
	DstPort uint16

	// Length is the "length" field of a UDP packet.
	Length uint16

	// Checksum is the "checksum" field of a UDP packet.
	Checksum uint16
}

// UDP represents a UDP header stored in a byte array.
type UDP []byte

const (
	// UDPMinimumSize is the minimum size of a valid UDP packet.
	UDPMinimumSize = 8

	// UDPProtocolNumber is UDP's transport protocol number.
	UDPProtocolNumber tcpip.TransportProtocolNumber = 17
)

// SourcePort returns the "source port" field of the udp header.
func (b UDP) SourcePort() uint16 {
	return binary.BigEndian.Uint16(b[udpSrcPort:])
}

// DestinationPort returns the "destination port" field of the udp header.
func (b UDP) DestinationPort() uint16 {
	return binary.BigEndian.Uint16(b[udpDstPort:])
}

// Length returns the "length" field of the udp header.
func (b UDP) Length() uint16 {
	return binary.BigEndian.Uint16(b[udpLength:])
}

// Payload returns the data contained in the UDP datagram.
func (b UDP) Payload() []byte {
	return b[UDPMinimumSize:]
}

// Checksum returns the "checksum" field of the udp header.
func (b UDP) Checksum() uint16 {
	return binary.BigEndian.Uint16(b[udpChecksum:])
}

// SetSourcePort sets the "source port" field of the udp header.
func (b UDP) SetSourcePort(port uint16) {
	binary.BigEndian.PutUint16(b[udpSrcPort:], port)
}

// SetDestinationPort sets the "destination port" field of the udp header.
func (b UDP) SetDestinationPort(port uint16) {
	binary.BigEndian.PutUint16(b[udpDstPort:], port)
}

// SetChecksum sets the "checksum" field of the udp header.
func (b UDP) SetChecksum(checksum uint16) {
	binary.BigEndian.PutUint16(b[udpChecksum:], checksum)
}

// CalculateChecksum calculates the checksum of the udp packet, given the total
// length of the packet and the checksum of the network-layer pseudo-header
// (excluding the total length) and the checksum of the payload.
func (b UDP) CalculateChecksum(partialChecksum uint16, totalLen uint16) uint16 {
	// Add the length portion of the checksum to the pseudo-checksum.
	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, totalLen)
	checksum := Checksum(tmp, partialChecksum)

	// Calculate the rest of the checksum.
	return Checksum(b[:UDPMinimumSize], checksum)
}

// Encode encodes all the fields of the udp header.
func (b UDP) Encode(u *UDPFields) {
	binary.BigEndian.PutUint16(b[udpSrcPort:], u.SrcPort)
	binary.BigEndian.PutUint16(b[udpDstPort:], u.DstPort)
	binary.BigEndian.PutUint16(b[udpLength:], u.Length)
	binary.BigEndian.PutUint16(b[udpChecksum:], u.Checksum)
}
