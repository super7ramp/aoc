package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Operator byte

const (
	Addition       = Operator('+')
	Subtraction    = Operator('-')
	Multiplication = Operator('*')
	Division       = Operator('/')
	Concatenation  = Operator('|')
)

func (o *Operator) apply(a, b float64) float64 {
	switch *o {
	case Addition:
		return a + b
	case Subtraction:
		return a - b
	case Multiplication:
		return a * b
	case Division:
		return a / b
	case Concatenation:
		concat := fmt.Sprintf("%d%d", int(a), int(b))
		result, _ := strconv.ParseFloat(concat, 64)
		return result
	}
	panic("Invalid operator")
}

type Equation struct {
	result   int
	operands []int
}

func EquationFrom(value string) *Equation {
	fields := strings.Split(value, ":")
	result, _ := strconv.Atoi(fields[0])
	operands := make([]int, 0)
	for _, operandValue := range strings.Split(strings.TrimSpace(fields[1]), " ") {
		operand, _ := strconv.Atoi(operandValue)
		operands = append(operands, operand)
	}
	return &Equation{result, operands}
}

// FindOperators returns the operators that make the equation valid, or nil if no such operators are found.
func (e *Equation) FindOperators(allowedOperators ...Operator) []Operator {
	testedOperators := make([]Operator, len(e.operands)-1)
	combinationCount := pow(len(allowedOperators), len(testedOperators))
	for combination := range combinationCount {
		for i := range testedOperators {
			elephantOperatorIndex := combination / (pow(len(allowedOperators), i)) % len(allowedOperators)
			testedOperators[i] = allowedOperators[elephantOperatorIndex]
		}
		if e.evaluate(testedOperators) {
			//fmt.Println(e.DebugString(testedOperators))
			return testedOperators
		}
	}
	return nil
}

// evaluate returns true if the equation evaluates to the result using the given operators.
func (e *Equation) evaluate(operators []Operator) bool {
	result := float64(e.operands[0])
	for i, operand := range e.operands[1:] {
		operator := operators[i]
		result = operator.apply(result, float64(operand))
	}
	return math.Abs(result-float64(e.result)) < 1e-9
}

func (e *Equation) DebugString(operators []Operator) string {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(e.result))
	sb.WriteRune('=')
	sb.WriteString(strconv.Itoa(e.operands[0]))
	for i, operand := range e.operands[1:] {
		sb.WriteRune(rune(operators[i]))
		sb.WriteString(strconv.Itoa(operand))
	}
	return sb.String()
}

func pow(a, b int) int {
	result := 1
	for range b {
		result *= a
	}
	return result
}

type Equations []Equation

func EquationsFrom(value string) Equations {
	equations := make(Equations, 0)
	for _, line := range strings.Split(value, "\n") {
		equations = append(equations, *EquationFrom(line))
	}
	return equations
}

func (e *Equations) TotalCalibrationResult(allowedOperators ...Operator) int {
	totalCalibrationResult := 0
	for _, equation := range *e {
		if operators := equation.FindOperators(allowedOperators...); operators != nil {
			totalCalibrationResult += equation.result
		}
	}
	return totalCalibrationResult
}

//go:embed input.txt
var input string

func main() {
	equations := EquationsFrom(input)
	totalCalibrationResult := equations.TotalCalibrationResult(Addition, Multiplication)
	fmt.Println("(Part 1) Total calibration result:", totalCalibrationResult)

	totalCalibrationResult = equations.TotalCalibrationResult(Addition, Multiplication, Concatenation)
	fmt.Println("(Part 2) Total calibration result:", totalCalibrationResult)
}
