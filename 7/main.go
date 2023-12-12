package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const Joker int = 1

var cardRanks = map[rune]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': Joker,
	'Q': 12,
	'K': 13,
	'A': 14,
}

type handType = int

const (
	UnknownHand handType = iota
	HighCard
	OnePair
	TwoPair
	ThreeOfAKind
	FullHouse
	FourOfAKind
	FiveOfAKind
)

type hand struct {
	cards [5]int
	bid   int
}

func (this hand) Type() handType {
	cardCounts := make(map[int]int)
	jokerCount := 0
	for _, card := range this.cards {
		if card == Joker {
			jokerCount++
			continue
		}
		c, ok := cardCounts[card]
		if !ok {
			c = 0
		}
		cardCounts[card] = c + 1
	}

	maxCount := 0
	maxCard := 0
	for card, c := range cardCounts {
		if c > maxCount {
			maxCount = c
			maxCard = card
		}
	}
	cardCounts[maxCard] += jokerCount

	var counts [6]int
	for _, c := range cardCounts {
		counts[c] += 1
	}
	if counts[5] == 1 {
		return FiveOfAKind
	}
	if counts[4] == 1 {
		return FourOfAKind
	}
	if counts[3] == 1 && counts[2] == 1 {
		return FullHouse
	}
	if counts[3] == 1 {
		return ThreeOfAKind
	}
	if counts[2] == 2 {
		return TwoPair
	}
	if counts[2] == 1 {
		return OnePair
	}
	if counts[1] == 5 {
		return HighCard
	}
	return UnknownHand
}

func (this hand) Less(other hand) bool {
	thisType := this.Type()
	otherType := other.Type()
	if thisType < otherType {
		return true
	}
	if thisType > otherType {
		return false
	}

	numCards := len(this.cards)
	for i := 0; i < numCards; i++ {
		if this.cards[i] < other.cards[i] {
			return true
		}
		if this.cards[i] > other.cards[i] {
			return false
		}
	}
	return false
}

func parseHand(s string) (h hand, err error) {
	cardsAndBid := strings.Split(s, " ")
	if len(cardsAndBid) != 2 {
		err = fmt.Errorf("Cannot parse hand from: %s", s)
		return
	}

	cards := cardsAndBid[0]
	if len(cards) != 5 {
		err = fmt.Errorf("Cannot parse cards from: %s", s)
		return
	}

	for i, c := range cards {
		rank, ok := cardRanks[c]
		if !ok {
			err = fmt.Errorf("Not a card: %s", c)
			return
		}
		h.cards[i] = rank
	}

	h.bid, err = strconv.Atoi(cardsAndBid[1])
	return
}

func loadHands(inputPath string) (hands []hand, err error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	hands = make([]hand, 0, 10) // arbitrary capacity
	for scanner.Scan() {
		token := scanner.Text()
		if token == "" {
			continue
		}
		var h hand
		h, err = parseHand(token)
		if err != nil {
			return
		}
		hands = append(hands, h)
	}

	err = scanner.Err()
	if err != nil {
		return
	}
	return
}

func totalScore(hands []hand) int {
	sort.Slice(hands, func(i, j int) bool { return hands[i].Less(hands[j]) })
	score := 0
	for i, h := range hands {
		score += (i + 1) * h.bid
	}
	return score
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	hands, err := loadHands(*inputPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(totalScore(hands))
}
