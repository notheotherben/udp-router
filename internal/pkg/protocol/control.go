package protocol

import (
	"encoding/gob"
)

func init() {
	gob.Register(PathAdvertisement{})
}

type PathAdvertisement struct {
	Source Address
	Dest   Address
	Port   int
	Cost   int
}

const PathAdvertisementSubtype = PacketSubtype(1)
