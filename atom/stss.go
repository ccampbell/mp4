package atom

import "encoding/binary"

// StssBox defines the stss box structure.
type StssBox struct {
	*Box
	Version byte
	Flags   uint32
	EntryCount   uint32
	SampleNumbers []uint32
}

func (b *StssBox) parse() error {
	data := b.ReadBoxData()

	b.Version = data[0]
	b.Flags = binary.BigEndian.Uint32(data[0:4])

	count := binary.BigEndian.Uint32(data[4:8])
	b.EntryCount = count

	b.SampleNumbers = make([]uint32, count)

	for i := 0; i < int(count); i++ {
		b.SampleNumbers[i] = binary.BigEndian.Uint32(data[(8 + 4*i):(12 + 4*i)])
	}

	return nil
}
