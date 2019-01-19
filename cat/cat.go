package cat

import (
	"io/ioutil"
	"os"
)

// TODO: add support for extracting text from .docx .pptx .od* .pdf

//Cat is cat
func Cat(file string) string {
	return readTxt(file)
}

func readTxt(file string) string {
	_, err := os.Stat(file)
	if err != nil {
		return "[!] File does not exist!"
	}
	outbyte, err := ioutil.ReadFile(file)
	if err != nil {
		return "[!] Error while reading file"
	}
	return string(outbyte)
}
