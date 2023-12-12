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

type iposition struct {
	row int
	col int
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

func adjacentGears(r irange, prevLine, line, nextLine string) []iposition {
	// Initialize to arbitrary small capacity.
	gears := make([]iposition, 0, 10)
	if r.start > 0 && line[r.start-1] == '*' {
		gears = append(gears, iposition{row: 0, col: r.start - 1})
	}
	if r.end < len(line) && line[r.end] == '*' {
		gears = append(gears, iposition{row: 0, col: r.end})
	}

	end := min(r.end+1, len(line))
	for i := max(0, r.start-1); i < end; i++ {
		if prevLine[i] == '*' {
			gears = append(gears, iposition{row: -1, col: i})
		}
		if nextLine[i] == '*' {
			gears = append(gears, iposition{row: 1, col: i})
		}
	}
	return gears
}

type gearMap map[iposition][]int

func (gears gearMap) updateParts(index int, prev, line, next string) error {
	sum := 0
	digitRanges := findDigitRanges(line)
	for _, r := range digitRanges {
		pos := adjacentGears(r, prev, line, next)
		if len(pos) == 0 {
			continue
		}
		num, err := strconv.Atoi(line[r.start:r.end])
		if err != nil {
			return err
		}

		for _, p := range pos {
			absp := iposition{row: index + p.row, col: p.col}
			gears[absp] = append(gears[absp], num)
		}
		sum += num
	}
	return nil
}

func makeString(b byte, l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = b
	}
	return string(bytes)
}

func computeGearRatioSum(inputPath string) (int, error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return -1, err
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)

	// Read first line to determine the length of all lines.
	lineBytes, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}

	gears := make(gearMap)
	// TODO: Fail if lines not equal in length.
	// TODO: Improve parsing and avoid the []byte - string dance.
	line := string(lineBytes)
	length := len(line)
	prevLine := makeString('.', length)
	i := 0
	for {
		nextLineBytes, _, err := reader.ReadLine()
		if err == io.EOF {
			break

		}
		if err != nil {
			return 0, err
		}

		nextLine := string(nextLineBytes)
		err = gears.updateParts(i, prevLine, line, nextLine)
		if err != nil {
			return 0, err
		}
		prevLine = line
		line = nextLine
		i++

	}
	nextLine := makeString('.', length)
	err = gears.updateParts(i, prevLine, line, nextLine)
	if err != nil {
		return 0, err
	}

	sum := 0
	for _, parts := range gears {
		if len(parts) == 2 {
			sum += parts[0] * parts[1]
		}
	}
	return sum, nil
}

func main() {
	inputPathFlag := flag.String("input_path", "", "Path to puzzle input file")
	flag.Parse()

	if *inputPathFlag == "" {
		log.Fatal("Flag --input_path must be non-empty!")
	}

	sum, err := computeGearRatioSum(*inputPathFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sum)
}
