package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

type Reader struct {
	file    *os.File
	tarball *tar.Reader
	gz      *gzip.Reader
}

func NewReader(path string) (*Reader, error) {
	archive := &Reader{
		file:    nil,
		tarball: nil,
		gz:      nil,
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(gzipReader)

	archive.file = file
	archive.tarball = tarReader
	archive.gz = gzipReader

	return archive, nil
}

func (r *Reader) Close() error {
	if err := r.file.Close(); err != nil {
		return err
	}
	r.file = nil
	r.tarball = nil
	r.gz = nil
	return nil
}

func (r *Reader) Uncompress(dest string) error {
	var header *tar.Header
	header, err := r.tarball.Next()
	for ; err == nil; header, err = r.tarball.Next() {
		path := dest + "/" + header.Name
		switch header.Typeflag {
		case tar.TypeReg:
			if err = r.uncompressFile(path); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err = r.uncompressSymlink(header.Linkname, path); err != nil {
				return err
			}
		default:
			continue
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}

func (r *Reader) uncompressFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, r.tarball); err != nil {
		return err
	}
	return nil
}

func (r *Reader) uncompressSymlink(linkName string, path string) error {
	return os.Symlink(linkName, path)
}
