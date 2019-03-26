package protocol

import "encoding/gob"

func init() {
	gob.Register(DataPayload{})
}

type DataPayload struct {
	Length int
	Data   []byte
}
