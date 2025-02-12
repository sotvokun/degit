package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

type Writer struct {
	file    *os.File
	tarball *tar.Writer
	gz      *gzip.Writer
}

func NewWriter(path string) (*Writer, error) {
	archive := &Writer{
		file:    nil,
		tarball: nil,
		gz:      nil,
	}

	if _, err := os.Stat(path); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(file)
	tarWriter := tar.NewWriter(gzipWriter)

	archive.file = file
	archive.tarball = tarWriter
	archive.gz = gzipWriter

	return archive, nil
}

func (w *Writer) Close() error {
	if err := w.tarball.Close(); err != nil {
		return err
	}
	if err := w.gz.Close(); err != nil {
		return err
	}
	if err := w.file.Close(); err != nil {
		return err
	}
	w.tarball = nil
	w.gz = nil
	w.file = nil
	return nil
}

func (w *Writer) File(file io.Reader, fileinfo os.FileInfo, filename string) error {
	header, err := tar.FileInfoHeader(fileinfo, "")
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filename)
	if err := w.tarball.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(w.tarball, file); err != nil {
		return err
	}
	return nil
}

func (w *Writer) Symlink(fileinfo os.FileInfo, filename string, target string) error {
	header, err := tar.FileInfoHeader(fileinfo, target)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filename)
	if err := w.tarball.WriteHeader(header); err != nil {
		return err
	}
	return nil
}
