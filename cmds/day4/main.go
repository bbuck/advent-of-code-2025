package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	scanner, cleanUp, err := input.GetScanner(configuration, "day4")
	if err != nil {
		panic(err)
	}
	defer cleanUp()

	var ranges []Range
	for scanner.Scan() {
		// stop scanning ranges
		if scanner.Text() == "" {
			break
		}

		newRange, err := ParseRange(scanner.Text())
		if err != nil {
			panic(fmt.Errorf("failed to parse range %q: %w\n", scanner.Text(), err))
		}

		ranges = append(ranges, newRange)
	}

	targetIndex := 0
	for {
		if len(ranges) < 2 {
			break
		}

		if targetIndex >= len(ranges) {
			break
		}

		var (
			merges    int
			newRanges []Range
			target    = ranges[targetIndex]
		)

		for i, r := range ranges {
			if i == targetIndex {
				continue
			}

			if target.Overlaps(r) {
				target = target.Merge(r)
				merges++
			} else {
				newRanges = append(newRanges, r)
			}
		}

		ranges = append(newRanges, target)

		if merges == 0 {
			targetIndex++
		} else {
			targetIndex = 0
		}
	}

	slices.SortFunc(ranges, func(a Range, b Range) int {
		if a.Minimum < b.Minimum {
			return -1
		}

		if a.Minimum > b.Minimum {
			return 1
		}

		return 0
	})

	var freshCount int
	for _, r := range ranges {
		freshCount += r.Count()
	}

	fmt.Println(freshCount)
}

type Range struct {
	Minimum int
	Maximum int
}

func ParseRange(input string) (Range, error) {
	startStr, endStr, found := strings.Cut(input, "-")
	var newRange Range

	if !found {
		return newRange, errors.New("not a valid range string")
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		return newRange, fmt.Errorf("range start is not a valid number: %w", err)
	}

	end, err := strconv.Atoi(endStr)
	if err != nil {
		return newRange, fmt.Errorf("range end is not a valid number: %w", err)
	}

	newRange.Minimum = start
	newRange.Maximum = end

	return newRange, nil
}

func (r Range) String() string {
	return fmt.Sprintf("%d - %d", r.Minimum, r.Maximum)
}

func (r Range) Contains(value int) bool {
	return value >= r.Minimum && value <= r.Maximum
}

func (r Range) Overlaps(other Range) bool {
	return r.Contains(other.Minimum) || r.Contains(other.Maximum)
}

func (r Range) Merge(other Range) Range {
	return Range{
		Minimum: min(r.Minimum, other.Minimum),
		Maximum: max(r.Maximum, other.Maximum),
	}
}

func (r Range) Count() int {
	return r.Maximum - r.Minimum + 1
}

func min(values ...int) int {
	minimum := values[0]

	for i := 1; i < len(values); i++ {
		if values[i] < minimum {
			minimum = values[i]
		}
	}

	return minimum
}

func max(values ...int) int {
	maximum := values[0]

	for i := 1; i < len(values); i++ {
		if values[i] > maximum {
			maximum = values[i]
		}
	}

	return maximum
}
