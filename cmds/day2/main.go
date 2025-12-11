package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	c := config.Parse()
	scanner, cleanUp, err := input.GetScanner(c, "day2")
	if err != nil {
		panic(err)
	}
	defer cleanUp()

	scanner.Scan()
	ranges := strings.Split(scanner.Text(), ",")

	values := make(chan int64, 100)
	var wg sync.WaitGroup

	for _, rangeInput := range ranges {
		wg.Go(func() {
			findBadIds(rangeInput, values)
		})
	}

	go func() {
		wg.Wait()
		close(values)
	}()

	var sum int64
	for value := range values {
		sum += value
	}

	fmt.Println(sum)
}

func findBadIds(rangeInput string, values chan<- int64) {
	rangeStart, rangeEnd := parseRange(rangeInput)

	for i := rangeStart; i <= rangeEnd; i++ {
		if matchesPart2(i) {
			values <- i
		}
	}
}

func matchesPart2(num int64) bool {
	if num < 10 {
		return false
	}

	if matchesPart1(num) {
		return true
	}

	str := strconv.Itoa(int(num))

	// all the same, like 111, 2222, 44444
	if isIdenticalParts([]byte(str)) {
		return true
	}

	length := len(str)
	for i := 2; i < length; i++ {
		chunks := chunkStr(str, i)

		if isIdenticalParts(chunks) {
			return true
		}
	}

	return false
}

func chunkStr(input string, size int) []string {
	var chunks []string
	builder := new(strings.Builder)
	for _, c := range input {
		builder.WriteRune(c)
		if builder.Len() == size {
			chunks = append(chunks, builder.String())
			builder.Reset()
		}
	}

	if builder.Len() != 0 {
		chunks = append(chunks, builder.String())
	}

	return chunks
}

func matchesPart1(num int64) bool {
	magnitude := orderOfMagnitude(num)
	if magnitude%2 == 1 {
		return false
	}

	half := magnitude / 2
	denominator := int64(math.Pow10(int(half)))
	denominator += 1

	if num%denominator == 0 {
		return true
	}

	return false
}

func isIdenticalParts[T comparable](items []T) bool {
	var first T
	allSame := true
	for i, c := range items {
		if i == 0 {
			first = c

			continue
		}

		if c != first {
			allSame = false
			break
		}
	}

	return allSame
}

var magnitudeCaps = []int64{
	10,
	100,
	1_000,
	10_000,
	100_000,
	1_000_000,
	10_000_000,
	100_000_000,
	1_000_000_000,
	10_000_000_000,
	100_000_000_000,
	1_000_000_000_000,
}

func orderOfMagnitude(i int64) int64 {
	for m, c := range magnitudeCaps {
		if i < c {
			return int64(m) + 1
		}
	}

	fmt.Fprintf(os.Stderr, "The number %d was too large", i)

	return -1
}

func parseRange(rangeInput string) (int64, int64) {
	rangeParts := strings.Split(rangeInput, "-")
	rangeStart, err := strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse %q as integer", rangeParts[0])
		os.Exit(1)
	}
	rangeEnd, err := strconv.ParseInt(rangeParts[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot parse %q as integer", rangeParts[1])
		os.Exit(2)
	}

	return rangeStart, rangeEnd
}
