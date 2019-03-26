package protocol

type Address int

type PacketHeader struct {
	Type    PacketType
	Subtype PacketSubtype

	Source Address
	Dest   Address
}

type Packet struct {
	PacketHeader

	Payload interface{}
}
