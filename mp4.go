package mp4

import (
	"os"
	"strings"
	"net/url"
	"github.com/alfg/mp4/atom"
	"github.com/DHowett/ranger"
)

func getReaderForPath(path string) (r *atom.Reader, err error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		u, _ := url.Parse(path)
		ranger, err := ranger.NewReader(&ranger.HTTPRanger{URL: u})
		if err != nil {
			return nil, err
		}

		r = &atom.Reader{
			File: ranger,
		}

		return r, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	r = &atom.Reader{
		File: file,
	}

	return
}

// Open opens a file and returns a &File{}.
func Open(path string) (f *atom.File, err error) {
	reader, err := getReaderForPath(path)

	f = &atom.File{
		Reader: reader,
	}

	return f, f.Parse()
}
