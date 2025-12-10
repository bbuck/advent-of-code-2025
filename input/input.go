package input

import (
	"bufio"
	"fmt"
	"os"

	"bbuck.dev/aoc2025/config"
)

// GetScanner determines the correct file to use and loads it.
func GetScanner(configuration config.Config, day string) (*bufio.Scanner, func(), error) {
	fileName := "sample"
	if configuration.Solve {
		fileName = "problem"
	}

	filePath := fmt.Sprintf("input/%s/%s.in", day, fileName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, func() {}, err
	}

	cleanUp := func() {
		file.Close()
	}

	scanner := bufio.NewScanner(file)

	return scanner, cleanUp, nil
}
