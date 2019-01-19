package shred

import (
	"crypto/rand"
	"os"
	"path/filepath"
)

// Path shreds all files in the location of the given path
// recursively. If remove is set to true files will be deleted
// after shredding. When a file is shredded its content
// is NOT recoverable so !!USE WITH CAUTION!!
func Path(path string, remove bool) {
	stats, err := os.Stat(path)
	if err != nil {
		return
	} else if stats.IsDir() {
		Dir(path, remove)
	} else {
		File(path, remove)
	}

}

// File overwrites a given File in the location of path
// with random Data followed by another overwrite with zeros.
func File(path string, remove bool) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return
	}
	size := fileinfo.Size()
	writeRandom(path, size)
	writeZeros(path, size)
	if remove {
		os.Remove(path)
	}

}

// Dir overwrites every File in root and all its subdirectories
func Dir(root string, remove bool) {
	filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			stats, _ := os.Stat(path)

			if stats.IsDir() == false {
				File(path, remove)
			}
			return nil
		})
}

// writeRandom overwrites a File with random stuff
func writeRandom(path string, size int64) {
	for i := 0; i < 3; i++ {
		file, err := os.OpenFile(path, os.O_RDWR, 0)
		defer file.Close()
		if err != nil {
			return
		}
		offset, err := file.Seek(0, 0)
		if err != nil {
			return
		}
		buff := make([]byte, size)
		rand.Read(buff)
		file.WriteAt(buff, offset)
		file.Sync()
		file.Close()
	}
}

// writeZeros overwrites a File with zeros
func writeZeros(path string, size int64) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	defer file.Close()
	if err != nil {
		return
	}

	offset, err := file.Seek(0, 0)
	if err != nil {
		return
	}
	buff := make([]byte, size)
	file.WriteAt(buff, offset)
	file.Sync()
	file.Close()
}
