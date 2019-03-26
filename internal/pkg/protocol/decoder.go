package protocol

import (
	"encoding/gob"
	"io"
)

type Decoder struct {
	gobDecoder *gob.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		gobDecoder: gob.NewDecoder(r),
	}
}

func (e *Decoder) Decode(into *Packet) error {
	return e.gobDecoder.Decode(into)
}
