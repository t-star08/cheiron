package ioFile

import (
	"bufio"
	"os"
)

func Puts(path string, source []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range source {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	
	writer.Flush()

	return nil
}

