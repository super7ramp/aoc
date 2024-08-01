package main

import (
	_ "embed"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type State rune

const (
	Operational State = '.'
	Damaged     State = '#'
	Unknown     State = '?'
)

type DamagedGroup struct {
	indexStart int
	indexEnd   int
}

func (group *DamagedGroup) size() int {
	return group.indexEnd - group.indexStart
}

type ConditionRecord struct {
	states            []State
	damagedGroupSizes []int
}

func ConditionRecordFrom(input string) ConditionRecord {
	fields := strings.Fields(input)
	states := make([]State, len(fields[0]))
	for i, char := range fields[0] {
		states[i] = State(char)
	}
	var damagedGroupSizes []int
	for _, damagedGroupSizeStr := range strings.Split(fields[1], ",") {
		damagedGroupSize, _ := strconv.Atoi(damagedGroupSizeStr)
		damagedGroupSizes = append(damagedGroupSizes, damagedGroupSize)
	}
	return ConditionRecord{states, damagedGroupSizes}
}

func (record *ConditionRecord) TargetDamagedCount() int {
	total := 0
	for _, damagedGroupSize := range record.damagedGroupSizes {
		total += damagedGroupSize
	}
	return total
}

func (record *ConditionRecord) DamagedCount() int {
	return record.count(Damaged)
}

func (record *ConditionRecord) DamagedToFindCount() int {
	return record.TargetDamagedCount() - record.DamagedCount()
}

func (record *ConditionRecord) count(state State) int {
	count := 0
	for _, s := range record.states {
		if s == state {
			count++
		}
	}
	return count
}

func (record *ConditionRecord) UnknownIndices() []int {
	indices := make([]int, 0)
	for i, state := range record.states {
		if state == Unknown {
			indices = append(indices, i)
		}
	}
	return indices
}

func (record *ConditionRecord) IsValid() bool {
	damagedGroups := record.DamagedGroups()
	if len(damagedGroups) != len(record.damagedGroupSizes) {
		return false
	}
	for i, damagedGroup := range damagedGroups {
		if damagedGroup.size() != record.damagedGroupSizes[i] {
			return false
		}
	}
	return true
}

func (record *ConditionRecord) DamagedGroups() []DamagedGroup {
	damagedGroups := make([]DamagedGroup, 0)
	currentGroup := DamagedGroup{indexStart: -1}
	for index, state := range record.states {
		if state == Damaged && currentGroup.indexStart == -1 {
			currentGroup.indexStart = index
		} else if state != Damaged && currentGroup.indexStart != -1 {
			currentGroup.indexEnd = index
			damagedGroups = append(damagedGroups, currentGroup)
			currentGroup = DamagedGroup{indexStart: -1}
		}
	}
	if currentGroup.indexStart != -1 {
		currentGroup.indexEnd = len(record.states)
		damagedGroups = append(damagedGroups, currentGroup)
	}
	return damagedGroups
}

type ConditionRecords []ConditionRecord

func ConditionRecordsFrom(input string) ConditionRecords {
	lines := strings.Split(input, "\n")
	conditionRecords := make(ConditionRecords, len(lines))
	for i, line := range lines {
		conditionRecords[i] = ConditionRecordFrom(line)
	}
	return conditionRecords
}

type Filler struct{}

func NewFiller() *Filler {
	return &Filler{}
}

func (filler *Filler) Fill(record *ConditionRecord) [][]State {
	candidates := filler.generateCandidates(record)
	candidates = slices.DeleteFunc(candidates, func(candidate []State) bool {
		candidateRecord := ConditionRecord{candidate, record.damagedGroupSizes}
		return !candidateRecord.IsValid()
	})
	return candidates
}

func (filler *Filler) generateCandidates(record *ConditionRecord) [][]State {
	unknownIndices := record.UnknownIndices()
	if len(unknownIndices) == 0 {
		return [][]State{record.states}
	}
	unknownIndicesFilledWithDamaged := filler.combinations(record.DamagedToFindCount(), unknownIndices)
	candidates := make([][]State, len(unknownIndicesFilledWithDamaged))
	for i, unknownIndexFilledWithDamaged := range unknownIndicesFilledWithDamaged {
		candidate := make([]State, len(record.states))
		copy(candidate, record.states)
		for _, unknownIndex := range unknownIndices {
			if slices.Contains(unknownIndexFilledWithDamaged, unknownIndex) {
				candidate[unknownIndex] = Damaged
			} else {
				candidate[unknownIndex] = Operational
			}
		}
		candidates[i] = candidate
	}
	return candidates
}

// See https://www.baeldung.com/cs/generate-k-combinations, lexicographic generation.
// Watch out, c indices start at 0 here.
func (filler *Filler) combinations(k int, elements []int) [][]int {
	result := make([][]int, 0)

	// Initialization
	c := make([]int, k+2)
	for i := 0; i < k; i++ {
		c[i] = i
	}
	c[k] = len(elements)
	c[k+1] = 0
	for {
		// Visit combination
		combination := make([]int, k)
		for i, ci := range c[0:k] {
			combination[i] = elements[ci]
		}
		result = append(result, combination)
		// Find j
		j := 0
		for c[j]+1 == c[j+1] {
			c[j] = j
			j++
		}
		// Done?
		if j < k {
			// Increment
			c[j] = c[j] + 1
		} else {
			// End
			break
		}
	}

	return result
}

//go:embed input.txt
var input string

func main() {
	conditionsRecords := ConditionRecordsFrom(input)
	filler := NewFiller()
	sum := 0
	for _, record := range conditionsRecords {
		solutions := filler.Fill(&record)
		fmt.Printf("%v solutions for %q: %q\n", len(solutions), record.states, solutions)
		sum += len(solutions)
	}
	fmt.Println("Total solution count: ", sum)
}
