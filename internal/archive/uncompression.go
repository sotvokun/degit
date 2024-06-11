package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Uncompress(file io.Reader, dest string) error {
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzReader)
	var header *tar.Header
	header, err = tarReader.Next()
	for ; err == nil; header, err = tarReader.Next() {
		path := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeReg:
			if err = uncompressFile(path, tarReader); err != nil {
				return fmt.Errorf("error while uncompressing file: %s", err)
			}
		case tar.TypeSymlink:
			if err = uncompressSymlink(header.Linkname, path); err != nil {
				return fmt.Errorf("error while uncompressing symlink: %s", err)
			}
		default:
			continue
		}
	}
	if err != io.EOF {
		return fmt.Errorf("error while reading tar: %s", err)
	}
	return nil
}

func uncompressFile(path string, reader io.Reader) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, reader); err != nil {
		return err
	}
	return nil
}

func uncompressSymlink(linkName string, path string) error {
	return os.Symlink(linkName, path)
}
