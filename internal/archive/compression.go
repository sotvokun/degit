package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

type Archive struct {
	file    *os.File
	tarball *tar.Writer
	gz      *gzip.Writer
}

func New(path string) (*Archive, error) {
	archive_ := &Archive{
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

	archive_.file = file
	archive_.tarball = tarWriter
	archive_.gz = gzipWriter

	return archive_, nil
}

func (a *Archive) Add(file io.Reader, fileinfo os.FileInfo, filename string) error {
	header, err := tar.FileInfoHeader(fileinfo, "")
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filename)
	if err := a.tarball.WriteHeader(header); err != nil {
		return err
	}
	if _, err := io.Copy(a.tarball, file); err != nil {
		return err
	}
	return nil
}

func (a *Archive) Symlink(fileinfo os.FileInfo, filename string, target string) error {
	header, err := tar.FileInfoHeader(fileinfo, target)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filename)
	if err := a.tarball.WriteHeader(header); err != nil {
		return err
	}
	return nil
}

func (a *Archive) Close() error {
	if err := a.tarball.Close(); err != nil {
		return err
	}
	if err := a.gz.Close(); err != nil {
		return err
	}
	if err := a.file.Close(); err != nil {
		return err
	}
	return nil
}
