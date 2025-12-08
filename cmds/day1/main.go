package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	inputFileName := os.Args[1]
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open the file %q\n", inputFileName)
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	dial := newDial()
	var password int64

	fmt.Printf("The dial starts by pointing at %d.\n", dial.value)
	for scanner.Scan() {
		line := scanner.Text()
		left := line[0] == 'L'

		countStr := line[1:]
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to convert the number %q\n", countStr)
			panic(err)
		}

		var sawZero int64
		if left {
			sawZero = dial.rotateLeft(count)
		} else {
			sawZero = dial.rotateRight(count)
		}

		fmt.Printf("The dial is rotated %s to point at %d", line, dial.value)
		if sawZero == 0 {
			fmt.Print(".\n")
		} else {
			fmt.Printf("; during this rotation it points at 0 %d time(s).\n", sawZero)
		}

		password += sawZero

		if dial.value == 0 {
			password++
		}
	}

	fmt.Println(password)
}

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}

	return i
}

type Dial struct {
	value int64
}

func newDial() *Dial {
	return &Dial{
		value: 50,
	}
}

func (d *Dial) rotateLeft(count int64) int64 {
	if d.value == 0 && count != 0 {
		d.value = 100
	}

	return d.rotate(-count)
}

func (d *Dial) rotateRight(count int64) int64 {
	return d.rotate(count)
}

func (d *Dial) rotate(count int64) int64 {
	sawZero := abs(count / 100)
	count = count % 100
	newValue := d.value + count
	wrapped := newValue < 0 || newValue > 99
	if newValue < 0 {
		newValue += 100
	} else if newValue > 99 {
		newValue -= 100
	}
	d.value = newValue
	if newValue != 0 && wrapped {
		sawZero += 1
	}

	if sawZero < 0 {
		return 0
	}

	return sawZero
}
