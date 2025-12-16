package main

import (
	"fmt"
	"strconv"
	"strings"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	scanner, cleanUp, err := input.GetScanner(configuration, "day6")
	if err != nil {
		panic(err)
	}
	defer cleanUp()

	var lines [][]rune
	for scanner.Scan() {
		lines = append(lines, []rune(scanner.Text()))
	}

	var (
		lineCount = len(lines) - 1
		column    = len(lines[0]) - 1
		answer    int
	)
	for column > 0 {
		var problem = NewProblem()

		for {
			var number int

			for i := range lineCount {
				if lines[i][column] == ' ' {
					continue
				}

				digit := lines[i][column] - '0'

				number = (number * 10) + int(digit)
			}

			problem.AddNumber(number)

			if lines[lineCount][column] != ' ' {
				problem.Operation = ParseOperation(lines[lineCount][column])

				answer += problem.Execute()

				// skip seperate column too
				column -= 2

				break
			}

			column--
		}
	}

	fmt.Println(answer)
}

func Map[Slice ~[]E, E any, U any](slice Slice, mapper func(E) U) []U {
	var newSlice []U

	for _, item := range slice {
		newSlice = append(newSlice, mapper(item))
	}

	return newSlice
}

type Operation int

const (
	OperationAdd Operation = iota
	OperationMultiply
)

func ParseOperation(char rune) Operation {
	switch char {
	case '+':
		return OperationAdd

	case '*':
		return OperationMultiply

	default:
		panic(fmt.Errorf("Unknown operation string: %c", char))
	}
}

func (o Operation) Execute(a, b int) int {
	switch o {
	case OperationAdd:
		return a + b

	case OperationMultiply:
		return a * b
	}

	return 0
}

func (o Operation) String() string {
	switch o {
	case OperationAdd:
		return "+"

	case OperationMultiply:
		return "*"

	default:
		return "~"
	}
}

type Problem struct {
	Numbers   []int
	Operation Operation
}

func NewProblem() *Problem {
	return &Problem{
		Numbers: make([]int, 0),
	}
}

func (p *Problem) AddNumber(num int) {
	p.Numbers = append(p.Numbers, num)
}

func (p Problem) Execute() int {
	if len(p.Numbers) == 0 {
		return 0
	}

	result := p.Numbers[0]
	for _, n := range p.Numbers[1:] {
		result = p.Operation.Execute(result, n)
	}

	return result
}

func (p Problem) String() string {
	builder := new(strings.Builder)

	for i, n := range p.Numbers {
		builder.WriteString(strconv.Itoa(n))
		if i < len(p.Numbers)-1 {
			builder.WriteRune(' ')
			builder.WriteString(p.Operation.String())
			builder.WriteRune(' ')
		}
	}

	builder.WriteString(" = ")
	builder.WriteString(strconv.Itoa(p.Execute()))

	return builder.String()
}
