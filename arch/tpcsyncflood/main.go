// The program is considerably larger than just getting stats since I was trying to document and structure it for clarity in future reads :)
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

type pcap struct {
	MagicNumber uint32 //Magic Number
	Major       uint16 //Major Version
	Minor       uint16 //Minor Version
	LinkType    uint32 //Link Type
	fileHeader  []byte //raw data
	Packets     []Packet
}

type tcpstats struct {
	synctr    uint64
	synackctr uint64
	ackctr    uint64
}

// todo remove global tcpstats
var tcpStats = &tcpstats{}

func (p *pcap) String() string {
	return fmt.Sprintf("Magic Number: %d\nMajor Version: %d\nMinor Version: %d\nLink Type: %d\n", p.MagicNumber, p.Major, p.Minor, p.LinkType)
}

type Packet struct {
	header []byte //raw data
	data   []byte //raw data
	packetheader
	ipv4
	tcp
}

type packetheader struct {
	TimestampSeconds      uint32 //Timestamp Seconds
	TimestampMicroseconds uint32 //Timestamp Microseconds
	CapturedLength        uint32 //Captured Length
	OriginalLength        uint32 //Original Length
	ProtocolType          uint32 //Protocol Type

}
type ipv4 struct {
	Version        uint8   //Version (4 bits): Specifies the IP protocol version, which is always 4 for IPv4.
	HeaderLength   uint8   //Header Length (4 bits): Specifies the length of the IPv4 header in 32-bit words.
	TypeOfService  uint8   //Type of Service (8 bits): Specifies the type of service and its priority level.
	TotalLength    uint16  //Total Length (16 bits): Specifies the length of the entire IPv4 packet (header and data) in bytes.
	Identification uint16  //Identification (16 bits): Used to identify fragments of a larger IP packet.
	Flags          uint8   //Flags (3 bits): Used for fragmentation control.
	FragmentOffset uint16  //Fragment Offset (13 bits): Used to reassemble fragmented IP packets.
	TTL            uint8   //Time to Live (8 bits): Specifies the number of hops (routers) that a packet is allowed to traverse before being discarded.
	Protocol       uint8   //Protocol (8 bits): Specifies the protocol used in the data portion of the IP packet.
	Checksum       uint16  //Header Checksum (16 bits): Used to verify the integrity of the IPv4 header.
	SourceIP       [4]byte //Source IP Address (32 bits): Specifies the IP address of the sender.
	DestIP         [4]byte //Destination IP Address (32 bits): Specifies the IP address of the receiver.
	Options        []byte  //Options (variable): Specifies additional options that are not included in the IPv4 header.
}

func (i *ipv4) String() string {
	return fmt.Sprintf("Version: %d\nHeader Length: %d\nType of Service: %d\nTotal Length: %d\nIdentification: %d\nFlags: %d\nFragment Offset: %d\nTTL: %d\nProtocol: %d\nChecksum: %d\nSource IP: %v\nDestination IP: %v\n", i.Version, i.HeaderLength, i.TypeOfService, i.TotalLength, i.Identification, i.Flags, i.FragmentOffset, i.TTL, i.Protocol, i.Checksum, i.SourceIP, i.DestIP)

}

type tcp struct {
	SourcePort      uint16 //Source Port (16 bits): Specifies the port number of the application sending the data.
	DestinationPort uint16 //Destination Port (16 bits): Specifies the port number of the application receiving the data.
	SequenceNumber  uint32 //Sequence Number (32 bits): Specifies the sequence number of the first data byte in this segment.
	Acknowledgment  uint32 //Acknowledgment Number (32 bits): Specifies the sequence number of the next data byte that the sender of the ACK is expecting.
	Offset          uint8  //Data Offset (4 bits): Specifies the length of the TCP header in 32-bit words.
	Reserved        uint8  //Reserved (3 bits): Reserved for future use.
	Flags           byte   //Flags (8 bits): Specifies the TCP flags that are set.
	WindowSize      uint16 //Window Size (16 bits): Specifies the number of bytes that the sender of this segment is willing to accept.
	Checksum        uint16 //Checksum (16 bits): Used to verify the integrity of the TCP header and data.
	UrgentPointer   uint16 //Urgent Pointer (16 bits): Specifies the sequence number of the last data byte in a segment when the URG flag is set.
	Options         []byte //Options (variable): Specifies additional options that are not included in the TCP header.
	Data            []byte //Data (variable): Specifies the data that is being sent.
}

func (p *packetheader) Parse(raw []byte) {
	//parse the PacketHeader
	p.TimestampSeconds = binary.LittleEndian.Uint32(raw[:4]) //Timestamp Seconds
	p.TimestampMicroseconds = binary.LittleEndian.Uint32(raw[4:8])
	p.CapturedLength = binary.LittleEndian.Uint32(raw[8:12])
	p.OriginalLength = binary.LittleEndian.Uint32(raw[12:16])

}
func (p *packetheader) String() string {
	return fmt.Sprintf("Timestamp Seconds: %d\nTimestamp Microseconds: %d\nCaptured Length: %d\nOriginal Length: %d\n", p.TimestampSeconds, p.TimestampMicroseconds, p.CapturedLength, p.OriginalLength)
}

// read file syn-flood/synflood.pcap
func ReadFile() *pcap {
	file, err := os.Open("syn-flood/synflood.pcap")
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fsz := fileInfo.Size()
	fmt.Printf("File size is %d\n", fileInfo.Size())

	pcap := &pcap{}
	//read pcap file header
	pcap.fileHeader = make([]byte, 24)
	file.Read(pcap.fileHeader)
	pcap.MagicNumber = binary.LittleEndian.Uint32(pcap.fileHeader[:4]) //Magic Number
	pcap.Major = binary.LittleEndian.Uint16(pcap.fileHeader[4:6])      //Major Version
	pcap.Minor = binary.LittleEndian.Uint16(pcap.fileHeader[6:8])      //Minor Version
	pcap.LinkType = binary.LittleEndian.Uint32(pcap.fileHeader[20:24]) //Link Type
	fmt.Printf("%s\n", pcap)

	readPtr := int64(24)

	for readPtr < fsz {
		//read packets
		pkt := &Packet{packetheader: packetheader{}}
		pkt.header = make([]byte, 16)
		file.Read(pkt.header)
		readPtr += 16
		pktHdr := &(pkt.packetheader)
		pktHdr.Parse(pkt.header)
		//read packet data
		pkt.data = make([]byte, pktHdr.CapturedLength)
		file.Read(pkt.data)
		readPtr += int64(pktHdr.CapturedLength)
		pcap.Packets = append(pcap.Packets, *pkt)

	}
	//print total number of packets
	fmt.Printf("Total number of packets: %d\n", len(pcap.Packets))

	defer file.Close()
	return pcap

}
func main() {
	pcap := ReadFile()

	for _, pkt := range pcap.Packets {
		parseIPV4(&pkt)
	}
	//print tcp stats
	fmt.Printf("TCP Stats: %v\n", tcpStats)

}

// https://networklessons.com/cisco/ccna-routing-switching-icnd1-100-105/ipv4-packet-header
// Version (4 bits): Specifies the IP protocol version, which is always 4 for IPv4.
// Header Length (4 bits): Specifies the length of the IPv4 header in 32-bit words.
// The minimum length of an IP header is 20 bytes (5 words), and the maximum length is 60 bytes (15 words)
//
// Type of Service (8 bits): Specifies the type of service and its priority level.
// skip
// Total Length (16 bits): Specifies the length of the entire IPv4 packet (header and data) in bytes.
// Identification (16 bits): Used to identify fragments of a larger IP packet.
// Flags (3 bits): Used for fragmentation control.
// Fragment Offset (13 bits): Used to reassemble fragmented IP packets.
// skip
// Time to Live (8 bits): Specifies the number of hops (routers) that a packet is allowed to traverse before being discarded.
// Protocol (8 bits): Specifies the protocol used in the data portion of the IP packet.
// Header Checksum (16 bits): Used to verify the integrity of the IPv4 header.
// skip
// Source IP Address (32 bits): Specifies the IP address of the sender.
// Destination IP Address (32 bits): Specifies the IP address of the receiver.
func parseIPV4(pkt *Packet) {

	pkt.ProtocolType = binary.NativeEndian.Uint32(pkt.data[:4])
	buffer := bytes.NewBuffer(pkt.data[4:])
	bufLen := buffer.Len()

	//IMPORTANT: From here on its Network byte order which is BigEndian
	if pkt.ProtocolType == 2 {
		fmt.Printf("IPv4 Packet\n")
		ipv4 := &ipv4{}
		nextByte := buffer.Next(1)[0]

		ipv4.Version = nextByte >> 4

		ipv4.HeaderLength = (nextByte & 0x0f) * 4

		ipv4.TypeOfService = buffer.Next(1)[0]

		ipv4.TotalLength = binary.BigEndian.Uint16(buffer.Next(2))

		_ = buffer.Next(4)

		ipv4.TTL = buffer.Next(1)[0]

		ipv4.Protocol = buffer.Next(1)[0]

		_ = buffer.Next(2)

		ipv4.SourceIP[0], ipv4.SourceIP[1], ipv4.SourceIP[2], ipv4.SourceIP[3] = buffer.Next(1)[0], buffer.Next(1)[0], buffer.Next(1)[0], buffer.Next(1)[0]

		ipv4.DestIP[0], ipv4.DestIP[1], ipv4.DestIP[2], ipv4.DestIP[3] = buffer.Next(1)[0], buffer.Next(1)[0], buffer.Next(1)[0], buffer.Next(1)[0]
		hasOptions := (int(ipv4.HeaderLength) > (bufLen - buffer.Len()))

		fmt.Printf("hasOptions: %v\n", hasOptions)
		fmt.Printf("IPV4 Details %s\n", ipv4)
		pkt.ipv4 = *ipv4
		pkt.tcp = *parseTCP(*buffer)

	} else {
		fmt.Printf("Can't parse %v Packet\n", pkt.packetheader.ProtocolType)

	}

}

func parseTCP(buf bytes.Buffer) *tcp {
	tcp := &tcp{}
	tcp.SourcePort = binary.BigEndian.Uint16(buf.Next(2))
	tcp.DestinationPort = binary.BigEndian.Uint16(buf.Next(2))
	fmt.Printf("Source Port: %d\n", tcp.SourcePort)
	fmt.Printf("Destination Port: %d\n", tcp.DestinationPort)
	tcp.SequenceNumber = binary.LittleEndian.Uint32(buf.Next(4))
	tcp.Acknowledgment = binary.LittleEndian.Uint32(buf.Next(4))
	_ = buf.Next(1)
	//tcp.Offset = buf.Next(1)[0] >> 4
	//tcp.Reserved = buf.Next(1)[0] & 0x0f
	tcp.Flags = buf.Next(1)[0]
	//print all bits in flags in binary
	fmt.Printf("%b\n", tcp.Flags)

	switch {
	case tcp.Flags == 2:
		fmt.Printf("SYN ONLY\n")
		tcpStats.synctr++
	case tcp.Flags == 18:
		fmt.Printf("SYN ACK\n")
		tcpStats.synackctr++
	case tcp.Flags == 16:
		fmt.Printf("ACK ONLY\n")
		//fmt.Scanf("%s")
		tcpStats.ackctr++
	default:
		fmt.Printf("Other\n")
	}

	tcp.WindowSize = binary.BigEndian.Uint16(buf.Next(2))
	tcp.Checksum = binary.BigEndian.Uint16(buf.Next(2))
	tcp.UrgentPointer = binary.BigEndian.Uint16(buf.Next(2))
	return tcp
}
