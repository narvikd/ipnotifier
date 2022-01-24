package fileio

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func ReadIP() (string, error) {
	var fileContents []string
	path := filepath.Clean("ip.txt")
	file, errOpenLogFile := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if errOpenLogFile != nil {
		return "", errors.Wrap(errOpenLogFile, "error opening log file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
		if len(fileContents) >= 1 {
			break
		}
	}
	errScan := scanner.Err()
	if errScan != nil {
		return "", errors.Wrap(errScan, "couldn't read file contents")
	}

	if len(fileContents) == 0 {
		return "", nil
	}

	return fileContents[0], nil
}

func WriteIP(ip string) error {
	path := filepath.Clean("ip.txt")
	errWrite := os.WriteFile(path, []byte(ip), 0600)
	if errWrite != nil {
		return errors.Wrap(errWrite, "couldn't write ip to file")
	}
	return nil
}
