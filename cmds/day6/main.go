package main

import (
	"fmt"
	"io"
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

	var (
		problems []*Problem
	)
	for scanner.Scan() {
		reader := strings.NewReader(scanner.Text())
		numbers, ok := readNumbers(reader)

		if ok {
			if problems == nil {
				for range len(numbers) {
					problems = append(problems, NewProblem())
				}
			}

			for i, num := range numbers {
				problems[i].AddNumber(num)
			}

			continue
		}

		operators, _ := readOperators(reader)
		for i, op := range operators {
			problems[i].Operation = op
		}
	}

	var answer int
	for _, problem := range problems {
		answer += problem.Execute()
	}

	fmt.Println(answer)
}

func readNumbers(reader io.Reader) ([]int, bool) {
	var (
		numbers []int
		current int
		err     error
	)

	for err == nil {
		_, err = fmt.Fscanf(reader, "%d", &current)
		if err == nil {
			numbers = append(numbers, current)
		}
	}

	return numbers, len(numbers) > 0
}

func readOperators(reader io.Reader) ([]Operation, bool) {
	var (
		operations []Operation
		current    string
		err        error
	)

	for {
		_, err = fmt.Fscanf(reader, "%s", &current)
		if err != nil {
			break
		}

		switch current {
		case "*":
			operations = append(operations, OperationMultiply)
		case "+":
			operations = append(operations, OperationAdd)
		default:
			panic(fmt.Errorf("Failed to identify operation %q", current))
		}
	}

	return operations, len(operations) > 0
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
