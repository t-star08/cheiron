package ioFile

import (
	"bufio"
	"os"
)

func Gets(path string) ([]string, error) {
	doc := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		return doc, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		doc = append(doc, line)
	}

	return doc, nil
}

