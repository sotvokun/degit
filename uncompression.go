package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func uncompressArchive(reader io.Reader, dest string) error {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)
	var header *tar.Header
	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		if header.Typeflag != tar.TypeReg {
			continue
		}
		path := filepath.Join(dest, header.Name)

		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %s", err)
		}

		outFile, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("error creating file: %s", err)
		}
		if _, err := io.Copy(outFile, tarReader); err != nil {
			outFile.Close()
			return fmt.Errorf("error writing file: %s", err)
		}
		if err := outFile.Close(); err != nil {
			return fmt.Errorf("error closing file: %s", err)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("error reading tar: %s", err)
	}
	return nil
}
