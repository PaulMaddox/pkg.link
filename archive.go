package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/GeertJohan/go.rice"
)

func TarGz(root string, info os.FileInfo, box *rice.Box) (io.Reader, error) {

	buffer := bytes.NewBuffer(nil)

	g := gzip.NewWriter(buffer)
	t := tar.NewWriter(g)

	defer g.Close()
	defer t.Close()

	box.Walk(root, func(path string, info os.FileInfo, err error) error {

		// Don't archive empty directories
		if info.IsDir() {
			return nil
		}

		file, err := box.Open(path)
		if err != nil {
			return errors.New("failed to archive " + path)
		}

		filename := strings.Replace(path, root+string(os.PathSeparator), "", 1)

		t.WriteHeader(&tar.Header{
			Name:       filename,
			Size:       info.Size(),
			Mode:       0666,
			ModTime:    info.ModTime(),
			AccessTime: info.ModTime(),
			ChangeTime: info.ModTime(),
		})

		contents, err := ioutil.ReadAll(file)
		if err != nil {
			return errors.New("failed to read " + path)
		}

		t.Write(contents)

		return nil

	})

	t.Close()
	g.Close()

	return buffer, nil

}
