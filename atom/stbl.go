package atom

// StblBox defines the stbl box structure.
type StblBox struct {
	*Box
	Stts *SttsBox
	Stsd *StsdBox
	Stss *StssBox
}

func (b *StblBox) parse() error {
	boxes := readBoxes(b.File, b.Start+BoxHeaderSize, b.Size-BoxHeaderSize)

	for _, box := range boxes {
		switch box.Name {
		case "stts":
			b.Stts = &SttsBox{Box: box}
			b.Stts.parse()
		case "stss":
			b.Stss = &StssBox{Box: box}
			b.Stss.parse()
		case "stsd":
			b.Stsd = &StsdBox{Box: box}
			b.Stsd.parse()
		}
	}
	return nil
}
