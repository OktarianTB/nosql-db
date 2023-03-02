package main

import "encoding/binary"

type meta struct {
	root         pgnum
	freelistPage pgnum
}

func newEmptyMeta() *meta {
	return &meta{}
}

func (meta *meta) serialize(buffer []byte) {
	pos := 0

	binary.LittleEndian.PutUint64(buffer[pos:], uint64(meta.root))
	pos += pageNumSize

	binary.LittleEndian.PutUint64(buffer[pos:], uint64(meta.freelistPage))
	pos += pageNumSize
}

func (meta *meta) deserialize(buffer []byte) {
	pos := 0

	meta.root = pgnum(binary.LittleEndian.Uint64(buffer[pos:]))
	pos += pageNumSize

	meta.freelistPage = pgnum(binary.LittleEndian.Uint64(buffer[pos:]))
	pos += pageNumSize
}
