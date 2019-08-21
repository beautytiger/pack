package blob

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"

	"github.com/docker/docker/pkg/ioutils"
	"github.com/pkg/errors"

	"github.com/buildpack/pack/internal/archive"
)

type Blob interface {
	Open() (io.ReadCloser, error)
}

type blob struct {
	Path string
}

func NewBlob(path string) Blob {
	return &blob{Path: path}
}

// Open returns an io.ReadCloser whose contents are in tar archive format
func (b blob) Open() (io.ReadCloser, error) {
	fi, err := os.Stat(b.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "read blob at path '%s'", b.Path)
	}
	if fi.IsDir() {
		return archive.ReadDirAsTar(b.Path, ".", 0, 0, -1), nil
	}
	fh, err := os.Open(b.Path)
	if err != nil {
		return nil, errors.Wrap(err, "open buildpack archive")
	}
	if ok, err := isGZip(fh); err != nil {
		return nil, errors.Wrap(err, "check header")
	} else if !ok {
		return fh, nil
	}
	gzr, err := gzip.NewReader(fh)
	if err != nil {
		return nil, errors.Wrap(err, "create gzip reader")
	}
	rc := ioutils.NewReadCloserWrapper(gzr, func() error {
		if err := gzr.Close(); err != nil {
			return err
		}
		return fh.Close()
	})
	return rc, nil
}

func isGZip(file *os.File) (bool, error) {
	b := make([]byte, 3)
	_, err := file.Read(b)
	if err != nil && err != io.EOF {
		return false, err
	} else if err == io.EOF {
		return false, nil
	}
	if _, err := file.Seek(0, 0); err != nil {
		return false, err
	}
	return bytes.Equal(b, []byte("\x1f\x8b\x08")), nil
}
