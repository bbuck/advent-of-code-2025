package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"sync"

	"bbuck.dev/aoc2025/config"
	"bbuck.dev/aoc2025/input"
)

func main() {
	configuration := config.Parse()
	lines, err := input.ReadInput(configuration, "day10")
	if err != nil {
		panic(err)
	}

	var machines []Machine
	for _, line := range lines {
		machine := ParseMachine(line)

		machines = append(machines, machine)
	}

	SolvePart2(machines)
}

func SolvePart1(machines []Machine) {
	var (
		solutionLengths = make(chan int, len(machines))
		wg              sync.WaitGroup
	)

	for _, machine := range machines {
		wg.Go(func() {
			solution := FindSolution(machine)

			solutionLengths <- len(solution)
		})
	}

	go func() {
		wg.Wait()
		close(solutionLengths)
	}()

	var result int
	for value := range solutionLengths {
		result += value
	}

	fmt.Println(result)
}

func SolvePart2(machines []Machine) {
	var sum int
	for i, machine := range machines {
		matrix := machine.Matrix()
		matrix.RowEchelonForm()

		presses, solved := matrix.Solve()

		if !solved {
			fmt.Println(i)
			fmt.Println(machine.Matrix())
			fmt.Println(matrix)

			return
		}

		sum += sumInts(presses)
	}

	fmt.Println(sum)
}

type Row []float64

func (r Row) String() string {
	builder := new(strings.Builder)

	for i, col := range r {
		if i > 0 {
			builder.WriteRune(' ')
		}

		if i == len(r)-1 {
			builder.WriteString("| ")
		}

		fmt.Fprintf(builder, "%f", col)
	}

	return builder.String()
}

func (r Row) Subtract(otherRow Row) Row {
	var newRow Row
	for i, value := range r {
		newRow = append(newRow, value-otherRow[i])
	}

	return newRow
}

func (r Row) Multiply(scalar float64) Row {
	var newRow Row
	for _, value := range r {
		newRow = append(newRow, value*scalar)
	}
	return newRow
}

func (r Row) Divide(scalar float64) Row {
	var newRow Row
	for _, value := range r {
		newRow = append(newRow, value/scalar)
	}

	return newRow
}

func (r Row) Sum() float64 {
	return r[len(r)-1]
}

func (r Row) Check(variables []int) bool {
	sum := r.Sum()
	var total float64
	for i := range len(r) - 1 {
		if variables[i] < 0 || variables[i] > topOut {
			return false
		}

		value := r[i] * float64(variables[i])

		total += value
	}

	return equalApprox(total, sum)
}

type Matrix []Row

func (m Matrix) String() string {
	builder := new(strings.Builder)
	for _, row := range m {
		builder.WriteString("[ ")
		builder.WriteString(row.String())
		builder.WriteString("]\n")
	}

	return builder.String()
}

func equalApprox(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}

func (m Matrix) RowEchelonForm() {
	targetRow := 0

	for column := range len(m[0]) {
		if targetRow >= len(m) {
			break
		}

		found := -1
		for i := targetRow; i < len(m); i++ {
			if !equalApprox(m[i][column], 0) {
				found = i

				break
			}
		}

		if found < 0 {
			continue
		}

		if !equalApprox(m[found][column], 1) {
			m[found] = m[found].Divide(m[found][column])
		}

		m[targetRow], m[found] = m[found], m[targetRow]

		for j := targetRow + 1; j < len(m); j++ {
			if !equalApprox(m[j][column], 0) {
				if equalApprox(m[j][column], 1) {
					m[j] = m[j].Subtract(m[targetRow])
				} else {
					m[j] = m[j].Subtract(m[targetRow].Multiply(m[j][column]))
				}
			}
		}

		targetRow++
		column++
	}
}

func (m Matrix) Solve() ([]int, bool) {
	variables := make([]int, len(m[0])-1)

	return m.solveUp(len(m)-1, len(variables), variables)
}

type Solution struct {
	variables []int
	presses   int
}

const topOut = 266

func (m Matrix) solveUp(row, solvedVariables int, variables []int) ([]int, bool) {
	if row < 0 {
		return variables, true
	}

	var solutions []Solution

	pivotCol := slices.IndexFunc(m[row], func(val float64) bool {
		return equalApprox(val, 1)
	})

	newVariables := make([]int, len(variables))
	copy(newVariables, variables)
	if pivotCol >= 0 && pivotCol < len(m[row])-1 && pivotCol < solvedVariables {
		var freeVariables []int
		for i := pivotCol + 1; i < solvedVariables; i++ {
			freeVariables = append(freeVariables, i)
		}

		if len(freeVariables) == 0 {
			sum := m[row].Sum()
			for i := pivotCol + 1; i < len(newVariables); i++ {
				sum -= m[row][i] * float64(newVariables[i])
			}
			newVariables[pivotCol] = int(math.Round(sum))

			if m[row].Check(newVariables) {
				return m.solveUp(row-1, pivotCol, newVariables)
			}

			return variables, false
		} else {
			guessValues := make([]int, len(freeVariables))
		parentLoop:
			for {
				if guessValues[len(guessValues)-1] > topOut {
					break
				}

				for i, fv := range freeVariables {
					newVariables[fv] = guessValues[i]
				}
				sum := m[row].Sum()
				newVariables[pivotCol] = 0
				for i := pivotCol + 1; i < len(newVariables); i++ {
					sum -= m[row][i] * float64(newVariables[i])
				}
				newVariables[pivotCol] = int(math.Round(sum))

				if m[row].Check(newVariables) {
					if finalVars, solved := m.solveUp(row-1, pivotCol, newVariables); solved {
						solution := Solution{
							variables: finalVars,
							presses:   sumInts(finalVars),
						}
						solutions = append(solutions, solution)
					}
				}

				guessValues[0]++
				for i, v := range guessValues {
					if v > topOut {
						if i >= len(guessValues)-1 {
							break parentLoop
						}
						guessValues[i] = 0
						guessValues[i+1]++

						break
					}
				}
			}
		}
	}

	if pivotCol < 0 && m[row].Check(newVariables) {
		return m.solveUp(row-1, solvedVariables, newVariables)
	}

	if len(solutions) > 0 {
		fmt.Println("Found multiple possible solutions...")
		for _, solution := range solutions {
			fmt.Println(" --", solution.variables, solution.presses)
		}
		best := solutions[0]
		for i := 1; i < len(solutions); i++ {
			if solutions[i].presses < best.presses {
				best = solutions[i]
			}
		}
		fmt.Println("Selected", best.variables, best.presses)

		return best.variables, m[row].Check(best.variables)
	}

	return variables, false
}

func sumInts(ints []int) int {
	var sum int
	for _, v := range ints {
		sum += v
	}

	return sum
}

func FindSolution(m Machine) []int {
	root := NewStep(m, -1, nil)

	solved, _ := root.RunStep()
	for solved == nil {
		solved, _ = root.RunStep()
	}

	return solved.Solution()
}

type Step struct {
	Parent      *Step
	ButtonIndex int
	Machine     Machine
	initialized bool
	Children    []*Step
}

func NewStep(m Machine, buttonIndex int, parent *Step) *Step {
	return &Step{
		Parent:      parent,
		ButtonIndex: buttonIndex,
		Machine:     m,
	}
}

func (s *Step) Solution() []int {
	var (
		solution []int
		current  = s
	)

	for current != nil {
		if current.ButtonIndex >= 0 {
			solution = append(solution, current.ButtonIndex)
		}
		current = current.Parent
	}

	slices.Reverse(solution)

	return solution
}

func (s *Step) Solved() bool {
	return s.Machine.Solved()
}

func (s *Step) RunStep() (*Step, bool) {
	if !s.initialized {
		return s.spawnChildren()
	}

	var (
		newChildren []*Step
		solvedChild *Step
	)
	for _, child := range s.Children {
		var prune bool
		solvedChild, prune = child.RunStep()

		if !prune {
			newChildren = append(newChildren, child)
		}

		if solvedChild != nil {
			break
		}
	}

	s.Children = newChildren

	return solvedChild, len(s.Children) == 0
}

func (s *Step) ShouldContinue() bool {
	if s.initialized {
		return len(s.Children) > 0
	}

	return s.Machine.CanComplete()
}

func (s *Step) spawnChildren() (*Step, bool) {
	for i, button := range s.Machine.Buttons {
		step := NewStep(s.Machine.Press(button), i, s)

		if step.Solved() || step.ShouldContinue() {
			s.Children = append(s.Children, step)
		}

		if step.Solved() {
			return step, false
		}
	}

	s.initialized = true

	return nil, len(s.Children) == 0
}

type Machine struct {
	Indicators      int
	IndicatorTarget int
	Buttons         []int
	AddJoltage      bool
	Joltages        []int
	JoltageTarget   []int
}

func ParseMachine(input string) Machine {
	parts := strings.Split(input, " ")

	indicatorInput := parts[0]
	joltageInput := parts[len(parts)-1]
	buttonInputs := parts[1 : len(parts)-1]

	var indicator int
	for i, light := range indicatorInput[1 : len(indicatorInput)-1] {
		if light == '.' {
			continue
		}

		indicator = indicator | (1 << i)
	}

	var (
		joltageTarget []int
		joltageInputs = strings.Split(joltageInput[1:len(joltageInput)-1], ",")
	)
	for _, joltageInput := range joltageInputs {
		joltage, err := strconv.Atoi(joltageInput)
		if err != nil {
			panic(err)
		}

		joltageTarget = append(joltageTarget, joltage)
	}

	var buttons []int
	for _, buttonInput := range buttonInputs {
		flips := strings.Split(buttonInput[1:len(buttonInput)-1], ",")

		var button int
		for _, flipIndicator := range flips {
			value, err := strconv.Atoi(flipIndicator)
			if err != nil {
				panic(err)
			}

			button = button | (1 << value)
		}

		buttons = append(buttons, button)
	}

	return Machine{
		IndicatorTarget: indicator,
		JoltageTarget:   joltageTarget,
		Joltages:        make([]int, len(joltageTarget)),
		Buttons:         buttons,
	}
}

func (m Machine) CanComplete() bool {
	if m.AddJoltage {
		for i := range len(m.Joltages) {
			if m.Joltages[i] > m.JoltageTarget[i] {
				return false
			}
		}
	}

	return true
}

func (m Machine) Press(button int) Machine {
	var next Machine
	if m.AddJoltage {
		next = m.pressJoltage(button)
	} else {
		next = m.pressIndicator(button)
	}

	return next
}

func (m Machine) pressIndicator(button int) Machine {
	m.Indicators = m.Indicators ^ button

	return m
}

func (m Machine) pressJoltage(button int) Machine {
	var joltages []int
	for i, joltage := range m.Joltages {
		newJoltage := joltage

		value := 1 << i
		if button&value == value {
			newJoltage++
		}

		joltages = append(joltages, newJoltage)
	}

	m.Joltages = joltages

	return m
}

func (m Machine) PressAt(buttonIndex int) Machine {
	button := m.Buttons[buttonIndex]

	return m.Press(button)
}

func (m Machine) Solved() bool {
	if m.AddJoltage {
		return m.solvedJoltages()
	}

	return m.solvedIndicators()
}

func (m Machine) solvedIndicators() bool {
	return m.Indicators == m.IndicatorTarget
}

func (m Machine) solvedJoltages() bool {
	for i := range len(m.Joltages) {
		if m.Joltages[i] != m.JoltageTarget[i] {
			return false
		}
	}

	return true
}

func (m Machine) Matrix() Matrix {
	var matrix Matrix
	for _, target := range m.JoltageTarget {
		row := make(Row, len(m.Buttons)+1)
		matrix = append(matrix, row)

		row[len(row)-1] = float64(target)
	}

	for i, button := range m.Buttons {
		for j := range len(m.JoltageTarget) {
			value := 1 << j
			if button&value == value {
				matrix[j][i] = 1
			}
		}
	}

	return matrix
}

func (m Machine) String() string {
	builder := new(strings.Builder)

	indicatorFormat := fmt.Sprintf("[%%0%db]", len(m.Joltages))
	buttonFormat := fmt.Sprintf("(%%0%db)", len(m.Joltages))

	builder.WriteRune(' ')
	fmt.Fprintf(builder, indicatorFormat, m.Indicators)
	label := "OFF"
	if m.Solved() {
		label = "ON"
	}
	fmt.Fprintf(builder, " %s\n", label)

	for _, button := range m.Buttons {
		fmt.Fprintf(builder, buttonFormat, button)
		builder.WriteRune(' ')
	}

	builder.WriteRune('\n')

	builder.WriteRune('[')
	fmt.Fprintf(builder, indicatorFormat, m.IndicatorTarget)
	builder.WriteString("]\n")

	builder.WriteString(" {")
	for i, joltage := range m.Joltages {
		if i > 0 {
			builder.WriteRune(' ')
		}
		fmt.Fprintf(builder, "%3d", joltage)
	}
	builder.WriteString("}\n{{")
	for i, target := range m.JoltageTarget {
		if i > 0 {
			builder.WriteRune(' ')
		}
		fmt.Fprintf(builder, "%3d", target)
	}
	builder.WriteString("}}\n")

	return builder.String()
}
