//go:build (linux && ignore) || ebpf
// +build linux,ignore ebpf

package dnsflow

import (
	"encoding/binary"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Weave
// Copyright 2016 Weaveworks Ltd.
// https://github.com/weaveworks/scope

// Gopacket doesn't provide direct support for DNS over TCP, see https://github.com/google/gopacket/issues/236
type tcpWithDNSSupport struct {
	tcp layers.TCP
}

func (m *tcpWithDNSSupport) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return m.tcp.DecodeFromBytes(data, df)
}

func (m *tcpWithDNSSupport) CanDecode() gopacket.LayerClass { return m.tcp.CanDecode() }

// Determine if a TCP segment contains a full DNS message (i.e. not fragmented).
func (m *tcpWithDNSSupport) hasSelfContainedDNSPayload() bool {
	payload := m.LayerPayload()
	if len(payload) < 2 {
		return false
	}

	// Assume it's a self-contained DNS message if the Length field
	// matches the length of the TCP segment
	dnsLengthField := binary.BigEndian.Uint16(payload)
	return int(dnsLengthField) == len(payload)-2
}

func (m *tcpWithDNSSupport) NextLayerType() gopacket.LayerType {
	// TODO: deal with TCP fragmentation and out-of-order segments
	if (m.tcp.SrcPort == 53 || m.tcp.DstPort == 53) && m.hasSelfContainedDNSPayload() {
		return layers.LayerTypeDNS
	}
	return m.tcp.NextLayerType()
}

func (m *tcpWithDNSSupport) LayerPayload() []byte {
	payload := m.tcp.LayerPayload()
	if len(payload) > 1 && (m.tcp.SrcPort == 53 || m.tcp.DstPort == 53) {
		// Omit the DNS length field, only included
		// in TCP, in order to reuse the DNS UDP parser
		payload = payload[2:]
	}
	return payload
}
