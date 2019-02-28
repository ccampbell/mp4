package atom

import (
	"encoding/binary"
	"fmt"
	"os"
	"github.com/DHowett/ranger"
)

const (
	// BoxHeaderSize Size of box header.
	BoxHeaderSize = int64(8)
)

type Reader struct {
	File interface{}
}

func (sr *Reader) Close() error {
	if _, ok := sr.File.(*ranger.Reader); ok {
		return nil
	}

	return sr.File.(*os.File).Close()
}

func (sr *Reader) Read(p []byte) (n int, err error) {
	if r, ok := sr.File.(*ranger.Reader); ok {
		return r.Read(p)
	}

	return sr.File.(*os.File).Read(p)
}

func (sr *Reader) ReadAt(p []byte, pos int64) (n int, err error) {
	if r, ok := sr.File.(*ranger.Reader); ok {
		return r.ReadAt(p, pos)
	}

	return sr.File.(*os.File).ReadAt(p, pos)
}

func (sr *Reader) GetSize() (s int64, err error) {
	if r, ok := sr.File.(*ranger.Reader); ok {
		return r.Length()
	}

	info, err := sr.File.(*os.File).Stat()
	if err != nil {
		return
	}

	return info.Size(), nil
}

// File defines a file structure.
type File struct {
	*Reader
	Ftyp *FtypBox
	Moov *MoovBox
	Mdat *MdatBox
	Size int64

	IsFragmented bool
}

// Parse parses an MP4 file for atom boxes.
func (f *File) Parse() (err error) {
	if f.Size, err = f.GetSize(); err != nil {
		return
	}

	boxes := readBoxes(f, int64(0), f.Size)
	for _, box := range boxes {
		switch box.Name {
		case "ftyp":
			f.Ftyp = &FtypBox{Box: box}
			f.Ftyp.parse()
		case "wide":
			// fmt.Println("found wide")
		case "mdat":
			f.Mdat = &MdatBox{Box: box}
			// No mdat boxes to parse
		case "moov":
			f.Moov = &MoovBox{Box: box}
			f.Moov.parse()

			f.IsFragmented = f.Moov.IsFragmented
		}
	}
	return
}

// ReadBoxAt reads a box from an offset.
func (f *File) ReadBoxAt(offset int64) (boxSize uint32, boxType string) {
	buf := f.ReadBytesAt(BoxHeaderSize, offset)
	boxSize = binary.BigEndian.Uint32(buf[0:4])
	offset += BoxHeaderSize

	boxType = string(buf[4:8])
	return boxSize, boxType
}

// ReadBytesAt reads a box at n and offset.
func (f *File) ReadBytesAt(n int64, offset int64) (word []byte) {
	buf := make([]byte, n)
	if _, error := f.ReadAt(buf, offset); error != nil {
		fmt.Println(error)
		return
	}
	return buf
}

// Box defines an Atom Box structure.
type Box struct {
	Name        string
	Size, Start int64
	File        *File
}

// ReadBoxData reads the box data from an atom box.
func (b *Box) ReadBoxData() []byte {
	if b.Size <= BoxHeaderSize {
		return nil
	}
	return b.File.ReadBytesAt(b.Size-BoxHeaderSize, b.Start+BoxHeaderSize)
}

func readBoxes(f *File, start int64, n int64) (l []*Box) {
	for offset := start; offset < start+n; {
		size, name := f.ReadBoxAt(offset)

		b := &Box{
			Name:  string(name),
			Size:  int64(size),
			File:  f,
			Start: offset,
		}

		l = append(l, b)
		offset += int64(size)
	}
	return l
}
