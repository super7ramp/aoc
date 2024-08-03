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

func (record *ConditionRecord) Unfold() ConditionRecord {
	unfoldedStates := make([]State, 0, 5*len(record.states)+4)
	unfoldedDamagedGroupSizes := make([]int, 0, 5*len(record.damagedGroupSizes))
	for repeat := 0; repeat < 5; repeat++ {
		for _, state := range record.states {
			unfoldedStates = append(unfoldedStates, state)
		}
		if repeat < 4 {
			unfoldedStates = append(unfoldedStates, Unknown)
		}
		for _, damagedGroupSize := range record.damagedGroupSizes {
			unfoldedDamagedGroupSizes = append(unfoldedDamagedGroupSizes, damagedGroupSize)
		}
	}
	return ConditionRecord{unfoldedStates, unfoldedDamagedGroupSizes}
}

func (record *ConditionRecord) String() string {
	return fmt.Sprintf("%v %v", string(record.states), record.damagedGroupSizes)
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

func (records ConditionRecords) Unfold() ConditionRecords {
	unfoldedConditionRecords := make(ConditionRecords, len(records))
	for i, record := range records {
		unfoldedConditionRecords[i] = record.Unfold()
	}
	return unfoldedConditionRecords
}

type NaiveFiller struct{}

func (filler *NaiveFiller) Fill(record *ConditionRecord) [][]State {
	unknownIndices := record.UnknownIndices()
	if len(unknownIndices) == 0 {
		return [][]State{record.states}
	}
	solutions := make([][]State, 0)
	filler.combinations(record.DamagedToFindCount(), unknownIndices, func(combination []int) {
		candidate := slices.Clone(record.states)
		for _, unknownIndex := range unknownIndices {
			if slices.Contains(combination, unknownIndex) {
				candidate[unknownIndex] = Damaged
			} else {
				candidate[unknownIndex] = Operational
			}
		}
		candidateRecord := ConditionRecord{candidate, record.damagedGroupSizes}
		if candidateRecord.IsValid() {
			solutions = append(solutions, candidate)
		}
	})
	return solutions
}

// See https://www.baeldung.com/cs/generate-k-combinations, lexicographic generation.
// Watch out, c indices start at 0 here.
func (filler *NaiveFiller) combinations(k int, elements []int, visitCombination func([]int)) {

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
		visitCombination(combination)

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
}

type LessNaiveFillCounter struct {
	cache map[string]int64
}

func NewLessNaiveFillCounter() *LessNaiveFillCounter {
	return &LessNaiveFillCounter{cache: make(map[string]int64)}
}

func (filler *LessNaiveFillCounter) CountFills(record *ConditionRecord) int64 {

	if cachedCount, ok := filler.cache[record.String()]; ok {
		return cachedCount
	}

	groupSizes := record.damagedGroupSizes
	states := record.states

	if len(groupSizes) == 0 {
		candidate := slices.Clone(states)
		for i, s := range candidate {
			if s == Damaged {
				return 0
			}
			candidate[i] = Operational
		}
		return 1
	}

	if filler.requiredSpace(groupSizes) > len(states) {
		return 0
	}

	candidate := make([]State, len(states))
	groupSize := groupSizes[0]
	count := int64(0)

outer:
	for groupStartIndex := 0; groupStartIndex < len(states)-groupSize+1; groupStartIndex++ {
		copy(candidate, states)
		for beforeStartIndex := 0; beforeStartIndex < groupStartIndex; beforeStartIndex++ {
			if candidate[beforeStartIndex] == Damaged {
				continue outer
			}
			candidate[beforeStartIndex] = Operational
		}
		for indexInGroup := 0; indexInGroup < groupSize; indexInGroup++ {
			if candidate[groupStartIndex+indexInGroup] == Operational {
				continue outer
			}
			candidate[groupStartIndex+indexInGroup] = Damaged
		}
		groupSeparatorIndex := groupStartIndex + groupSize
		if len(candidate) > groupSeparatorIndex {
			if candidate[groupSeparatorIndex] == Damaged {
				continue
			}
			candidate[groupSeparatorIndex] = Operational
		}
		if len(candidate) > groupSeparatorIndex+1 {
			subRecord := &ConditionRecord{candidate[groupSeparatorIndex+1:], groupSizes[1:]}
			count += filler.CountFills(subRecord)
		} else if len(groupSizes) == 1 {
			count++
		}
	}

	filler.cache[record.String()] = count
	return count
}

func (filler *LessNaiveFillCounter) requiredSpace(groupSizes []int) int {
	count := 0
	for _, groupSize := range groupSizes {
		count += groupSize
	}
	count += len(groupSizes) - 1
	return count
}

//go:embed input.txt
var input string

func main() {
	conditionsRecords := ConditionRecordsFrom(input)
	filler := NaiveFiller{}
	sum := int64(0)
	for _, record := range conditionsRecords {
		solutions := filler.Fill(&record)
		fmt.Printf("%v solutions for %q\n", len(solutions), record.states)
		sum += int64(len(solutions))
	}
	fmt.Println("Total solution count (part 1, naive filler): ", sum)

	fillCounter := NewLessNaiveFillCounter()
	sum = 0
	for _, record := range conditionsRecords {
		count := fillCounter.CountFills(&record)
		fmt.Printf("%v solutions for %q\n", count, record.states)
		sum += count
	}
	fmt.Println("Total solution count (part 1, less naive fill counter): ", sum)

	sum = 0
	conditionsRecords = conditionsRecords.Unfold()
	for _, record := range conditionsRecords {
		count := fillCounter.CountFills(&record)
		fmt.Printf("%v solutions for %q\n", count, record.states)
		sum += count
	}
	fmt.Println("Total solution count (part 2, less naive fill counter): ", sum)

}
