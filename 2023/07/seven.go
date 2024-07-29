package main

import (
	_ "embed"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	joker = iota
	two   = iota + 1
	three
	four
	five
	six
	seven
	height
	nine
	ten
	jack
	queen
	king
	ace
)

const (
	highCard = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

type Card int

func CardFrom(r rune) Card {
	switch r {
	case '2':
		return two
	case '3':
		return three
	case '4':
		return four
	case '5':
		return five
	case '6':
		return six
	case '7':
		return seven
	case '8':
		return height
	case '9':
		return nine
	case 'T':
		return ten
	case 'J':
		//return jack // part 1
		return joker // part 2
	case 'Q':
		return queen
	case 'K':
		return king
	case 'A':
		return ace
	}
	return -1
}

type Hand [5]Card

type HandWithBid struct {
	Hand
	bid int
}

func HandFrom(s string) Hand {
	var h Hand
	for i, r := range s {
		h[i] = CardFrom(r)
	}
	return h
}

func (h *Hand) isFiveOfAKind() bool {
	countsPerKind := countsOf(h[:]...)
	return len(countsPerKind) == 1
}

func (h *Hand) isFourOfAKind() bool {
	countsPerKind := countsOf(h[:]...)
	if len(countsPerKind) != 2 {
		return false
	}
	counts := valuesOf(countsPerKind)
	slices.Sort(counts)
	return slices.Equal(counts, []int{1, 4})
}

func (h *Hand) isFullHouse() bool {
	countsPerKind := countsOf(h[:]...)
	if len(countsPerKind) != 2 {
		return false
	}
	counts := valuesOf(countsPerKind)
	slices.Sort(counts)
	return slices.Equal(counts, []int{2, 3})
}

func (h *Hand) isThreeOfAKind() bool {
	countsPerKind := countsOf(h[:]...)
	if len(countsPerKind) != 3 {
		return false
	}
	counts := valuesOf(countsPerKind)
	slices.Sort(counts)
	return slices.Equal(counts, []int{1, 1, 3})
}

func (h *Hand) isTwoPair() bool {
	countsPerKind := countsOf(h[:]...)
	if len(countsPerKind) != 3 {
		return false
	}
	counts := valuesOf(countsPerKind)
	slices.Sort(counts)
	return slices.Equal(counts, []int{1, 2, 2})
}

func (h *Hand) isOnePair() bool {
	countsPerKind := countsOf(h[:]...)
	return len(countsPerKind) == 4
}

func (h *Hand) isHighCard() bool {
	countsPerKind := countsOf(h[:]...)
	return len(countsPerKind) == 5
}

func (h *Hand) Type() int {
	if h.isFiveOfAKind() {
		return fiveOfAKind
	}
	if h.isFourOfAKind() {
		return fourOfAKind
	}
	if h.isFullHouse() {
		return fullHouse
	}
	if h.isThreeOfAKind() {
		return threeOfAKind
	}
	if h.isTwoPair() {
		return twoPair
	}
	if h.isOnePair() {
		return onePair
	}
	return highCard
}

func (h *Hand) TypeAsString() string {
	if h.isFiveOfAKind() {
		return "fiveOfAKind"
	}
	if h.isFourOfAKind() {
		return "fourOfAKind"
	}
	if h.isFullHouse() {
		return "fullHouse"
	}
	if h.isThreeOfAKind() {
		return "threeOfAKind"
	}
	if h.isTwoPair() {
		return "twoPair"
	}
	if h.isOnePair() {
		return "onePair"
	}
	return "highCard"
}

func countsOf(items ...Card) map[Card]int {
	counts := make(map[Card]int)
	for _, item := range items {
		counts[item]++
	}
	if jokerCount := counts[joker]; jokerCount > 0 {
		delete(counts, joker)
		// find the best card the jokers should represent
		var bestCard Card
		for card, count := range counts {
			if count > counts[bestCard] || counts[bestCard] <= 1 && card > bestCard {
				bestCard = card
			}
		}
		counts[bestCard] += jokerCount
	}
	return counts
}

func valuesOf[K, V comparable](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

//go:embed input.txt
var input string

func main() {
	lines := strings.Split(input, "\n")
	hands := make([]HandWithBid, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		hand := HandFrom(fields[0])
		bid, _ := strconv.Atoi(fields[1])
		hands[i] = HandWithBid{hand, bid}
	}

	for _, hand := range hands {
		fmt.Println(hand, " -> ", hand.TypeAsString())
	}

	sort.Slice(hands, func(cardAIndex, cardBIndex int) bool {
		cardA := hands[cardAIndex].Hand
		cardB := hands[cardBIndex].Hand
		if cardA.Type() == cardB.Type() {
			for i := 0; i < len(cardA); i++ {
				if cardA[i] == cardB[i] {
					continue
				}
				return cardA[i] < cardB[i]
			}
		}
		return cardA.Type() < cardB.Type()
	})
	fmt.Println("Sorted hands: ", hands)

	winnings := 0
	for i, hand := range hands {
		winnings += hand.bid * (i + 1)
	}
	fmt.Println("Winnings: ", winnings)
}
