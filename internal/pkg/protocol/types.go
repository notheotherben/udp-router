package protocol

type PacketType int8
type PacketSubtype int8

const ControlPacketType = PacketType(1)
const DataPacketType = PacketType(2)
