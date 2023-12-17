package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

type ArchiveWriter struct {
	File    *os.File
	Tarball *tar.Writer
	Archive *gzip.Writer
}

func createArchive(path string) (ArchiveWriter, error) {
	DEFAULT_ARCHIVE := ArchiveWriter{
		File:    nil,
		Tarball: nil,
	}

	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return DEFAULT_ARCHIVE, err
	}

	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return DEFAULT_ARCHIVE, err
	}

	file, err := os.Create(path)
	if err != nil {
		return DEFAULT_ARCHIVE, err
	}

	gziper := gzip.NewWriter(file)
	writer := tar.NewWriter(gziper)
	return ArchiveWriter{
		File:    file,
		Tarball: writer,
		Archive: gziper,
	}, nil
}

func (writer *ArchiveWriter) Add(file billy.File, fileinfo *os.FileInfo, filename string) error {
	header, err := tar.FileInfoHeader(*fileinfo, "")
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filename)
	if err := writer.Tarball.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(writer.Tarball, file); err != nil {
		return err
	}

	return nil
}

func (writer *ArchiveWriter) Symlink(fileinfo *os.FileInfo, filename string, target string) error {
	header, err := tar.FileInfoHeader(*fileinfo, target)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(filename)
	if err := writer.Tarball.WriteHeader(header); err != nil {
		return err
	}
	return nil
}

func (writer *ArchiveWriter) Close() error {
	if err := writer.Tarball.Close(); err != nil {
		return err
	}
	if err := writer.Archive.Close(); err != nil {
		return err
	}
	if err := writer.File.Close(); err != nil {
		return err
	}
	return nil
}
