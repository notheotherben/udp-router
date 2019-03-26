package protocol

import (
	"encoding/gob"
	"io"
)

type Encoder struct {
	gobEncoder *gob.Encoder
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		gobEncoder: gob.NewEncoder(w),
	}
}

func (e *Encoder) Encode(packet *Packet) error {
	return e.gobEncoder.Encode(packet)
}
