package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Almanac struct {
	seeds                 []int
	seedToSoil            Associations
	soilToFertilizer      Associations
	fertilizerToWater     Associations
	waterToLight          Associations
	lightToTemperature    Associations
	temperatureToHumidity Associations
	humidityToLocation    Associations
}

func AlmanachFrom(in string) *Almanac {
	sections := strings.Split(in, "\n\n")
	return &Almanac{
		seeds:                 seedsFrom(sections[0]),
		seedToSoil:            associationsFrom(sections[1]),
		soilToFertilizer:      associationsFrom(sections[2]),
		fertilizerToWater:     associationsFrom(sections[3]),
		waterToLight:          associationsFrom(sections[4]),
		lightToTemperature:    associationsFrom(sections[5]),
		temperatureToHumidity: associationsFrom(sections[6]),
		humidityToLocation:    associationsFrom(sections[7]),
	}
}

func seedsFrom(seedSection string) []int {
	seedFields := strings.Fields(seedSection)[1:]
	seeds := make([]int, len(seedFields))
	for i, seedField := range seedFields {
		seeds[i], _ = strconv.Atoi(seedField)
	}
	return seeds
}

func (a *Almanac) LocationForSeed(seed int) int {
	soil := a.seedToSoil.Destination(seed)
	fertilizer := a.soilToFertilizer.Destination(soil)
	water := a.fertilizerToWater.Destination(fertilizer)
	light := a.waterToLight.Destination(water)
	temperature := a.lightToTemperature.Destination(light)
	humidity := a.temperatureToHumidity.Destination(temperature)
	return a.humidityToLocation.Destination(humidity)
}

func (a *Almanac) SeedRanges() []Range {
	seedRanges := make([]Range, len(a.seeds)/2)
	for i := 0; i < len(a.seeds); i += 2 {
		seedRanges[i/2] = Range{a.seeds[i], a.seeds[i+1]}
	}
	return seedRanges
}

type Associations []Association

func associationsFrom(mapSection string) Associations {
	associationLines := strings.Split(mapSection, "\n")[1:]
	associations := make([]Association, len(associationLines))
	for i, associationLine := range associationLines {
		fields := strings.Fields(associationLine)
		targetRangerStart, _ := strconv.Atoi(fields[0])
		sourceRangeStart, _ := strconv.Atoi(fields[1])
		rangeLength, _ := strconv.Atoi(fields[2])
		associations[i] = Association{targetRangerStart, sourceRangeStart, rangeLength}
	}
	return associations
}

func (associations Associations) Destination(source int) int {
	for _, association := range associations {
		if dest := association.Destination(source); dest != -1 {
			return dest
		}
	}
	return source
}

type Association struct {
	targetRangerStart int
	sourceRangeStart  int
	rangeLength       int
}

func (a *Association) Destination(source int) int {
	if source < a.sourceRangeStart || source >= a.sourceRangeStart+a.rangeLength {
		return -1
	}
	return a.targetRangerStart + source - a.sourceRangeStart
}

type Range struct {
	start  int
	length int
}

func (r *Range) end() int {
	return r.start + r.length
}

//go:embed input.txt
var input string

func main() {
	almanac := AlmanachFrom(input)

	minLocation := math.MaxInt
	for _, seed := range almanac.seeds {
		if location := almanac.LocationForSeed(seed); location < minLocation {
			minLocation = location
		}
	}
	fmt.Println("Min location for seeds (part 1): ", minLocation)

	// the following takes 2 min on a M2, there must be a cleverer way 😅
	minLocation = math.MaxInt
	for _, seedRange := range almanac.SeedRanges() {
		for seed := seedRange.start; seed < seedRange.end(); seed++ {
			if location := almanac.LocationForSeed(seed); location < minLocation {
				minLocation = location
			}
		}
	}
	fmt.Println("Min location for seeds (part 2): ", minLocation)
}
