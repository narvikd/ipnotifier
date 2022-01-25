package fileutils

import (
	"bufio"
	"fmt"
	"ipnotifier/pkg/errorsutils"
	"os"
	"path/filepath"
)

// Read receives a file path and returns the contents of the file as a slice.
func Read(path string) ([]string, error) {
	var contents []string
	fromPath := fmt.Sprintf("from path: \"%s\"", path)

	file, errOpenLogFile := os.OpenFile(filepath.Clean(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if errOpenLogFile != nil {
		return nil, errorsutils.Wrap(errOpenLogFile, "error opening file "+fromPath)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		contents = append(contents, scanner.Text())
	}

	errScan := scanner.Err()
	if errScan != nil {
		return nil, errorsutils.Wrap(errScan, "couldn't read file contents "+fromPath)
	}

	return contents, file.Close() // https://www.joeshaw.org/dont-defer-close-on-writable-files/
}

// Write receives a string and writes it to a file overriding its contents.
func Write(text string, path string) error {
	errWrite := os.WriteFile(filepath.Clean(path), []byte(text), 0600)
	if errWrite != nil {
		return errorsutils.Wrap(errWrite, "couldn't write text to file")
	}
	return nil
}
