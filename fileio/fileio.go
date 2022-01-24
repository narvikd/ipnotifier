package fileio

import (
	"bufio"
	"errors"
	"ipnotifier/iputils"
	"ipnotifier/pkg/errorsutils"
	"os"
	"path/filepath"
)

func ReadIP(path string) (string, error) {
	fileContents := ""
	file, errOpenLogFile := os.OpenFile(filepath.Clean(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if errOpenLogFile != nil {
		return "", errorsutils.Wrap(errOpenLogFile, "error opening log file")
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContents = scanner.Text()
		break
	}

	errScan := scanner.Err()
	if errScan != nil {
		return "", errorsutils.Wrap(errScan, "couldn't read file contents")
	}
	if fileContents != "" && !iputils.IsIPValid(fileContents) {
		return "", errors.New("ip file has bogus content inside")
	}

	return fileContents, file.Close() // https://www.joeshaw.org/dont-defer-close-on-writable-files/
}

func WriteIP(ip string, path string) error {
	errWrite := os.WriteFile(filepath.Clean(path), []byte(ip), 0600)
	if errWrite != nil {
		return errorsutils.Wrap(errWrite, "couldn't write ip to file")
	}
	return nil
}
