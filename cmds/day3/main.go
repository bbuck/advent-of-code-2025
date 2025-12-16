package main

import (
	"fmt"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	scanner, cleanUp, err := input.GetScanner(configuration, "day3")
	if err != nil {
		panic(err)
	}
	defer cleanUp()

	var sum int
	for scanner.Scan() {
		bank := convertToBank(scanner.Text())
		joltage := getJoltageLarge(bank)
		sum += joltage
	}

	fmt.Println(sum)
}

func convertToBank(input string) []int {
	var bank []int
	for _, c := range input {
		bank = append(bank, int(c-'0'))
	}

	return bank
}

func getJoltageLarge(bank []int) int {
	var (
		joltage int
		bankLen = len(bank)
		stop    = -1
	)

	for i := 12; i > 0; i-- {
		var (
			start   = bankLen - i
			foundAt = -1
			digit   = -1
		)

		for j := start; j > stop; j-- {
			if bank[j] >= digit {
				digit = bank[j]
				foundAt = j
			}
		}

		stop = foundAt
		joltage = (joltage * 10) + digit
	}

	return joltage
}

// func getJoltageSimple(bank []int) int {
// 	maxTen := -1
// 	maxOne := -1
//
// 	bankLen := len(bank)
// 	for i := 0; i < bankLen; i++ {
// 		if bank[i] < maxTen {
// 			continue
// 		}
//
// 		oldMaxOne := maxOne
// 		if bank[i] > maxTen {
// 			maxOne = -1
// 		}
//
// 		wasSet := false
// 		for j := i + 1; j < bankLen; j++ {
// 			if bank[j] <= maxOne {
// 				continue
// 			}
//
// 			maxOne = bank[j]
// 			wasSet = true
// 		}
//
// 		if wasSet {
// 			maxTen = bank[i]
// 		} else {
// 			maxOne = oldMaxOne
// 		}
// 	}
//
// 	return (maxTen * 10) + maxOne
// }
