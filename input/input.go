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

	filePath := fmt.Sprintf("%s/%s.in", day, fileName)

	file, err := os.OpenInRoot("inputs", filePath)
	if err != nil {
		return nil, func() {}, err
	}

	cleanUp := func() {
		file.Close()
	}

	scanner := bufio.NewScanner(file)

	return scanner, cleanUp, nil
}

// ReadInput reads all the input from the problem file into a string slice.
func ReadInput(configuration config.Config, day string) ([]string, error) {
	scanner, cleanUp, err := GetScanner(configuration, day)
	if err != nil {
		return nil, err
	}
	defer cleanUp()

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return lines, nil
}
