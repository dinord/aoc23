package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// TODO: Make error propagation less intrusive.
// There must be a way...

func parseInts(s string) ([]int, error) {
	ints := make([]int, 0, 10) // arbitrary capacity
	tokens := strings.Split(strings.Trim(s, " "), " ")
	for _, t := range tokens {
		if t == "" {
			continue
		}
		i, err := strconv.Atoi(t)
		if err != nil {
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func parseScratchCard(s string) (wins []int, scratches []int, err error) {
	// Skip the prefix, e.g. "Card 123456:"
	start := strings.Index(s, ":")
	if start == -1 {
		return nil, nil, fmt.Errorf("Expected card after `:`, got: %s", s)
	}

	card := s[start+1:]
	split := strings.Index(card, "|")
	if split == -1 {
		return nil, nil, fmt.Errorf("Expected card with `|`, got: %s", card)
	}

	winsText := card[:split]
	scratchesText := card[split+1:]

	wins, err = parseInts(winsText)
	if err != nil {
		return nil, nil, err
	}
	scratches, err = parseInts(scratchesText)
	if err != nil {
		return nil, nil, err
	}
	return wins, scratches, nil
}

func computeScratchScore(wins []int, scratches []int) int {
	matches := 0
	winSet := make(map[int]struct{})
	for _, w := range wins {
		winSet[w] = struct{}{}
	}

	for _, s := range scratches {
		if _, match := winSet[s]; match {
			matches++
		}
	}
	if matches > 0 {
		return 1 << (matches - 1)
	} else {
		return 0
	}
}

func computeTotalScratchScore(inputPath string) (points int, err error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return -1, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	score := 0
	for scanner.Scan() {
		wins, scratches, err := parseScratchCard(scanner.Text())
		if err != nil {
			return -1, err
		}
		score += computeScratchScore(wins, scratches)
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}

	return score, nil
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	points, err := computeTotalScratchScore(*inputPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(points)

}
