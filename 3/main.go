package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"unicode"
)

type irange struct {
	start int // inclusive
	end   int // not inclusive
}

func findDigitRanges(s string) []irange {
	ranges := make([]irange, 0, 10) // Arbitrary initial capacity.
	inRange := false
	for i, r := range s {
		isDigit := unicode.IsDigit(r)
		if isDigit && !inRange {
			inRange = true
			ranges = append(ranges, irange{start: i})
		} else if !isDigit && inRange {
			inRange = false
			ranges[len(ranges)-1].end = i
		}
	}
	if inRange {
		ranges[len(ranges)-1].end = len(s)
	}
	return ranges
}

func isSymbol(b byte) bool {
	return b != '.' && !unicode.IsDigit(rune(b))
}

func isPartNumber(r irange, prevLine, line, nextLine string) bool {
	if r.start > 0 && isSymbol(line[r.start-1]) {
		return true
	}
	if r.end < len(line) && isSymbol(line[r.end]) {
		return true
	}

	end := min(r.end+1, len(line))
	for i := max(0, r.start-1); i < end; i++ {
		if isSymbol(prevLine[i]) || isSymbol(nextLine[i]) {
			return true
		}
	}

	return false
}

func partNumberSum(prevLine, line, nextLine string) (int, error) {
	sum := 0
	digitRanges := findDigitRanges(line)
	for _, r := range digitRanges {
		if !isPartNumber(r, prevLine, line, nextLine) {
			continue
		}
		num, err := strconv.Atoi(line[r.start:r.end])
		if err != nil {
			return 0, err
		}
		sum += num
	}
	return sum, nil
}

func makeString(b byte, l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = b
	}
	return string(bytes)
}

func computePartNumberSum(inputPath string) (int, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return -1, err
	}
	defer inputFile.Close()

	sum := 0
	reader := bufio.NewReader(inputFile)

	// Read first line to determine the length of all lines.
	lineBytes, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}
	// TODO: Fail if lines not equal in length.
	// TODO: Improve parsing and avoid the []byte - string dance.
	line := string(lineBytes)
	length := len(line)
	prevLine := makeString('.', length)
	for {
		nextLineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break

		}
		if err != nil {
			return 0, err
		}

		nextLine := string(nextLineBytes)
		lineSum, err := partNumberSum(prevLine, line, nextLine)
		if err != nil {
			return 0, err
		}
		sum += lineSum
		prevLine = line
		line = nextLine

	}
	nextLine := makeString('.', length)
	s, err := partNumberSum(prevLine, line, nextLine)
	sum += s
	return sum, err
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	sum, err := computePartNumberSum(*inputPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sum)
}
